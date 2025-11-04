package main

import (
	"fmt"
	"net/http"
	"net/url"
	"nycaresprojectwelcomer/internal/httphelper"
	"strings"
)

// buildPayload constructs the form data for login
func buildPayload(username, password string) map[string]string {
    return map[string]string{
        "username": username,
        "password": password,
        "form_id":  "user_login_form",
        "op":       "Log in",
    }
}

// encodePayload URL-encodes Form data
func encodePayload(payload map[string]string) string {
    formData := url.Values{}
    for k, v := range payload {
        formData.Set(k, v)
    }
    return formData.Encode()
}

// buildRequest creates a POST request with headers
func buildLoginRequest(baseURL string, creds Credentials) (*http.Request, error) {
    payload := buildPayload(creds.Username, creds.Password)
    encoded := encodePayload(payload)

    req, err := http.NewRequest("POST", baseURL, strings.NewReader(encoded))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("User-Agent", "Mozilla/5.0")
    return req, nil
}

// extractCookies extracts cookies from the response
func extractCookies(resp *http.Response) (map[string]string, error) {
    cookies := make(map[string]string)
    for _, c := range resp.Cookies() {
        cookies[c.Name] = c.Value
    }
    if len(cookies) == 0 {
        return nil, fmt.Errorf("no cookies were set by the server")
    }
    return cookies, nil
}

// Credentials holds username/password for login
type Credentials struct {
	Username string
	Password string
}

// Login performs the login and returns only the cookies
func Login(client *http.Client, baseURL string, creds Credentials) (map[string]string, error) {
    if creds.Username == "" {
        return nil, fmt.Errorf("ACCOUNT_USERNAME is required")
    }
    if creds.Password == "" {
        return nil, fmt.Errorf("ACCOUNT_PASSWORD is required")
    }

    req, err := buildLoginRequest(baseURL, creds)
    if err != nil {
        return nil, err
    }

    resp, err := httphelper.SendRequest(client, req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close() // close after we read cookies

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("login attempt failed with status %d", resp.StatusCode)
    }

    cookies, err := extractCookies(resp)
    if err != nil {
        return nil, err
    }

    return cookies, nil
}