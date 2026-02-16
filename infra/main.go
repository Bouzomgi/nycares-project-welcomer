package main

import (
	"nycares-project-welcomer-infra/stack"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	stack.ProjectNotifierStack(app, "LambdaStack", &stack.LambdaStackProps{
		StackProps: awscdk.StackProps{
			Env: &awscdk.Environment{
				Account: jsii.String("000000000000"),
				Region:  jsii.String("us-east-1"),
			},
		},
	})

	app.Synth(nil)
}
