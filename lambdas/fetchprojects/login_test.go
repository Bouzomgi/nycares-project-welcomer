package main

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
)

// mock server for login
func setupMockLoginServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that method is POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		username := r.Form.Get("name")
		password := r.Form.Get("pass")

		// Basic credential check
		if username == "testuser" && password == "testpass" {
			// Set a cookie to simulate successful login
			http.SetCookie(w, &http.Cookie{Name: "sessionid", Value: "12345", Path: "/"})
			w.WriteHeader(http.StatusOK)
			return
		}

		// Invalid credentials
		w.WriteHeader(http.StatusUnauthorized)
	}))
}

func TestLogin(t *testing.T) {
	server := setupMockLoginServer()
	defer server.Close()

	// Create HTTP client with cookie jar
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	creds := Credentials{
		Username: "testuser",
		Password: "testpass",
	}

	// Perform login
	err := Login(client, server.URL+LoginPath, creds)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Check that cookie is stored in client
	u, _ := url.Parse(server.URL)
	cookies := client.Jar.Cookies(u)
	if len(cookies) == 0 {
		t.Fatal("Expected cookies to be set after login, but got none")
	}

	// Verify cookie name and value
	if cookies[0].Name != "sessionid" || cookies[0].Value != "12345" {
		t.Fatalf("Unexpected cookie: %+v", cookies[0])
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	server := setupMockLoginServer()
	defer server.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	creds := Credentials{
		Username: "wronguser",
		Password: "wrongpass",
	}

	err := Login(client, server.URL, creds)
	if err == nil {
		t.Fatal("Expected login to fail with invalid credentials, but it succeeded")
	}
}
