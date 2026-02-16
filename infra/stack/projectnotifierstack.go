package stack

import (
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/iancoleman/strcase"
)

type LambdaStackProps struct {
	awscdk.StackProps
}

func ProjectNotifierStack(scope constructs.Construct, id string, props *LambdaStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	lambdas := []string{
		"Login",
		"FetchProjects",
		"ComputeMessageToSend",
		"RequestApprovalToSend",
		"SendAndPinMessage",
		"RecordMessage",
		"NotifyCompletion",
		"DLQNotifier",
	}

	for _, name := range lambdas {
		lowerName := strings.ToLower(name)
		kebabName := strcase.ToKebab(name)

		awslambda.NewFunction(stack, jsii.String(name), &awslambda.FunctionProps{
			Runtime: awslambda.Runtime_PROVIDED_AL2023(),
			Handler: jsii.String("bootstrap"),
			Code: awslambda.Code_FromAsset(
				jsii.String("../lambda-build/"+lowerName),
				nil,
			),
			FunctionName: jsii.String(kebabName),
			Architecture: awslambda.Architecture_ARM_64(),
		})
	}

	awsstepfunctions.NewStateMachine(stack, jsii.String("ProjectNotifierStateMachine"), &awsstepfunctions.StateMachineProps{
		StateMachineName: jsii.String("project-notifier-workflow"),
		DefinitionBody:   awsstepfunctions.DefinitionBody_FromFile((jsii.String("DailyProjectNotificationWorkflow.json")), nil),
	})

	return stack
}
