package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/middleware"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/routes/mockresponses"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/utils"
	"github.com/gorilla/mux"
)

func RegisterUpcomingProjectsRoute(r *mux.Router) {
	r.HandleFunc("/api/registrations/dashboard/upcoming/{userId}/{page}", middleware.RequireCookieMiddleware(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["userId"]

		if !utils.ValidateInternalID(userId) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		projects, err := GetProjectsFromCookie(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		resp := mockresponses.MockUpcomingResponse(projects)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})).Methods("GET")
}

func RegisterTodayProjectsRoute(r *mux.Router) {
	r.HandleFunc("/api/registrations/dashboard/today/{userId}/{page}", middleware.RequireCookieMiddleware(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["userId"]

		if !utils.ValidateInternalID(userId) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		projects, err := GetProjectsFromCookie(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		resp := mockresponses.MockTodayResponse(projects)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})).Methods("GET")
}
