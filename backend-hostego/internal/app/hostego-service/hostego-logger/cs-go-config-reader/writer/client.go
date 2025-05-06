// tayyar hoke bhandare me gaya tha, andar gaya toh halwa khatam
// bahar aaya toh chappal gayab

package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const LOCAL_CONFIG_FILE_PATH = "../config.json"

type ConfigWriter struct {
	BaseURL           string
	HostType          string
	Tenant            string
	Service           string
	MakerChecker      string
	RequestTimeoutSec int
}

func NewConfigWriter(baseURL, hostType, tenant, service, makerChecker string, requestTimeoutSec int) *ConfigWriter {
	return &ConfigWriter{
		BaseURL:           baseURL,
		HostType:          hostType,
		Tenant:            tenant,
		Service:           service,
		MakerChecker:      makerChecker,
		RequestTimeoutSec: requestTimeoutSec,
	}
}

type ConfigWriterException struct {
	message string
}

func (e *ConfigWriterException) Error() string {
	return e.message
}

func (cw *ConfigWriter) SetConfig(key, field string, value interface{}, headers map[string]string) (map[string]interface{}, error) {
	if cw.isLocalHost() {
		return nil, &ConfigWriterException{message: fmt.Sprintf("Please add your config key value in %s file", LOCAL_CONFIG_FILE_PATH)}
	}

	body := map[string]interface{}{
		"config_key": map[string]string{
			"tenant":  cw.Tenant,
			"service": cw.Service,
			"key":     key,
		},
		"maker_checker": cw.MakerChecker,
		"field":         field,
		"value":         value,
	}

	url := cw.BaseURL + "/api/v1/config/partial"
	log.Printf("calling config service with body %+v\n", body)

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: time.Duration(cw.RequestTimeoutSec) * time.Second,
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, &ConfigWriterException{message: string(responseBody)}
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		return nil, err
	}

	return responseData["data"].(map[string]interface{}), nil
}

func (cw *ConfigWriter) isLocalHost() bool {
	return containsIgnoreCase(cw.HostType, "local")
}

func containsIgnoreCase(str, substr string) bool {
	return len(str) >= len(substr) && bytes.EqualFold([]byte(str)[len(str)-len(substr):], []byte(substr))
}

// func main() {
// 	// Example usage
// 	configWriter := NewConfigWriter("http://example.com", "SomeHostType", "SomeTenant", "SomeService", "SomeMakerChecker", 10)

// 	headers := map[string]string{"Authorization": "Bearer your_token"}
// 	key := "some_key"
// 	field := "some_field"
// 	value := "some_value"

// 	response, err := configWriter.SetConfig(key, field, value, headers)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Printf("Response: %+v\n", response)
// }
