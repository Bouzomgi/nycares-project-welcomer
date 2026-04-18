package bedrockservice

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

const ModelID = "us.amazon.nova-lite-v1:0"

type GenerationService interface {
	GenerateThankYouMessage(ctx context.Context, writingSample, refinementContext string) (string, error)
}

type BedrockService struct {
	client *bedrockruntime.Client
}

func NewBedrockService(client *bedrockruntime.Client) *BedrockService {
	return &BedrockService{client: client}
}

func (s *BedrockService) GenerateThankYouMessage(ctx context.Context, writingSample, refinementContext string) (string, error) {
	prompt := fmt.Sprintf(
		"You are a volunteer project team leader writing thank-you messages to your volunteers. Here are several example messages you have written — match their style exactly:\n\n%s\n\nWrite a new, unique thank-you message (2-3 sentences) to the volunteers who helped with today's project. Match the style of the examples but do not repeat any of them.",
		writingSample,
	)
	if refinementContext != "" {
		prompt += fmt.Sprintf(" Additional context to incorporate: %s", refinementContext)
	}

	resp, err := s.client.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId: aws.String(ModelID),
		Messages: []types.Message{
			{
				Role: types.ConversationRoleUser,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{Value: prompt},
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("bedrock Converse failed: %w", err)
	}

	output, ok := resp.Output.(*types.ConverseOutputMemberMessage)
	if !ok || len(output.Value.Content) == 0 {
		return "", fmt.Errorf("bedrock returned unexpected output shape")
	}

	textBlock, ok := output.Value.Content[0].(*types.ContentBlockMemberText)
	if !ok {
		return "", fmt.Errorf("bedrock first content block is not text")
	}

	return textBlock.Value, nil
}
