package main

import (
	"os"

	"nycares-project-welcomer-infra/stack"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	env := &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}

	var mockServerUrl *string
	if os.Getenv("DEPLOY_MOCKSERVER") == "true" {
		_, url := stack.MockServerStack(app, "NYCaresMockServerStack", &awscdk.StackProps{Env: env})
		mockServerUrl = url
	}

	stack.ProjectNotifierStack(app, "NycaresWelcomerLambdaStack", &stack.LambdaStackProps{
		StackProps: awscdk.StackProps{
			Env: env,
		},
		MockServerUrl: mockServerUrl,
	})

	app.Synth(nil)
}
