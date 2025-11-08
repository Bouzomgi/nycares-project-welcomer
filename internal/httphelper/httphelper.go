package httphelper

import (
	"io"
	"net/http"
	"net/url"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

// sendRequest executes the HTTP request and returns the response
func SendRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	return client.Do(req) // resp.Body is closed by caller
}

// ReadBody is a helper to read response body as []byte.
func ReadBody(resp *http.Response) ([]byte, error) {
	return io.ReadAll(resp.Body)
}

// setCookiesOnClient adds cookies from Auth into the client jar.
func SetCookiesOnClient(client *http.Client, baseUrl string, auth models.Auth) error {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}

	for _, c := range auth.Cookies {
		client.Jar.SetCookies(u, []*http.Cookie{
			{
				Name:   c.Name,
				Value:  c.Value,
				Domain: c.Domain,
				Path:   c.Path,
			},
		})
	}

	return nil
}
