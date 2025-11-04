package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	// Mock server to simulate login
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the method is POST
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// Check form values
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.FormValue("username") != "testuser" {
			t.Errorf("expected username=testuser, got %s", r.FormValue("username"))
		}
		if r.FormValue("password") != "testpass" {
			t.Errorf("expected password=testpass, got %s", r.FormValue("password"))
		}

		// Set a cookie to simulate successful login
		http.SetCookie(w, &http.Cookie{Name: "sessionid", Value: "12345"})
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := &http.Client{}

	creds := Credentials {
		Username: "testuser",
		Password: "testpass",
	}

	cookies, err := Login(client, ts.URL, creds)
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	if cookies["sessionid"] != "1234j" {
		t.Errorf("expected cookie sessionid=12345, got %v", cookies)
	}
}
