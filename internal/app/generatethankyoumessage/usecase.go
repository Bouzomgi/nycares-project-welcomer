package generatethankyoumessage

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	bedrockservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/bedrock"
	s3service "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/s3"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)

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

func (u *GenerateThankYouMessageUseCase) Execute(ctx context.Context, projectName, refinementContext string) (string, error) {
	samplesRef := fmt.Sprintf("s3://%s/%s/thankyou-samples.md", u.bucketName, toKebabCase(projectName))

	writingSample, err := u.s3Service.GetMessageContent(ctx, samplesRef)
	if err != nil {
		return "", fmt.Errorf("failed to fetch writing samples: %w", err)
	}

	return u.bedrockService.GenerateThankYouMessage(ctx, writingSample, refinementContext)
}

func toKebabCase(s string) string {
	parts := nonAlphanumeric.Split(s, -1)
	words := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			words = append(words, strings.ToLower(p))
		}
	}
	return strings.Join(words, "-")
}
