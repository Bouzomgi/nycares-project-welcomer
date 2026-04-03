package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/routes/mockresponses"
)

type projectInput struct {
	Name       string `json:"name"`
	Date       string `json:"date"`
	Id         string `json:"id"`
	CampaignId string `json:"campaignId"`
}

func GetProjectsFromCookie(r *http.Request) ([]mockresponses.ProjectConfig, error) {
	cookie, err := r.Cookie("session")
	if err != nil || !strings.HasPrefix(cookie.Value, "mock-session:") {
		return nil, nil
	}
	b64 := strings.TrimPrefix(cookie.Value, "mock-session:")
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("bad cookie: %w", err)
	}
	var stored []projectInput
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, fmt.Errorf("unmarshal cookie: %w", err)
	}
	projects := make([]mockresponses.ProjectConfig, 0, len(stored))
	for _, p := range stored {
		date, err := time.Parse("2006-01-02", p.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date %q: %w", p.Date, err)
		}
		projects = append(projects, mockresponses.ProjectConfig{
			Name:       p.Name,
			Date:       date,
			Id:         p.Id,
			CampaignId: p.CampaignId,
		})
	}
	return projects, nil
}
