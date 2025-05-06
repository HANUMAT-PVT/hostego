package custom_metrics

// import (
// 	"bitbucket.org/coinswitch/cs-gometrics/enums"
// 	"bitbucket.org/coinswitch/cs-gometrics/factory"
// )

// func GetAllCustomMetrics() []factory.CustomMetricFactory {
// 	samplePercentage := getSamplePercentage()
// 	data := make(map[string]interface{})
// 	data["sample_percentage"] = samplePercentage
// 	buckets := []float64{10, 20, 30, 50, 100, 200, 500, 750, 1000, 1500, 2000, 2500, 3000, 4000, 5000}

// 	sourceCountLabels := []string{"api_name", "user_id", "source"}
// 	apiNameUserIdSourceErrorCountGauge := factory.NewCustomMetric(API_NAME_USERID_SOURCE_ERROR_GAUGE, "API Trading gauge for errors", enums.GAUGE, sourceCountLabels, data)

// 	apiNameUserIdSourceTotalReqCountGauge := factory.NewCustomMetric(API_NAME_USERID_SOURCE_TOTAL_GAUGE, "API Trading gauge for api calls", enums.GAUGE, sourceCountLabels, data)

// 	sourceRequestTimeLatenciesLabels := []string{"api_name", "user_id", "status_code", "source"}
// 	sourceRequestTimeData := make(map[string]interface{})
// 	sourceRequestTimeData["buckets"] = buckets
// 	sourceRequestTimeData["sample_percentage"] = samplePercentage
// 	apiNameUserIdSourceRequestTimeHistogram := factory.NewCustomMetric(API_NAME_USERID_SOURCE_HISTOGRAM_LATENCY, "API Trading histogram for api latencies", enums.HISTOGRAM, sourceRequestTimeLatenciesLabels, sourceRequestTimeData)

// 	sourceExchangeLabels := []string{"api_name", "user_id", "source", "exchange"}
// 	apiNameUserIdSourceExchangeTotalReqCountGauge := factory.NewCustomMetric(API_NAME_USERID_SOURCE_EXCHANGE_GAUGE, "API Trading gauge for api calls per exchange", enums.GAUGE, sourceExchangeLabels, data)

// 	sourceExchangeRequesTimeLabels := []string{"api_name", "user_id", "source", "exchange", "status_code"}
// 	apiNameUserIdSourceExchangeRequestTimeHistogram := factory.NewCustomMetric(API_NAME_USERID_SOURCE_EXCHANGE_HISTOGRAM_LATENCY, "API Trading histogram for api latencies for create and cancel", enums.HISTOGRAM, sourceExchangeRequesTimeLabels, sourceRequestTimeData)

// 	var customMetrics = []factory.CustomMetricFactory{
// 		apiNameUserIdSourceErrorCountGauge,
// 		apiNameUserIdSourceTotalReqCountGauge,
// 		apiNameUserIdSourceRequestTimeHistogram,
// 		apiNameUserIdSourceExchangeTotalReqCountGauge,
// 		apiNameUserIdSourceExchangeRequestTimeHistogram,
// 	}

// 	return customMetrics
// }
