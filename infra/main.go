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

	suffix := os.Getenv("ENV_SUFFIX")

	envName := "prod"
	if suffix != "" {
		envName = suffix[1:] // strip leading "-", e.g. "-ci" → "ci"
	}

	awscdk.Tags_Of(app).Add(jsii.String("Project"), jsii.String("nycares-project-welcomer"), nil)
	awscdk.Tags_Of(app).Add(jsii.String("Environment"), jsii.String(envName), nil)

	var mockServerUrl *string
	if os.Getenv("DEPLOY_MOCKSERVER") == "true" {
		_, url := stack.MockServerStack(app, "NYCaresMockServerStack"+suffix, &awscdk.StackProps{Env: env})
		mockServerUrl = url
	}

	stack.ProjectNotifierStack(app, "NYCaresLambdaStack"+suffix, &stack.LambdaStackProps{
		StackProps: awscdk.StackProps{
			Env: env,
		},
		MockServerUrl: mockServerUrl,
		EnvSuffix:     suffix,
	})

	app.Synth(nil)
}
