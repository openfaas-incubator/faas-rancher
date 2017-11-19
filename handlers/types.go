package handlers

import "net/http"

// HTTPDoer is a Http Wrapper
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
