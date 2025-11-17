package httpservice

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

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

	mockHttpService, err := NewHttpService(server.URL)
	if err != nil {
		t.Fatalf("Could not create mock server: %v", err)
	}

	creds := models.Credentials{
		Username: "testuser",
		Password: "testpass",
	}

	ctx := context.Background()
	auth, err := mockHttpService.Login(ctx, creds)
	if err != nil {
		t.Fatalf("Expect login to suceed, but got error: %v", err)
	}

	if len(auth.Cookies) == 0 {
		t.Fatal("Expect cookies to set after successful login, but got none")
	}

	if auth.Cookies[0].Name != "sessionid" || auth.Cookies[0].Value != "12345" {
		t.Fatalf("Unexpected cookie: %+v", auth.Cookies[0])
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	server := setupMockLoginServer()
	defer server.Close()

	mockHttpService, err := NewHttpService(server.URL)
	if err != nil {
		t.Fatalf("Could not create mock server: %v", err)
	}

	creds := models.Credentials{
		Username: "wronguser",
		Password: "wrongpass",
	}

	ctx := context.Background()
	_, loginErr := mockHttpService.Login(ctx, creds)
	if loginErr == nil {
		t.Fatal("Expected error for invalid credential, but got nil")
	}
}

func TestLogin_MissingUsername(t *testing.T) {
	server := setupMockLoginServer()
	defer server.Close()

	mockHttpService, err := NewHttpService(server.URL)
	if err != nil {
		t.Fatalf("Could not create mock server: %v", err)
	}

	creds := models.Credentials{
		Username: "",
		Password: "testpass",
	}

	ctx := context.Background()
	_, loginErr := mockHttpService.Login(ctx, creds)
	if loginErr == nil || loginErr.Error() != "username is required" {
		t.Fatalf("Expected missing username error, got %v", loginErr)
	}
}

func TestLogin_MissingPassword(t *testing.T) {
	server := setupMockLoginServer()
	defer server.Close()

	mockHttpService, err := NewHttpService(server.URL)
	if err != nil {
		t.Fatalf("Could not create mock server: %v", err)
	}

	creds := models.Credentials{
		Username: "testuser",
		Password: "",
	}

	ctx := context.Background()
	_, loginErr := mockHttpService.Login(ctx, creds)
	if loginErr == nil || loginErr.Error() != "password is required" {
		t.Fatalf("Expected missing password error, got %v", loginErr)
	}
}

func TestBuildLoginPassword(t *testing.T) {
	server := setupMockLoginServer()
	defer server.Close()

	mockHttpService, err := NewHttpService(server.URL)
	if err != nil {
		t.Fatalf("Could not create mock server: %v", err)
	}

	creds := models.Credentials{Username: "testuser", Password: "testpass"}
	req, err := mockHttpService.buildLoginRequest(creds)
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
