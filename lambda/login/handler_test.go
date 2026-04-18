package main

import (
	"net/http"
	"testing"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
)

func TestToResponseAuth(t *testing.T) {
	tests := []struct {
		name        string
		auth        domain.Auth
		executionId string
		wantCookies int
		wantIntId   string
		wantExecId  string
	}{
		{
			name: "single cookie is converted correctly",
			auth: domain.Auth{
				Cookies:    []*http.Cookie{{Name: "session", Value: "abc123"}},
				InternalId: "user-99",
			},
			executionId: "exec-1",
			wantCookies: 1,
			wantIntId:   "user-99",
			wantExecId:  "exec-1",
		},
		{
			name: "multiple cookies are all converted",
			auth: domain.Auth{
				Cookies: []*http.Cookie{
					{Name: "a", Value: "1"},
					{Name: "b", Value: "2"},
					{Name: "c", Value: "3"},
				},
				InternalId: "user-42",
			},
			executionId: "exec-2",
			wantCookies: 3,
			wantIntId:   "user-42",
			wantExecId:  "exec-2",
		},
		{
			name: "nil cookies slice produces empty cookies slice",
			auth: domain.Auth{
				Cookies:    nil,
				InternalId: "user-0",
			},
			executionId: "exec-3",
			wantCookies: 0,
			wantIntId:   "user-0",
			wantExecId:  "exec-3",
		},
		{
			name: "empty cookies slice produces empty cookies slice",
			auth: domain.Auth{
				Cookies:    []*http.Cookie{},
				InternalId: "user-1",
			},
			executionId: "exec-4",
			wantCookies: 0,
			wantIntId:   "user-1",
			wantExecId:  "exec-4",
		},
		{
			name: "cookie fields are preserved",
			auth: domain.Auth{
				Cookies: []*http.Cookie{
					{
						Name:     "token",
						Value:    "xyz",
						Domain:   "example.com",
						Path:     "/",
						HttpOnly: true,
						Secure:   true,
					},
				},
				InternalId: "user-5",
			},
			executionId: "exec-5",
			wantCookies: 1,
			wantIntId:   "user-5",
			wantExecId:  "exec-5",
		},
		{
			name: "empty executionId is forwarded",
			auth: domain.Auth{
				Cookies:    []*http.Cookie{{Name: "s", Value: "v"}},
				InternalId: "user-6",
			},
			executionId: "",
			wantCookies: 1,
			wantIntId:   "user-6",
			wantExecId:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ToResponseAuth(tt.auth, tt.executionId)

			if len(output.Auth.Cookies) != tt.wantCookies {
				t.Errorf("len(Cookies) = %d, want %d", len(output.Auth.Cookies), tt.wantCookies)
			}
			if output.InternalId != tt.wantIntId {
				t.Errorf("InternalId = %q, want %q", output.InternalId, tt.wantIntId)
			}
			if output.ExecutionId != tt.wantExecId {
				t.Errorf("ExecutionId = %q, want %q", output.ExecutionId, tt.wantExecId)
			}

			// Verify cookie field preservation for the first cookie when present.
			if tt.wantCookies > 0 && tt.auth.Cookies[0] != nil {
				orig := tt.auth.Cookies[0]
				got := output.Auth.Cookies[0]
				if got.Name != orig.Name {
					t.Errorf("Cookies[0].Name = %q, want %q", got.Name, orig.Name)
				}
				if got.Value != orig.Value {
					t.Errorf("Cookies[0].Value = %q, want %q", got.Value, orig.Value)
				}
			}
		})
	}
}
