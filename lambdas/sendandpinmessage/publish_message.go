package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

func GetProjectDetails(projectId string) models.CampaignResponse {
	url := fmt.Sprintf("https://www.newyorkcares.org/api/campaign/retrieve/%s", projectId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Go-http-client/1.1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var campaign models.CampaignResponse
	if err := json.Unmarshal(body, &campaign); err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	return campaign
	// I can then rip out the AWSChimeChannelId for the below function
}

func SendMessage(channelId string) models.MessageResponse {
	url := fmt.Sprintf("https://www.newyorkcares.org/api/messenger/channel/%s/messages/post", channelId)

	// Prepare multipart form body
	var bodyBuffer bytes.Buffer
	writer := multipart.NewWriter(&bodyBuffer)
	message := "helloooooo"

	if err := writer.WriteField("message", message); err != nil {
		log.Fatalf("Failed to write message field: %v", err)
	}

	writer.Close()

	// Build request
	req, err := http.NewRequest("POST", url, &bodyBuffer)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Basic headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://www.newyorkcares.org")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// Add cookies if needed
	// req.Header.Set("Cookie", "SSESS901807e183ba6c96f257b0692ef12e9c=your_cookie_here")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var messageResp models.MessageResponse
	if err := json.Unmarshal(respBody, &messageResp); err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Optional: print response for debugging
	fmt.Printf("Response: %+v\n", messageResp)
	return messageResp
}

// func PinMessage(channelId string) models.PinResponse {

// }
