package custom_metrics

// import (
// 	"net/http"
// 	"sync"

// 	"bitbucket.org/coinswitch/cs-gometrics/destinations"
// 	"bitbucket.org/coinswitch/cs-gometrics/factory"
// 	"bitbucket.org/coinswitch/cs-gometrics/models"
// )

// type CsGoMetrics struct {
// 	customMetrics []factory.CustomMetricFactory
// 	Client        destinations.PromClient
// 	HttpHandler   http.Handler
// }

// var once sync.Once
// var csGoMetrics *CsGoMetrics

// func GetInstance() *CsGoMetrics {
// 	once.Do(func() {
// 		customMetrics := GetAllCustomMetrics()
// 		client := destinations.NewPromClient("hostego", customMetrics)
// 		httpHandler := client.RegisterMetrics()
// 		csGoMetrics = &CsGoMetrics{customMetrics: customMetrics, Client: client, HttpHandler: httpHandler}
// 	})
// 	return csGoMetrics
// }

// func GetMetric(name string) (models.Metrics, error) {
// 	return csGoMetrics.Client.Dest.Get(name)
// }
