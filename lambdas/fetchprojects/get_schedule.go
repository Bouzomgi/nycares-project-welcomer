package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/httphelper"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

func buildScheduleRequest(baseUrl, internalId string) (*http.Request, error) {
	urlStr := fmt.Sprintf("%s/%s", baseUrl+GetSchedulePath, internalId)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	return req, nil
}

// flattenScheduleList converts a ScheduleResponse into a slice of Projects
func flattenScheduleList(resp models.ScheduleResponse) []models.Project {
	var result []models.Project

	for _, schedule := range resp.Data.ScheduleList {
		result = append(result, schedule)
	}

	return result
}

func GetSchedule(client *http.Client, baseUrl, internalId string) ([]models.Project, error) {
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

	var respArray []models.ScheduleResponse
	if err := json.Unmarshal(body, &respArray); err != nil {
		log.Fatal(err)
	}
	schedule := respArray[0]

	scheduleList := flattenScheduleList(schedule)

	return scheduleList, nil
}
