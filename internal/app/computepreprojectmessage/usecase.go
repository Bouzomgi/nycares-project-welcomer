package computepreprojectmessage

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
)

type ComputeMessageUseCase struct{}

func NewComputeMessageUseCase() *ComputeMessageUseCase {
	return &ComputeMessageUseCase{}
}

func (u *ComputeMessageUseCase) Execute(messageBucketName, projectName string, messageType domain.NotificationType) (string, error) {
	return computeS3MessageRefPath(messageBucketName, projectName, messageType)
}

func computeS3MessageRefPath(s3BucketName, projectName string, messageType domain.NotificationType) (string, error) {
	messageTypeStr := strings.ToLower(messageType.String())
	s3MessageRef := fmt.Sprintf("s3://%s/%s/%s.md", s3BucketName, toKebabCase(projectName), messageTypeStr)

	if isValidS3URI(s3MessageRef) {
		return s3MessageRef, nil
	}

	return "", fmt.Errorf("could not compute valid s3 bucket reference for message")
}

func toKebabCase(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	for i := range words {
		words[i] = strings.ToLower(words[i])
	}

	return strings.Join(words, "-")
}

var basicS3URIRegex = regexp.MustCompile(`^s3://[a-z0-9\-]{3,63}/.+$`)

func isValidS3URI(uri string) bool {
	return basicS3URIRegex.MatchString(uri)
}
