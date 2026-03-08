package stack

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssnssubscriptions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/iancoleman/strcase"
)

type LambdaStackProps struct {
	awscdk.StackProps
	MockServerUrl *string
}

func lambdaArchitecture() awslambda.Architecture {
	if os.Getenv("NYCARES_LAMBDA_ARCH") == "amd64" {
		return awslambda.Architecture_X86_64()
	}
	return awslambda.Architecture_ARM_64()
}

func ProjectNotifierStack(scope constructs.Construct, id string, props *LambdaStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// --- AWS Resources ---

	table := awsdynamodb.NewTable(stack, jsii.String("SentNotifications"), &awsdynamodb.TableProps{
		TableName: jsii.String("nycares-project-welcomer-notifications"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("ProjectName"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("ProjectDate"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	})

	bucket := awss3.NewBucket(stack, jsii.String("MessageTemplates"), &awss3.BucketProps{
		BucketName:    jsii.String("nycares-project-welcomer-messages"),
		RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	})

	topic := awssns.NewTopic(stack, jsii.String("NotificationTopic"), &awssns.TopicProps{
		TopicName: jsii.String("nycares-notifications"),
	})

	debugQueue := awssqs.NewQueue(stack, jsii.String("DebugNotificationQueue"), &awssqs.QueueProps{
		QueueName: jsii.String("nycares-notifications-debug"),
	})
	topic.AddSubscription(awssnssubscriptions.NewSqsSubscription(debugQueue, nil))

	// --- Shared environment variables ---
	// Env var names must match viper config paths: prefix NYCARES_ + path with . replaced by _

	const ssmPath = "/nycares-project-welcomer/"

	sharedEnv := &map[string]*string{
		"NYCARES_AWS_DYNAMO_TABLENAME": table.TableName(),
		"NYCARES_AWS_DYNAMO_REGION":    stack.Region(),
		"NYCARES_AWS_S3_BUCKETNAME":    bucket.BucketName(),
		"NYCARES_AWS_SNS_TOPICARN":     topic.TopicArn(),
		"NYCARES_SSM_PATH":             jsii.String(ssmPath),
	}

	// Tag all stack resources with the commit SHA
	if commitSha := os.Getenv("COMMIT_SHA"); commitSha != "" {
		awscdk.Tags_Of(stack).Add(jsii.String("CommitSha"), jsii.String(commitSha), nil)
	}

	// Passthrough env vars from deploy environment
	passthroughEnvVars := []string{
		"NYCARES_API_BASE_URL",
		"NYCARES_CURRENT_DATE",
		"NYCARES_MOCK_SENDMESSAGE",
		"NYCARES_ACCOUNT_USERNAME",
		"NYCARES_ACCOUNT_PASSWORD",
		"NYCARES_AWS_SF_APPROVALSECRET",
	}
	for _, key := range passthroughEnvVars {
		if val := os.Getenv(key); val != "" {
			(*sharedEnv)[key] = jsii.String(val)
		}
	}

	// --- Lambda Functions ---

	lambdaNames := []string{
		"Login",
		"FetchProjects",
		"ComputeMessageToSend",
		"RequestApprovalToSend",
		"SendAndPinMessage",
		"RecordMessage",
		"NotifyCompletion",
		"DLQNotifier",
	}

	lambdaFns := make(map[string]awslambda.Function)

	for _, name := range lambdaNames {
		lowerName := strings.ToLower(name)
		kebabName := strcase.ToKebab(name)

		fn := awslambda.NewFunction(stack, jsii.String(name), &awslambda.FunctionProps{
			Runtime: awslambda.Runtime_PROVIDED_AL2023(),
			Handler: jsii.String("bootstrap"),
			Code: awslambda.Code_FromAsset(
				jsii.String("../lambda-build/"+lowerName),
				nil,
			),
			FunctionName: jsii.String(kebabName),
			Architecture: lambdaArchitecture(),
			Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
			Environment:  sharedEnv,
			LogRetention: awslogs.RetentionDays_THREE_MONTHS,
		})

		lambdaFns[name] = fn
	}

	// Override API base URL for SendAndPinMessage in dry-run mode
	if props != nil && props.MockServerUrl != nil {
		lambdaFns["SendAndPinMessage"].AddEnvironment(
			jsii.String("NYCARES_API_BASE_URL"),
			props.MockServerUrl,
			nil,
		)
	}

	// --- IAM Permissions ---

	// All lambdas need to read SSM parameters
	ssmArn := fmt.Sprintf("arn:aws:ssm:%s:%s:parameter%s*", *stack.Region(), *stack.Account(), ssmPath)
	for _, name := range lambdaNames {
		lambdaFns[name].AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Actions:   jsii.Strings("ssm:GetParametersByPath"),
			Resources: jsii.Strings(ssmArn),
		}))
	}

	// ComputeMessageToSend needs DynamoDB read
	table.GrantReadData(lambdaFns["ComputeMessageToSend"])

	// RecordMessage needs DynamoDB read/write
	table.GrantReadWriteData(lambdaFns["RecordMessage"])

	// SendAndPinMessage and RequestApprovalToSend need S3 read
	bucket.GrantRead(lambdaFns["SendAndPinMessage"], nil)
	bucket.GrantRead(lambdaFns["RequestApprovalToSend"], nil)

	// SNS publish for notification lambdas
	topic.GrantPublish(lambdaFns["RequestApprovalToSend"])
	topic.GrantPublish(lambdaFns["NotifyCompletion"])
	topic.GrantPublish(lambdaFns["DLQNotifier"])

	// --- Approval Callback Lambda (outside loop — needs API Gateway wiring) ---

	approvalCallbackFn := awslambda.NewFunction(stack, jsii.String("ApprovalCallback"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("../lambda-build/approvalcallback"), nil),
		FunctionName: jsii.String("approval-callback"),
		Architecture: lambdaArchitecture(),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
		Environment:  sharedEnv,
		LogRetention: awslogs.RetentionDays_THREE_MONTHS,
	})

	approvalCallbackFn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("states:SendTaskSuccess", "states:SendTaskFailure"),
		Resources: jsii.Strings("*"),
	}))
	approvalCallbackFn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ssm:GetParametersByPath"),
		Resources: jsii.Strings(ssmArn),
	}))

	// --- API Gateway ---

	api := awsapigateway.NewRestApi(stack, jsii.String("ApprovalCallbackApi"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("approval-callback-api"),
	})

	callbackResource := api.Root().AddResource(jsii.String("callback"), nil)
	callbackResource.AddMethod(
		jsii.String("GET"),
		awsapigateway.NewLambdaIntegration(approvalCallbackFn, nil),
		nil,
	)

	// Set the callback endpoint from the API Gateway URL (resolved at deploy time)
	lambdaFns["RequestApprovalToSend"].AddEnvironment(
		jsii.String("NYCARES_AWS_SF_CALLBACKENDPOINT"),
		api.Url(),
		nil,
	)

	// --- Step Functions State Machine ---

	// Build substitution map: workflow JSON uses ${LoginLambdaArn}, etc.
	definitionSubs := make(map[string]*string)
	for _, name := range lambdaNames {
		definitionSubs[name+"LambdaArn"] = lambdaFns[name].FunctionArn()
	}

	stateMachine := awsstepfunctions.NewStateMachine(stack, jsii.String("ProjectNotifierStateMachine"), &awsstepfunctions.StateMachineProps{
		StateMachineName:        jsii.String("project-notifier-workflow"),
		DefinitionBody:          awsstepfunctions.DefinitionBody_FromFile(jsii.String("DailyProjectNotificationWorkflow.json"), nil),
		DefinitionSubstitutions: &definitionSubs,
	})

	// Grant the state machine permission to invoke all workflow lambdas
	for _, name := range lambdaNames {
		lambdaFns[name].GrantInvoke(stateMachine)
	}

	return stack
}
