package httpclient

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPError represents an HTTP-specific error
type HTTPError struct {
	StatusCode int
	Status     string
	URL        string
	Method     string
	Body       []byte
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	return fmt.Sprint("%s %s returned %d %s", e.Method, e.URL, e.StatusCode, e.Status)
}

// IsSuccess chrecks if the response status code indicates success
func IsSuccess(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// CheckResponse creates an error from an unsuccessful response
func CheckResponse(resp *http.Response) error {
	if IsSuccess(resp) {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return &HTTPError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		URL:        resp.Request.URL.String(),
		Method:     resp.Request.Method,
		Body:       body,
	}
}
