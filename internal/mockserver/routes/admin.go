package routes

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/routes/mockresponses"
	"github.com/gorilla/mux"
)

var (
	adminMu       sync.RWMutex
	adminProjects []mockresponses.ProjectConfig
)

func GetAdminProjects() []mockresponses.ProjectConfig {
	adminMu.RLock()
	defer adminMu.RUnlock()
	return adminProjects
}

func SetAdminProjects(projects []mockresponses.ProjectConfig) {
	adminMu.Lock()
	defer adminMu.Unlock()
	adminProjects = projects
}

type setProjectsRequest struct {
	Projects []projectInput `json:"projects"`
}

type projectInput struct {
	Name       string `json:"name"`
	Date       string `json:"date"`
	Id         string `json:"id"`
	CampaignId string `json:"campaignId"`
}

func RegisterAdminRoutes(r *mux.Router) {
	r.HandleFunc("/admin/set-projects", func(w http.ResponseWriter, r *http.Request) {
		var req setProjectsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		var projects []mockresponses.ProjectConfig
		for _, p := range req.Projects {
			date, err := time.Parse("2006-01-02", p.Date)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid date: " + p.Date})
				return
			}
			projects = append(projects, mockresponses.ProjectConfig{
				Name:       p.Name,
				Date:       date,
				Id:         p.Id,
				CampaignId: p.CampaignId,
			})
		}

		SetAdminProjects(projects)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("POST")
}
