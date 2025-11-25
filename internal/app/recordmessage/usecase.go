package recordmessage

import (
	"context"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	dynamoservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/dynamo/service"
)

type RecordMessageUseCase struct {
	dynamoSvc dynamoservice.StoredNotificationService
}

func NewRecordMessageUseCase(dynamoSvc dynamoservice.StoredNotificationService) *RecordMessageUseCase {
	return &RecordMessageUseCase{dynamoSvc}
}

func (u *RecordMessageUseCase) Execute(ctx context.Context, existingProjectNotification domain.ProjectNotification, sentMessageType domain.NotificationType) (domain.ProjectNotification, error) {

	updatedHasSentWelcome := existingProjectNotification.HasSentWelcome || (sentMessageType == domain.Welcome)
	updatedHasSentReminder := existingProjectNotification.HasSentReminder || (sentMessageType == domain.Reminder)

	updatedProjectNotification := &domain.ProjectNotification{
		Name:             existingProjectNotification.Name,
		Date:             existingProjectNotification.Date,
		Id:               existingProjectNotification.Id,
		HasSentWelcome:   updatedHasSentWelcome,
		HasSentReminder:  updatedHasSentReminder,
		ShouldStopNotify: existingProjectNotification.ShouldStopNotify,
	}

	writtenProjectNotification, err := u.dynamoSvc.UpsertProjectNotification(ctx, *updatedProjectNotification)
	if err != nil {
		return domain.ProjectNotification{}, err
	}

	return *writtenProjectNotification, nil
}
