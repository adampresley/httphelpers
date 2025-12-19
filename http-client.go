package httphelpers

import "net/http"

/*
HttpClient is an interface for making HTTP requests.
*/
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
