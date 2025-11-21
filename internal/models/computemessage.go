package models

type ComputeMessageInput struct {
	Auth    Auth    `json:"auth"`
	Project project `json:"project"`
}

type ComputeMessageOutput struct {
	Auth          Auth    `json:"auth"`
	Project       project `json:"project"`
	MessageToSend message `json:"message"`
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
