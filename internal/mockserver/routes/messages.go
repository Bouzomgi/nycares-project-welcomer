package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/middleware"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/routes/mockresponses"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/utils"
	"github.com/gorilla/mux"
)

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
		_, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "cannot read body")
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid message body:", err)
			return
		}

		resp := mockresponses.MockPostPinResponse()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	})).Methods("POST")
}
