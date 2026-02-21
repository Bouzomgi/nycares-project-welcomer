package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type LoginResponse struct {
	Message string `json:"message"`
}

func RegisterLoginRoute(r *mux.Router) {
	r.HandleFunc("/user/login", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := LoginResponse{
			Message: "Login successful!",
		}

		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")
}
