package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func MockServerStack(scope constructs.Construct, id string, props *awscdk.StackProps) (awscdk.Stack, *string) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	fn := awslambda.NewFunction(stack, jsii.String("MockServer"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("bootstrap"),
		Architecture: awslambda.Architecture_ARM_64(),
		FunctionName: jsii.String("mock-server"),
		Code: awslambda.Code_FromAsset(jsii.String("../"), &awss3assets.AssetOptions{
			Bundling: &awscdk.BundlingOptions{
				Image: awscdk.DockerImage_FromRegistry(jsii.String("golang:1.25-alpine")),
				Command: jsii.Strings(
					"sh", "-c",
					"CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /asset-output/bootstrap ./internal/mockserver",
				),
				OutputType: awscdk.BundlingOutput_NOT_ARCHIVED,
			},
		}),
	})

	fnUrl := awslambda.NewFunctionUrl(stack, jsii.String("MockServerUrl"), &awslambda.FunctionUrlProps{
		Function: fn,
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	return stack, fnUrl.Url()
}
