package sendandpinmessage

import (
	"context"
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http"
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

func (u *SendAndPinMessageUseCase) Execute(ctx context.Context, auth domain.Auth, projectId, messageRef string) error {

	u.httpService.SetCookies(auth.Cookies)

	messageContent, err := u.s3Service.GetMessageContent(ctx, messageRef)
	if err != nil {
		return err
	}

	channelId, err := u.httpService.GetProjectChannelId(ctx, projectId)
	if err != nil {
		return err
	}

	fmt.Printf("channelId: %s\n\n", channelId)

	fmt.Printf("messageContent: %s\n\n", messageContent)

	println("We are NOT sending and pinning yet!!")

	// messageId, err := u.httpService.SendMessage(ctx, channelId, messageContent)
	// if err != nil {
	// 	return err
	// }

	// err = u.httpService.PinMessage(ctx, channelId, messageId)
	// if err != nil {
	// 	return err
	// }

	return nil
}
