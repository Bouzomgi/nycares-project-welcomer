package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/httphelper"
)

func buildScheduleRequest(baseUrl, internalId string) (*http.Request, error) {
	getScheduleBaseUrl := baseUrl + endpoints.GetSchedulePath
	urlStr := fmt.Sprintf("%s/%s", getScheduleBaseUrl, internalId)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	return req, nil
}

// flattenScheduleList converts a ScheduleResponse into a slice of Projects
func flattenScheduleList(resp ScheduleResponse) []CompleteProject {
	var result []CompleteProject

	for _, schedule := range resp.Data.ScheduleList {
		result = append(result, schedule)
	}

	return result
}

func GetSchedule(client *http.Client, baseUrl, internalId string) ([]CompleteProject, error) {
	if internalId == "" {
		return nil, fmt.Errorf("internalId is required")
	}

	req, err := buildScheduleRequest(baseUrl, internalId)
	if err != nil {
		return nil, err
	}

	resp, err := httphelper.SendRequest(client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	body, err := httphelper.ReadBody(resp)
	if err != nil {
		return nil, err
	}

	var respArray []ScheduleResponse
	if err := json.Unmarshal(body, &respArray); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedule response: %w", err)
	}

	schedule := respArray[0]

	scheduleList := flattenScheduleList(schedule)

	return scheduleList, nil
}
