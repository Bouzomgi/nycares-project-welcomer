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
		Account: jsii.String("000000000000"),
		Region:  jsii.String("us-east-1"),
	}

	var mockServerUrl *string
	if os.Getenv("DEPLOY_MOCKSERVER") == "true" {
		_, url := stack.MockServerStack(app, "MockServerStack", &awscdk.StackProps{Env: env})
		mockServerUrl = url
	}

	stack.ProjectNotifierStack(app, "LambdaStack", &stack.LambdaStackProps{
		StackProps: awscdk.StackProps{
			Env: env,
		},
		MockServerUrl: mockServerUrl,
	})

	app.Synth(nil)
}
