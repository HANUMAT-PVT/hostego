package custom_metrics

// import (
// 	"backend-hostego/internal/hostego/constants/config_constants"

// 	"github.com/spf13/viper"
// )

// const (
// 	API_NAME_USERID_SOURCE_ERROR_GAUGE                    = "api_name_userid_source_error_gauge"
// 	API_NAME_USERID_SOURCE_TOTAL_GAUGE                    = "api_name_userid_source_total_gauge"
// 	API_NAME_USERID_SOURCE_HISTOGRAM_LATENCY              = "api_name_userid_source_histogram_latency"
// 	API_NAME_USERID_SOURCE_EXCHANGE_GAUGE                 = "api_name_userid_source_exchange_gauge"
// 	API_NAME_USERID_SOURCE_EXCHANGE_HISTOGRAM_LATENCY     = "api_name_userid_source_exchange_histogram_latency"
// 	DefaultSamplePercentage                           int = 100
// )

// func getSamplePercentage() int {
// 	samplePercentage := viper.GetInt(config_constants.CUSTOM_METRICS_SAMPLE_PERCENTAGE)
// 	if samplePercentage == 0 {
// 		return DefaultSamplePercentage
// 	}
// 	return samplePercentage
// }
