package myHttp

import (
	"net/http"
)

// HttpClientInterface helps mocking http.Client, we generally use only Do function of the client
// using this interface we can pass http.Client as an implementation and mock it in implementation
type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}
