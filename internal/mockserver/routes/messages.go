package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/middleware"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/routes/mockresponses"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/utils"
	"github.com/gorilla/mux"
)

var lastSentMessage string

// Master function to register all MessageService mock routes
func RegisterMessageRoutes(r *mux.Router) {
	registerCampaignRoute(r)
	registerChannelMessagesRoute(r)
	registerSendMessageRoute(r)
	registerPinMessageRoute(r)
}

// --- GET /campaign/{projectId} ---
func registerCampaignRoute(r *mux.Router) {
	r.HandleFunc("/api/campaign/retrieve/{campaignId}", middleware.RequireCookieMiddleware(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		campaignId := vars["campaignId"]

		if !utils.ValidateInternalID(campaignId) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := mockresponses.MockCampaignResponse()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})).Methods("GET")
}

// --- GET /messenger/channel/{channelId}/messages ---
func registerChannelMessagesRoute(r *mux.Router) {
	r.HandleFunc("/api/messenger/channel/{channelId}/messages", middleware.RequireCookieMiddleware(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId := vars["channelId"]

		if !utils.ValidateUUID(channelId) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid channelId")
			return
		}

		resp := mockresponses.MockChannelMessagesResponse()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})).Methods("GET")
}

// --- POST /api/messenger/channel/{channelId}/message/post ---
func registerSendMessageRoute(r *mux.Router) {
	r.HandleFunc("/api/messenger/channel/{channelId}/messages/post", middleware.RequireCookieMiddleware(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId := vars["channelId"]

		if !utils.ValidateUUID(channelId) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		message := r.FormValue("message")

		if message == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "message required")
			return
		}

		lastSentMessage = message

		resp := mockresponses.MockSendMessageResponse()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	})).Methods("POST")
}

// --- POST /api/messenger/create-pin-message/{campaignId} ---
func registerPinMessageRoute(r *mux.Router) {
	r.HandleFunc("/api/messenger/create-pin-message/{campaignId}", middleware.RequireCookieMiddleware(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		campaignId := vars["campaignId"]

		if !utils.ValidateCampaignID(campaignId) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid campaign id")

			return
		}

		defer r.Body.Close()
		if _, err := io.ReadAll(r.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "cannot read body")
			return
		}

		msgType := "reminder"
		if strings.Contains(strings.ToLower(lastSentMessage), "welcome") {
			msgType = "welcome"
		}

		projectName := campaignId
		projectDate := ""
		for _, p := range GetAdminProjects() {
			if p.Id == campaignId {
				projectName = p.Name
				projectDate = p.Date.Format("2006-01-02")
				break
			}
		}

		fmt.Printf("[mock] published and pinned: project=%q date=%s type=%s\n", projectName, projectDate, msgType)

		resp := mockresponses.MockPostPinResponse()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	})).Methods("POST")
}
