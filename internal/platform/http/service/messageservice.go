package httpservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/dto"
)

type MessageService interface {
	GetProjectChannelId(ctx context.Context, projectId string) (string, error)
	SendMessage(ctx context.Context, channelId, messageContent string) (string, error)
	PinMessage(ctx context.Context, channelId, messageId string) error
	SetCookies(cookies []*http.Cookie) error
}

func (s *HttpService) GetProjectChannelId(ctx context.Context, projectId string) (string, error) {

	req, err := s.buildCampaignRequest(projectId)
	if err != nil {
		return "", fmt.Errorf("failed to build schedule request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return "", fmt.Errorf("schedule request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return "", fmt.Errorf("schedule request failed: %w", err)
	}

	body, err := s.ReadBody(resp)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var campaignResp dto.CampaignResponse
	if err := json.Unmarshal(body, &campaignResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal campaign response: %w", err)
	}

	if len(campaignResp) == 0 {
		return "", fmt.Errorf("campaign response was empty: %w", err)
	}

	channelId := campaignResp[0].Campaign.AWSChimeChannelID

	return channelId, nil
}

func (s *HttpService) buildCampaignRequest(projectId string) (*http.Request, error) {
	getCampaignBaseUrl := endpoints.JoinPaths(endpoints.BaseUrl, endpoints.GetCampaignPath)
	urlStr := fmt.Sprintf("%s/%s", getCampaignBaseUrl, projectId)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

///////

// TODO fix all these endpoints
func (s *HttpService) GetFirstMessageId(ctx context.Context, channelId string) (string, error) {

	req, err := s.buildMessagesRequest(channelId)
	if err != nil {
		return "", fmt.Errorf("failed to build schedule request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return "", fmt.Errorf("schedule request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return "", fmt.Errorf("schedule request failed: %w", err)
	}

	body, err := s.ReadBody(resp)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var messageResp dto.ChannelMessagesResponse
	if err := json.Unmarshal(body, &messageResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal campaign response: %w", err)
	}

	if len(messageResp) == 0 {
		return "", fmt.Errorf("messages response was empty: %w", err)
	}

	if len(messageResp[0].ChannelMessages) == 0 {
		return "", fmt.Errorf("no channel messages in first element")
	}

	return messageResp[0].ChannelMessages[0].MessageId, nil
}

func (s *HttpService) buildMessagesRequest(channelId string) (*http.Request, error) {
	urlStr := endpoints.JoinPaths(endpoints.BaseUrl, "/api/messenger/channel/", channelId, "messages")

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

///////

func (s *HttpService) SendMessage(ctx context.Context, channelId, messageContent string) (string, error) {
	req, err := s.buildSendMessageRequest(channelId, messageContent)
	if err != nil {
		return "", fmt.Errorf("failed to build schedule request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return "", fmt.Errorf("schedule request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return "", fmt.Errorf("schedule request failed: %w", err)
	}

	body, err := s.ReadBody(resp)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var sendMessageResp dto.SendMessageResponse
	if err := json.Unmarshal(body, &sendMessageResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal campaign response: %w", err)
	}

	messageId := sendMessageResp.Data.MessageId
	return messageId, nil
}

func (s *HttpService) buildSendMessageRequest(channelId, messageContent string) (*http.Request, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormField("message")
	if err != nil {
		return nil, err
	}

	io.WriteString(part, messageContent)
	writer.Close()

	urlStr := endpoints.JoinPaths(endpoints.BaseUrl, "/api/messenger/channel/", channelId, "/message/post")

	req, err := http.NewRequest("POST", urlStr, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	return req, nil
}

///////

func (s *HttpService) PinMessage(ctx context.Context, channelId, messageId string) error {
	req, err := s.buildPinMessageRequest(channelId, messageId)
	if err != nil {
		return fmt.Errorf("failed to build schedule request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("schedule request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return fmt.Errorf("schedule request failed: %w", err)
	}

	return nil
}

func (s *HttpService) buildPinMessageRequest(channelId, messageId string) (*http.Request, error) {

	body := map[string]string{
		"MessageId": messageId,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	urlStr := endpoints.JoinPaths(endpoints.BaseUrl, "/api/messenger/create-pin-message/", channelId)

	req, err := http.NewRequest(
		"POST",
		urlStr,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	return req, nil
}
