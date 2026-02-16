package stack

import (
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
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

	// Create ApprovalCallback lambda (outside loop — needs API Gateway wiring)
	approvalCallbackFn := awslambda.NewFunction(stack, jsii.String("ApprovalCallback"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("../lambda-build/approvalcallback"), nil),
		FunctionName: jsii.String("approval-callback"),
		Architecture: awslambda.Architecture_ARM_64(),
	})

	// Grant SFN permissions to the callback lambda
	approvalCallbackFn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("states:SendTaskSuccess", "states:SendTaskFailure"),
		Resources: jsii.Strings("*"),
	}))

	// Create REST API with GET /callback route
	api := awsapigateway.NewRestApi(stack, jsii.String("ApprovalCallbackApi"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("approval-callback-api"),
	})

	callbackResource := api.Root().AddResource(jsii.String("callback"), nil)
	callbackResource.AddMethod(
		jsii.String("GET"),
		awsapigateway.NewLambdaIntegration(approvalCallbackFn, nil),
		nil,
	)

	awsstepfunctions.NewStateMachine(stack, jsii.String("ProjectNotifierStateMachine"), &awsstepfunctions.StateMachineProps{
		StateMachineName: jsii.String("project-notifier-workflow"),
		DefinitionBody:   awsstepfunctions.DefinitionBody_FromFile((jsii.String("DailyProjectNotificationWorkflow.json")), nil),
	})

	return stack
}
