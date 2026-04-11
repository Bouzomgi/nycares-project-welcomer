package bedrockservice

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

const ModelID = "anthropic.claude-3-5-haiku-20241022"

type GenerationService interface {
	GenerateThankYouMessage(ctx context.Context, writingSample, projectName string) (string, error)
}

type BedrockService struct {
	client *bedrockruntime.Client
}

func NewBedrockService(client *bedrockruntime.Client) *BedrockService {
	return &BedrockService{client: client}
}

func (s *BedrockService) GenerateThankYouMessage(ctx context.Context, writingSample, projectName string) (string, error) {
	systemPrompt := fmt.Sprintf("You are writing thank-you messages on behalf of a volunteer program coordinator. Here are several example messages they have written — match their style exactly:\n\n%s", writingSample)
	userPrompt := fmt.Sprintf("Write a new, unique thank-you message (2-3 sentences) for a team leader who led the volunteer project \"%s\" today. Match the style of the examples but do not repeat any of them.", projectName)

	resp, err := s.client.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId: aws.String(ModelID),
		System: []types.SystemContentBlock{
			&types.SystemContentBlockMemberText{Value: systemPrompt},
		},
		Messages: []types.Message{
			{
				Role: types.ConversationRoleUser,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{Value: userPrompt},
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
