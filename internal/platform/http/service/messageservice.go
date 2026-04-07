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
	SendMessage(ctx context.Context, channelId, messageContent string) (string, error)
	PinMessage(ctx context.Context, campaignId, messageId string) error
	SetCookies(cookies []*http.Cookie) error
}

func (s *HttpService) SendMessage(ctx context.Context, channelId, messageContent string) (string, error) {
	req, err := s.buildSendMessageRequest(channelId, messageContent)
	if err != nil {
		return "", fmt.Errorf("failed to build send message request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return "", fmt.Errorf("send message request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return "", fmt.Errorf("send message request failed: %w", err)
	}

	body, err := s.ReadBody(resp)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var sendMessageResp dto.SendMessageResponse
	if err := json.Unmarshal(body, &sendMessageResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal send message response: %w", err)
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

	urlStr := endpoints.JoinPaths(s.baseUrl, "/api/messenger/channel/", channelId, "/messages/post")

	req, err := http.NewRequest("POST", urlStr, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	return req, nil
}

///////

func (s *HttpService) PinMessage(ctx context.Context, campaignId, messageId string) error {
	req, err := s.buildPinMessageRequest(campaignId, messageId)
	if err != nil {
		return fmt.Errorf("failed to build pin message request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("pin message request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return fmt.Errorf("pin message request failed: %w", err)
	}

	return nil
}

func (s *HttpService) buildPinMessageRequest(campaignId, messageId string) (*http.Request, error) {

	body := map[string]string{
		"MessageId": messageId,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	urlStr := endpoints.JoinPaths(s.baseUrl, "/api/messenger/create-pin-message/", campaignId)

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
