package sendandpinmessage

import (
	"context"
	"strings"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/service"
	s3service "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/s3"
)

type SendAndPinMessageUseCase struct {
	s3Service   s3service.ContentService
	httpService httpservice.MessageService
}

func NewSendAndPinMessageUseCase(s3Service s3service.ContentService, httpService httpservice.MessageService) *SendAndPinMessageUseCase {
	return &SendAndPinMessageUseCase{
		s3Service,
		httpService,
	}
}

func (u *SendAndPinMessageUseCase) Execute(ctx context.Context, auth domain.Auth, projectId, channelId, messageRef, generatedContent, projectName string) error {

	u.httpService.SetCookies(auth.Cookies)

	var messageContent string
	if generatedContent != "" {
		messageContent = generatedContent
	} else {
		var err error
		messageContent, err = u.s3Service.GetMessageContent(ctx, messageRef)
		if err != nil {
			return err
		}
		messageContent = strings.ReplaceAll(messageContent, "{{projectName}}", projectName)
	}

	messageId, err := u.httpService.SendMessage(ctx, channelId, messageContent)
	if err != nil {
		return err
	}

	err = u.httpService.PinMessage(ctx, projectId, messageId)
	if err != nil {
		return err
	}

	return nil
}
