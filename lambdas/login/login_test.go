package main

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
)

// --- mock server for login ---
func setupMockLoginServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		username := r.Form.Get("name")
		password := r.Form.Get("pass")

		if username == "testuser" && password == "testpass" {
			http.SetCookie(w, &http.Cookie{Name: "sessionid", Value: "12345", Path: "/"})
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
	}))
}

func TestLogin_Success(t *testing.T) {
	server := setupMockLoginServer()
	defer server.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	creds := Credentials{
		Username: "testuser",
		Password: "testpass",
	}

	loginURL := server.URL
	cookies, err := Login(client, loginURL, creds)
	if err != nil {
		t.Fatalf("Expected login to succeed, but got error: %v", err)
	}

	if len(cookies) == 0 {
		t.Fatal("Expected cookies to be set after successful login, but got none")
	}

	if cookies[0].Name != "sessionid" || cookies[0].Value != "12345" {
		t.Fatalf("Unexpected cookie: %+v", cookies[0])
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	server := setupMockLoginServer()
	defer server.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	creds := Credentials{
		Username: "wronguser",
		Password: "wrongpass",
	}

	_, err := Login(client, server.URL, creds)
	if err == nil {
		t.Fatal("Expected error for invalid credentials, but got nil")
	}
}

func TestLogin_MissingUsername(t *testing.T) {
	server := setupMockLoginServer()
	defer server.Close()

	client := &http.Client{}

	creds := Credentials{
		Username: "",
		Password: "testpass",
	}

	_, err := Login(client, server.URL, creds)
	if err == nil || err.Error() != "ACCOUNT_USERNAME is required" {
		t.Fatalf("Expected missing username error, got: %v", err)
	}
}

func TestLogin_MissingPassword(t *testing.T) {
	server := setupMockLoginServer()
	defer server.Close()

	client := &http.Client{}

	creds := Credentials{
		Username: "testuser",
		Password: "",
	}

	_, err := Login(client, server.URL, creds)
	if err == nil || err.Error() != "ACCOUNT_PASSWORD is required" {
		t.Fatalf("Expected missing password error, got: %v", err)
	}
}

func TestExtractCookies_NoCookies(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	_, err := extractCookies(client, server.URL)
	if err == nil {
		t.Fatal("Expected error when no cookies are set, but got nil")
	}
}

func TestBuildLoginRequest(t *testing.T) {
	creds := Credentials{Username: "testuser", Password: "testpass"}
	req, err := buildLoginRequest("http://example.com", creds)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if req.Method != http.MethodPost {
		t.Fatalf("Expected POST method, got %s", req.Method)
	}

	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Fatalf("Expected content-type header, got %s", req.Header.Get("Content-Type"))
	}

	bodyBytes := make([]byte, req.ContentLength)
	req.Body.Read(bodyBytes)
	body := string(bodyBytes)
	expectedParams := url.Values{
		"form_id": {"user_login_form"},
		"name":    {"testuser"},
		"pass":    {"testpass"},
	}.Encode()

	if body != expectedParams {
		t.Fatalf("Unexpected form body: got %s, want %s", body, expectedParams)
	}
}
