package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterLoginRoute(r *mux.Router) {
	r.HandleFunc("/user/login", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid form")
			return
		}

		formId := r.FormValue("form_id")
		username := r.FormValue("name")
		password := r.FormValue("pass")

		if username == "" || password == "" || formId != "user_login_form" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "form_id, username, password required")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "mock-session-id",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   3600, // 1 hour
		})

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `<!DOCTYPE html><html><body><div data-endpoint="api/schedule/retrieve/003MOCK00000000001"></div></body></html>`)
	}).Methods("POST")
}
