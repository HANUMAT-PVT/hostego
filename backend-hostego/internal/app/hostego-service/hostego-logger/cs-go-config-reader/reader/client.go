package reader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ConfigReader struct {
	baseURL             string
	hostType            string
	tenant              string
	service             string
	configType          string
	requestTimeoutSec   time.Duration
	localConfigFilePath string
}

func NewConfigReader(baseURL, hostType, tenant, service string, requestTimeoutSec time.Duration, localConfigFilePath string) *ConfigReader {
	return &ConfigReader{
		baseURL:             baseURL,
		hostType:            hostType,
		tenant:              tenant,
		service:             service,
		configType:          "stable",
		requestTimeoutSec:   requestTimeoutSec,
		localConfigFilePath: localConfigFilePath,
	}
}

func (c *ConfigReader) SetConfigType(configType string) {
	if configType == "" {
		configType = ""
	}
	c.configType = configType
}

func (c *ConfigReader) GetConfig(key string, options map[string]string) (string, error) {
	key = strings.ToUpper(key)
	if strings.Contains(strings.ToLower(c.hostType), "local") {
		localConfig := c.getConfigFromLocal(key)
		if localConfig != "" {
			fmt.Printf("config returned from local file for %s, value: %v\n", key, localConfig)
			return localConfig, nil
		}
	}

	headers := make(http.Header)
	for k, v := range options {
		headers.Set(k, v)
	}

	params := map[string]string{
		"tenant":   c.tenant,
		"service":  c.service,
		"is_draft": fmt.Sprintf("%t", strings.ToLower(c.configType) == "draft"),
		"key":      key,
	}

	paramStr := ""
	for key, value := range params {
		prefix := "&"
		if paramStr == "" {
			prefix = "?"
		}

		paramStr = fmt.Sprintf("%s%s%s=%s", paramStr, prefix, key, value)
	}

	url := fmt.Sprintf("%s/api/v1/config%s", c.baseURL, paramStr)

	client := &http.Client{Timeout: c.requestTimeoutSec * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("An error occurred while getting key: %v, error: %v", params, err)
	}
	defer response.Body.Close()

	fmt.Printf("called config reader with URL %s\n", response.Request.URL)

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Status code %d received while fetching %v", response.StatusCode, params)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("An error occurred while reading response body: %v", err)
	}

	var configValue map[string]interface{}
	err = json.Unmarshal(data, &configValue)
	if err != nil {
		return "", fmt.Errorf("An error occurred while unmarshalling response data: %v", err)
	}

	result, ok := configValue["is_non_dict_data"]
	if ok {
		isNonConfigData := result.(bool)
		if isNonConfigData {
			return configValue["config"].(string), nil
		}
	}

	return configValue["data"].(string), nil
}

func (c *ConfigReader) getConfigFromLocal(key string) string {
	fileContent, err := ioutil.ReadFile(c.localConfigFilePath)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			errMessage := fmt.Sprintf("The file '%s' does not exist. Please create this file.", c.localConfigFilePath)
			fmt.Println(errMessage)
			return ""
		} else {
			errMessage := fmt.Sprintf("An error occurred while reading file '%s': %v", c.localConfigFilePath, err)
			fmt.Println(errMessage)
			return ""
		}
	}

	var configs map[string]string
	err = json.Unmarshal(fileContent, &configs)
	if err != nil {
		errMessage := fmt.Sprintf("An error occurred while unmarshalling local config file: %v", err)
		fmt.Println(errMessage)
		return ""
	}

	if value, ok := configs[key]; ok {
		return value
	}

	return ""
}
