package httphelper

import (
	"io"
	"net/http"
)

// sendRequest executes the HTTP request and returns the response
func SendRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	return client.Do(req) // resp.Body is closed by caller
}

// ReadBody is a helper to read response body as []byte.
func ReadBody(resp *http.Response) ([]byte, error) {
	return io.ReadAll(resp.Body)
}
