package dto

import (
	configconstants2 "backend-hostego/internal/app/hostego-service/constants/config_constants"
	"fmt"

	"github.com/spf13/viper"
)

// TODO : Add timeout parameter
type HttpResponse struct {
	StatusCode int
	Header     map[string][]string
	Body       []byte
}

func (httpResponse HttpResponse) String() string {
	return fmt.Sprintf("status_code is %v, body is %v", httpResponse.StatusCode, string(httpResponse.Body))
}

type HttpRequest struct {
	RequestType        string
	RequestFormat      string
	Url                string
	Header             map[string]string
	Body               interface{}
	Timeout            int
	RetryCount         int
	EncryptionRequired bool
	DecryptionRequired bool
	ChecksumRequired   bool
}

func (requestType *HttpRequest) SetBody(body interface{}) {
	requestType.Body = body
}

func (httpRequest HttpRequest) String() string {
	if viper.GetString(configconstants2.VKEYS_HOST_TYPE) != configconstants2.PROD_HOST {
		return fmt.Sprintf("request_type is %v, url is %v, header is %v, body is %v, timeout is %v",
			httpRequest.RequestType, httpRequest.Url, httpRequest.Header, httpRequest.Body, httpRequest.Timeout)
	}
	return fmt.Sprintf(
		"request_type is %v, url is %v, body is %v, timeout is %v, retryCount is %v",
		httpRequest.RequestType, httpRequest.Url, httpRequest.Body, httpRequest.Timeout, httpRequest.RetryCount)
}
