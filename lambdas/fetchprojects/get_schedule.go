package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nycaresprojectwelcomer/internal/httphelper"
	"nycaresprojectwelcomer/internal/models"
	"strings"
)

// buildCookieHeader converts a map of cookies to a single "Cookie" header string.
func buildCookieHeader(cookies map[string]string) (string, error) {
	if len(cookies) == 0 {
		return "", fmt.Errorf("no cookies provided")
	}

	var cookieHeader strings.Builder
	for k, v := range cookies {
		cookieHeader.WriteString(fmt.Sprintf("%s=%s; ", k, v))
	}

	// remove trailing "; "
	header := cookieHeader.String()
	return header[:len(header)-2], nil
}

func buildScheduleRequest(internalId string, cookies map[string]string) (*http.Request, error) {
	urlStr := fmt.Sprintf("%s/%s", GetScheduleUrl, internalId)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	cookieHeader, err := buildCookieHeader(cookies)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", cookieHeader)
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

func GetSchedule(client *http.Client, internalId string,  cookies map[string]string) ([]models.Project, error) {
	if internalId == "" {
		return nil, fmt.Errorf("internalId is required")
	}

	req, err := buildScheduleRequest(internalId, cookies)
	if err != nil {
		return nil, err
	}

	resp, err := httphelper.SendRequest(client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // ensure body is always closed

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	body, err := httphelper.ReadBody(resp)
	if err != nil {
		return nil, err
	}

	var schedule models.ScheduleResponse
	if err := json.Unmarshal(body, &schedule); err != nil {
		return nil, err
	}

	scheduleList := flattenScheduleList(schedule)

	return scheduleList, nil
}