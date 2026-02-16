package approvalcallback

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type SFNClient interface {
	SendTaskSuccess(ctx context.Context, params *sfn.SendTaskSuccessInput, optFns ...func(*sfn.Options)) (*sfn.SendTaskSuccessOutput, error)
	SendTaskFailure(ctx context.Context, params *sfn.SendTaskFailureInput, optFns ...func(*sfn.Options)) (*sfn.SendTaskFailureOutput, error)
}

type ApprovalCallbackUseCase struct {
	sfnClient SFNClient
}

func NewApprovalCallbackUseCase(sfnClient SFNClient) *ApprovalCallbackUseCase {
	return &ApprovalCallbackUseCase{sfnClient: sfnClient}
}

func (u *ApprovalCallbackUseCase) Execute(ctx context.Context, taskToken string, approved bool) error {
	if taskToken == "" {
		return fmt.Errorf("taskToken must be defined")
	}

	if approved {
		_, err := u.sfnClient.SendTaskSuccess(ctx, &sfn.SendTaskSuccessInput{
			TaskToken: &taskToken,
			Output:    strPtr(`{"approved": true}`),
		})
		return err
	}

	_, err := u.sfnClient.SendTaskFailure(ctx, &sfn.SendTaskFailureInput{
		TaskToken: &taskToken,
		Error:     strPtr("rejected"),
		Cause:     strPtr("User rejected the approval request"),
	})
	return err
}

func strPtr(s string) *string {
	return &s
}
