package generatethankyoumessage

import (
	"context"
	"fmt"
	"strings"

	bedrockservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/bedrock"
	s3service "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/s3"
)

type GenerateThankYouMessageUseCase struct {
	s3Service      s3service.ContentService
	bedrockService bedrockservice.GenerationService
	bucketName     string
}

func NewGenerateThankYouMessageUseCase(s3Svc s3service.ContentService, bedrockSvc bedrockservice.GenerationService, bucketName string) *GenerateThankYouMessageUseCase {
	return &GenerateThankYouMessageUseCase{
		s3Service:      s3Svc,
		bedrockService: bedrockSvc,
		bucketName:     bucketName,
	}
}

func (u *GenerateThankYouMessageUseCase) Execute(ctx context.Context, projectName string) (string, error) {
	samplesRef := fmt.Sprintf("s3://%s/%s/thankyou-samples.md", u.bucketName, toKebabCase(projectName))

	writingSample, err := u.s3Service.GetMessageContent(ctx, samplesRef)
	if err != nil {
		return "", fmt.Errorf("failed to fetch writing samples: %w", err)
	}

	return u.bedrockService.GenerateThankYouMessage(ctx, writingSample, projectName)
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
