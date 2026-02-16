package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/middleware"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/routes/mockresponses"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/utils"
	"github.com/gorilla/mux"
)

func RegisterScheduleRoute(r *mux.Router) {
	r.HandleFunc("/api/schedule/retrieve/{internalId}", middleware.RequireCookieMiddleware(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		internalId := vars["internalId"]

		if !utils.ValidateInternalID(internalId) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		projects := GetAdminProjects()
		resp := mockresponses.MockScheduleResponse(projects)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})).Methods("GET")
}
