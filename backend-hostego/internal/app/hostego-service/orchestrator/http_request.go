package orchestrator

import (
	"backend-hostego/internal/app/hostego-service/constants"
	"backend-hostego/internal/app/hostego-service/constants/config_constants"
	dto2 "backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/pkg/common_utils"
	myHttp "backend-hostego/internal/pkg/http"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/viper"
)

var (
	once   sync.Once
	client *http.Client
)

func init() {
	t := http.DefaultTransport.(*http.Transport).Clone() // TODO : Make all the three values through config
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	once.Do(func() {
		client = &http.Client{
			Transport: t,
		}
	})
}

type HttpOrchestrator struct {
	Client myHttp.HttpClientInterface
}

func NewHttpOrchestrator() *HttpOrchestrator {
	return &HttpOrchestrator{
		Client: client,
	}
}

func (h *HttpOrchestrator) GetRequest(rCtx dto2.ReqCtx, request dto2.HttpRequest) (dto2.HttpResponse, error) {
	log := rCtx.Log
	var req *http.Request
	timeout := getTimeoutForRequest(request.Timeout)
	// This timeout context is set on the http client. Decide what timeout we need to put in config
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, request.Url, nil)
	setHeaderInPayload(req, request.Header)
	log.Infof("Request Details from Orchestrator, request = %v", request)

	return h.makeRequest(rCtx, req, request)
}

func (h *HttpOrchestrator) DeleteRequest(rCtx dto2.ReqCtx, request dto2.HttpRequest) (dto2.HttpResponse, error) {
	log := rCtx.Log
	var req *http.Request
	timeout := getTimeoutForRequest(request.Timeout)
	// This timeout context is set on the http client. Decide what timeout we need to put in config
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	req, _ = http.NewRequestWithContext(ctx, http.MethodDelete, request.Url, nil)
	setHeaderInPayload(req, request.Header)
	log.Infof("Request Details from Orchestrator, request = %v", request)

	return h.makeRequest(rCtx, req, request)
}

func (h *HttpOrchestrator) PostRequest(rCtx dto2.ReqCtx, request dto2.HttpRequest) (dto2.HttpResponse, error) {
	log := rCtx.Log
	var req *http.Request

	timeout := getTimeoutForRequest(request.Timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	// Make bifurcation on the basis of request format
	// json marshall if request type is json
	switch request.RequestFormat {
	case constants.REQUEST_FORMAT_XML, constants.REQUEST_FORMAT_PIPE_DELIMITED:
		log.Infof("payload for the request of request-type: %v is : %v", request.RequestFormat, request.Body)
		req, _ = http.NewRequestWithContext(ctx, request.RequestType, request.Url, bytes.NewBuffer([]byte(fmt.Sprint(request.Body))))
	case constants.REQUEST_FORMAT_JSON:
		jsonPayload, err := json.Marshal(request.Body)
		log.Infof("payload for the request request-type: %v is : %v", request.RequestFormat, string(jsonPayload))
		if err != nil {
			log.Errorf("request is %v, error while coverting struct to Json, Error is %v", request, err.Error())
			return dto2.HttpResponse{}, err
		}
		req = getHttpRequestWithContext(ctx, request, request.Body.([]byte))
	// TODO: refactor to make format a mandatory field and change in all outbound configs, return empty response and error in default
	default:
		jsonPayload, err := json.Marshal(request.Body)
		log.Infof("payload for the request request-type: %v is : %v", request.RequestFormat, string(jsonPayload))
		if err != nil {
			log.Errorf("request is %v, error while coverting struct to Json, Error is %v", request, err.Error())
			return dto2.HttpResponse{}, err
		}
		req = getHttpRequestWithContext(ctx, request, jsonPayload)
	}

	setHeaderInPayload(req, request.Header)

	log.Infof("Request Details from Orchestrator, request = %v", request)

	return h.makeRequest(rCtx, req, request)
}

func getHttpRequestWithContext(ctx context.Context, request dto2.HttpRequest, payload []byte) *http.Request {
	if val, ok := request.Header["Content-Type"]; ok {
		if val == config_constants.CONTENT_TYPE_XXX_FORM_URLENCODED {
			return getRequestForXXXFormEncodedURL(ctx, request)
		}
	}
	req, _ := http.NewRequestWithContext(ctx, request.RequestType, request.Url, bytes.NewBuffer(payload))
	return req
}

func getRequestForXXXFormEncodedURL(ctx context.Context, request dto2.HttpRequest) *http.Request {
	data := url.Values{}
	for key, value := range request.Body.(map[string]interface{}) {
		if keyVal, ok := value.(string); ok {
			data.Add(key, keyVal)
		}
	}
	encodedData := data.Encode()
	req, _ := http.NewRequestWithContext(ctx, request.RequestType, request.Url, strings.NewReader(encodedData))
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	return req
}

func (h *HttpOrchestrator) makeRequest(rCtx dto2.ReqCtx, req *http.Request, request dto2.HttpRequest) (dto2.HttpResponse, error) {
	log := rCtx.Log
	newrelicExternalSegment := newrelic.StartExternalSegment(rCtx.NewRelicTxn, req)
	externalApiTatTimer := time.Now()
	retryCount := request.RetryCount

	// Make request 1st time
	resp, err := h.Client.Do(req)

	defer func() {
		statusCode := -1
		if resp != nil {
			statusCode = resp.StatusCode
		}
		newrelicExternalSegment.SetStatusCode(statusCode)
		newrelicExternalSegment.End()
	}()
	if err != nil {
		// This will be the case when - connection not initialized | there is some issue from client sid
		// We are checking this here because if err!=nil then resp will be nil
		log.Errorf("Request URL: %v Error while fetching response is %v. Error is: %v", req.URL.Host+req.URL.Path, resp, err.Error())
		return dto2.HttpResponse{}, err
	}

	// RETRY_STATUS_CODE = 502
	if constants.RETRY_STATUS_CODE == resp.StatusCode {
		currentRetryCount := 1
		log.Infof("Request is eligible for retry for URL: %v. Current status code:%v",
			req.URL.Host+req.URL.Path, resp.StatusCode)

		// If retry count is present in external entity config for a specific use case,
		for currentRetryCount <= retryCount {
			currentRetryCount += 1

			// Retrying upon receiving 502
			resp, err = h.Client.Do(req)

			// timeout or due to client closed
			if err != nil {
				// This can be the scenario when resp1 status!=502, and we got some error on retry
				log.Errorf("Request URL: %v, Current response :%v, Current retry count: %v. Error while calling request, error is: %v",
					req.URL.Host+req.URL.Path, resp, currentRetryCount, err.Error())
				return dto2.HttpResponse{}, err
			}

			// Checking if we need to retry again
			if constants.RETRY_STATUS_CODE == resp.StatusCode {
				log.Errorf("Request URL: %v, Current status code:%v, Current retry count: %v", req.URL.Host+req.URL.Path, resp.StatusCode, currentRetryCount)

				continue
			} else {
				break
			}
		}
	}

	// Clean up activities
	// TODO : check if this defer function should be before return
	defer func() {
		statusCode := 0
		if resp != nil {
			statusCode = resp.StatusCode
		}

		log.Infof("external API call details, external_api_url: %v, external_api_status_code: %v, external_api_tat: %v",
			req.URL.Host+req.URL.Path, statusCode, time.Since(externalApiTatTimer).Milliseconds())

		resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		log.Errorf("error in parsing response while calling %v. Request is %v. Error is %v", request.Url, request, err)

		return dto2.HttpResponse{StatusCode: resp.StatusCode}, err
	}
	httpResponse := dto2.HttpResponse{Header: resp.Header, Body: body, StatusCode: resp.StatusCode}

	//Needs to disable logs for some requests, which contains critical data
	for _, urlString := range constants.DISABLE_EXTERNAL_URL_LOGS {
		if !(strings.Contains(request.Url, urlString) && common_utils.Contains(config_constants.ProdHosts, viper.GetString(config_constants.VKEYS_HOST_TYPE))) {
			log.Infof("external api response details, response: %v", httpResponse)
		}
	}

	return httpResponse, nil
}

func setHeaderInPayload(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		(*req).Header.Set(key, value)
	}
}

func getTimeoutForRequest(requestTimeout int) int {
	if requestTimeout == 0 {
		return 10 // TODO : need to make this a config value. Will do it once start making DB calls
	} else {
		return requestTimeout
	}
}
