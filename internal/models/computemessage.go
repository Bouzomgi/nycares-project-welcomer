package models

import (
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/utils"
)

type ComputeMessageInput struct {
	Auth    Auth    `json:"auth"`
	Project project `json:"project"`
}

type ComputeMessageOutput struct {
	Auth                        Auth                `json:"auth"`
	ExistingProjectNotification projectNotification `json:"existingProjectNotification"`
	MessageToSend               message             `json:"message"`
}

type projectNotification struct {
	Name             string `json:"projectName"`
	Date             string `json:"projectDate"`
	Id               string `json:"projectId"`
	HasSentWelcome   bool   `json:"hasSentWelcome"`
	HasSentReminder  bool   `json:"hasSentReminder"`
	ShouldStopNotify bool   `json:"shouldStopNotify"`
}

type message struct {
	Type        string `json:"type"`
	TemplateRef string `json:"templateRef"`
}

func BuildMessage(messageType, templateRef string) message {
	return message{
		Type:        messageType,
		TemplateRef: templateRef,
	}
}

func ConvertDomainProjectNotification(pn domain.ProjectNotification) projectNotification {
	return projectNotification{
		Name:             pn.Name,
		Date:             pn.Date,
		Id:               pn.Id,
		HasSentWelcome:   pn.HasSentWelcome,
		HasSentReminder:  pn.HasSentReminder,
		ShouldStopNotify: pn.ShouldStopNotify,
	}
}

func ConvertModelProjectNotification(pn projectNotification) domain.ProjectNotification {
	return domain.ProjectNotification{
		Name:             pn.Name,
		Date:             pn.Date,
		Id:               pn.Id,
		HasSentWelcome:   pn.HasSentWelcome,
		HasSentReminder:  pn.HasSentReminder,
		ShouldStopNotify: pn.ShouldStopNotify,
	}
}

func ConvertProjectNotificationToDomainProject(pn projectNotification) (domain.Project, error) {
	projectDate, err := utils.StringToDate(pn.Date)

	if err != nil {
		return domain.Project{}, fmt.Errorf("could not parse project date")
	}

	project := domain.Project{
		Name: pn.Name,
		Date: projectDate,
		Id:   pn.Id,
	}

	return project, nil
}
