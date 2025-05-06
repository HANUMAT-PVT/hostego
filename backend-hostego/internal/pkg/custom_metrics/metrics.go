package custom_metrics

// import (
// 	"backend-hostego/internal/hostego/utilities"
// 	"backend-hostego/internal/pkg/logger"
// 	"fmt"
// )

// var log = logger.GetLogger()

// type MetricData struct {
// 	ApiPath    string
// 	UserId     string
// 	Source     string
// 	StatusCode int
// 	Exchange   string
// 	Latency    int64
// }

// func SendErrorCountMetrics(apiPath, userId, source string) {
// 	defer utilities.GoRoutinePanicHandler()
// 	metric, err := GetMetric(API_NAME_USERID_SOURCE_ERROR_GAUGE)
// 	if err != nil {
// 		log.Errorf("Not able to get custom metric for %v due to : %v", API_NAME_USERID_SOURCE_ERROR_GAUGE, err.Error())
// 		return
// 	}

// 	metric.Inc(1, apiPath, userId, source)
// }

// func SendRequestCountMetrics(apiPath, userId, source string) {
// 	defer utilities.GoRoutinePanicHandler()
// 	metric, err := GetMetric(API_NAME_USERID_SOURCE_TOTAL_GAUGE)
// 	if err != nil {
// 		log.Errorf("Not able to get custom metric for %v due to : %v", API_NAME_USERID_SOURCE_TOTAL_GAUGE, err.Error())
// 		return
// 	}

// 	metric.Inc(1, apiPath, userId, source)
// }

// func SendHistogramLatencyMetrics(apiPath, userId, statusCode, source string, latency float64) {
// 	defer utilities.GoRoutinePanicHandler()
// 	metric, err := GetMetric(API_NAME_USERID_SOURCE_HISTOGRAM_LATENCY)
// 	if err != nil {
// 		log.Errorf("Not able to get custom metric for %v due to : %v", API_NAME_USERID_SOURCE_HISTOGRAM_LATENCY, err.Error())
// 		return
// 	}

// 	metric.SetValue(latency, apiPath, userId, fmt.Sprintf("%v", statusCode), source)
// }

// func SendExchangeRequestCountMetrics(apiPath, userId, source, exchange string) {
// 	defer utilities.GoRoutinePanicHandler()
// 	metric, err := GetMetric(API_NAME_USERID_SOURCE_EXCHANGE_GAUGE)
// 	if err != nil {
// 		log.Errorf("Not able to get custom metric for %v due to : %v", API_NAME_USERID_SOURCE_EXCHANGE_GAUGE, err.Error())
// 		return
// 	}

// 	metric.Inc(1, apiPath, userId, source, exchange)
// }

// func SendErrorAndLatencyMetrics(apiPath, userId, source, exchange string, statusCode int, latency int64) {
// 	defer utilities.GoRoutinePanicHandler()
// 	SendExchangeHistogramLatencyMetrics(apiPath, userId, source, exchange, statusCode, latency)
// 	SendErrorCountMetrics(apiPath, userId, source)
// }

// func SendExchangeHistogramLatencyMetrics(apiPath, userId, source, exchange string, statusCode int, latency int64) {
// 	defer utilities.GoRoutinePanicHandler()
// 	metric, err := GetMetric(API_NAME_USERID_SOURCE_EXCHANGE_HISTOGRAM_LATENCY)
// 	if err != nil {
// 		log.Errorf("Not able to get custom metric for %v due to : %v", API_NAME_USERID_SOURCE_EXCHANGE_HISTOGRAM_LATENCY, err.Error())
// 		return
// 	}

// 	metric.SetValue(float64(latency), apiPath, userId, source, exchange, fmt.Sprintf("%v", statusCode))
// }
