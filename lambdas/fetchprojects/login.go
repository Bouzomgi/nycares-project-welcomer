package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/httphelper"
)

// buildPayload constructs the form data for login
func buildPayload(username, password string) url.Values {
	form := url.Values{}
	form.Set("form_id", "user_login_form")
	form.Set("name", username)
	form.Set("pass", password)

	return form
}

// buildRequest creates a POST request with headers
func buildLoginRequest(baseUrl string, creds Credentials) (*http.Request, error) {
	payload := buildPayload(creds.Username, creds.Password)
	encoded := payload.Encode()
	loginUrl := baseUrl + LoginPath

	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(encoded))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func validateLogin(client *http.Client, baseUrl string) error {
	u, _ := url.Parse(baseUrl)
	if client.Jar.Cookies(u) == nil {
		return fmt.Errorf("Login failed. No cookie set on %s", baseUrl)
	}
	return nil
}

// Credentials holds username/password for login
type Credentials struct {
	Username string
	Password string
}

// Login performs the login and returns only the cookies
func Login(client *http.Client, baseUrl string, creds Credentials) error {
	if creds.Username == "" {
		return fmt.Errorf("ACCOUNT_USERNAME is required")
	}
	if creds.Password == "" {
		return fmt.Errorf("ACCOUNT_PASSWORD is required")
	}

	req, err := buildLoginRequest(baseUrl, creds)
	if err != nil {
		return err
	}

	resp, err := httphelper.SendRequest(client, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("login attempt failed with status %d", resp.StatusCode)
	}

	err = validateLogin(client, baseUrl)
	if err != nil {
		return err
	}

	return nil
}
