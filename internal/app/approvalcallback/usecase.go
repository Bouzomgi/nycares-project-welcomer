package approvalcallback

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type SFNClient interface {
	SendTaskSuccess(ctx context.Context, params *sfn.SendTaskSuccessInput, optFns ...func(*sfn.Options)) (*sfn.SendTaskSuccessOutput, error)
}

type ApprovalCallbackUseCase struct {
	sfnClient SFNClient
}

func NewApprovalCallbackUseCase(sfnClient SFNClient) *ApprovalCallbackUseCase {
	return &ApprovalCallbackUseCase{sfnClient: sfnClient}
}

// Execute resolves the task token based on the action:
//   - "approve"     → SendTaskSuccess with action="approve"
//   - "regenerate"  → SendTaskSuccess with action="regenerate" (fresh generation, no context)
//   - "refine"      → SendTaskSuccess with action="refine" and the supplied refinementContext
//   - "reject"      → SendTaskSuccess with action="reject" (routes to EndProjectIteration, not DLQ)
//   - "manual"      → SendTaskSuccess with action="manual" and the supplied manualMessage
func (u *ApprovalCallbackUseCase) Execute(ctx context.Context, taskToken, action, refinementContext, manualMessage string) error {
	if taskToken == "" {
		return fmt.Errorf("taskToken must be defined")
	}

	switch action {
	case "approve", "regenerate", "refine", "reject", "manual":
		type successOutput struct {
			Action            string `json:"action"`
			RefinementContext string `json:"refinementContext"`
			ManualMessage     string `json:"manualMessage"`
		}
		out, err := json.Marshal(successOutput{
			Action:            action,
			RefinementContext: refinementContext,
			ManualMessage:     manualMessage,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal task output: %w", err)
		}
		outStr := string(out)
		_, err = u.sfnClient.SendTaskSuccess(ctx, &sfn.SendTaskSuccessInput{
			TaskToken: &taskToken,
			Output:    &outStr,
		})
		return err

	default:
		return fmt.Errorf("unknown action %q: must be approve, regenerate, refine, reject, or manual", action)
	}
}
