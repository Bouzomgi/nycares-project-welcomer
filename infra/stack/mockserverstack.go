package stack

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func MockServerStack(scope constructs.Construct, id string, props *awscdk.StackProps) (awscdk.Stack, *string) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	suffix := os.Getenv("ENV_SUFFIX")

	fn := awslambda.NewFunction(stack, jsii.String("MockServer"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("bootstrap"),
		Architecture: lambdaArchitecture(),
		FunctionName: jsii.String("mock-server" + suffix),
		Code:         awslambda.Code_FromAsset(jsii.String("../lambda-build/mockserver"), lambdaAssetOptions()),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
		LogGroup:     lambdaLogGroup(stack, "MockServerLogGroup", "/aws/lambda/mock-server"+suffix),
	})

	fnUrl := awslambda.NewFunctionUrl(stack, jsii.String("MockServerUrl"), &awslambda.FunctionUrlProps{
		Function: fn,
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	awscdk.NewCfnOutput(stack, jsii.String("MockServerUrlOutput"), &awscdk.CfnOutputProps{
		Value:       fnUrl.Url(),
		Description: jsii.String("Mock server function URL"),
		ExportName:  jsii.String("MockServerUrl" + suffix),
	})

	return stack, fnUrl.Url()
}
