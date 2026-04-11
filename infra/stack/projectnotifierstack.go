package stack

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
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
	EnvSuffix     string // e.g. "-ci"; empty for LocalStack/production
}

func lambdaArchitecture() awslambda.Architecture {
	if os.Getenv("NYCARES_LAMBDA_ARCH") == "amd64" {
		return awslambda.Architecture_X86_64()
	}
	return awslambda.Architecture_ARM_64()
}

// lambdaAssetOptions returns asset options that force a new upload on every CI deploy.
// When COMMIT_SHA is set (i.e. in GitHub Actions), CDK uses it as a custom hash so
// CloudFormation always sees a changed S3 key and updates the Lambda code.
// The name suffix makes each lambda's hash unique so CDK does not deduplicate
// different binaries onto the same S3 asset.
// Locally (no COMMIT_SHA), CDK falls back to its default content-hash behavior.
func lambdaAssetOptions(name ...string) *awss3assets.AssetOptions {
	if sha := os.Getenv("COMMIT_SHA"); sha != "" {
		hash := sha
		if len(name) > 0 && name[0] != "" {
			hash = sha + "-" + name[0]
		}
		return &awss3assets.AssetOptions{
			AssetHash:     jsii.String(hash),
			AssetHashType: awscdk.AssetHashType_CUSTOM,
		}
	}
	return nil
}

// lambdaLogGroup imports a pre-existing log group by name.
// Lambda auto-creates log groups on first invocation, so importing (rather than
// creating) avoids CloudFormation conflicts when the group already exists.
func lambdaLogGroup(scope constructs.Construct, constructId, logGroupName string) awslogs.ILogGroup {
	return awslogs.LogGroup_FromLogGroupName(scope, jsii.String(constructId), jsii.String(logGroupName))
}

func ProjectNotifierStack(scope constructs.Construct, id string, props *LambdaStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	suffix := ""
	if props != nil {
		suffix = props.EnvSuffix
	}

	// --- AWS Resources ---

	tableName := "nycares-project-welcomer-notifications" + suffix
	bucketName := "nycares-project-welcomer-messages" + suffix

	var table awsdynamodb.ITable
	var bucket awss3.IBucket

	if suffix != "" {
		// Create owned resources for ephemeral PR environments
		table = awsdynamodb.NewTable(stack, jsii.String("SentNotifications"), &awsdynamodb.TableProps{
			TableName: jsii.String(tableName),
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String("ProjectName"),
				Type: awsdynamodb.AttributeType_STRING,
			},
			SortKey: &awsdynamodb.Attribute{
				Name: jsii.String("ProjectDate"),
				Type: awsdynamodb.AttributeType_STRING,
			},
			BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		})

		bucket = awss3.NewBucket(stack, jsii.String("MessageTemplates"), &awss3.BucketProps{
			BucketName:    jsii.String(bucketName),
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		})

	} else {
		// Import pre-existing resources (LocalStack / production without suffix)
		table = awsdynamodb.Table_FromTableName(stack, jsii.String("SentNotifications"), jsii.String(tableName))
		bucket = awss3.Bucket_FromBucketName(stack, jsii.String("MessageTemplates"), jsii.String(bucketName))
	}

	topic := awssns.NewTopic(stack, jsii.String("NotificationTopic"), &awssns.TopicProps{
		TopicName: jsii.String("nycares-notifications" + suffix),
	})

	debugQueue := awssqs.NewQueue(stack, jsii.String("DebugNotificationQueue"), &awssqs.QueueProps{
		QueueName: jsii.String("nycares-notifications-debug" + suffix),
	})
	topic.AddSubscription(awssnssubscriptions.NewSqsSubscription(debugQueue, nil))

	// --- Shared environment variables ---
	// Env var names must match viper config paths: prefix NYCARES_ + path with . replaced by _

	const ssmPath = "/nycares-project-welcomer/"

	sharedEnv := &map[string]*string{
		"NYCARES_AWS_DYNAMO_TABLENAME": jsii.String(tableName),
		"NYCARES_AWS_DYNAMO_REGION":    stack.Region(),
		"NYCARES_AWS_S3_BUCKETNAME":    jsii.String(bucketName),
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
		"NYCARES_MOCK_GENERATETHANKYOU",
		"NYCARES_ACCOUNT_USERNAME",
		"NYCARES_ACCOUNT_PASSWORD",
		"NYCARES_AWS_SF_APPROVALSECRET",
		"NYCARES_AWS_SES_SENDER",
		"NYCARES_AWS_SES_RECIPIENT",
	}
	for _, key := range passthroughEnvVars {
		if val := os.Getenv(key); val != "" {
			(*sharedEnv)[key] = jsii.String(val)
		}
	}

	// When a mock server URL is provided AND this is an ephemeral environment (suffix != ""),
	// override NYCARES_API_BASE_URL globally — integration tests mock the full API including
	// Login and FetchProjects. In production (no suffix), only SendAndPinMessage should route
	// to the mock server; Login/FetchProjects must hit the real NYC Cares API. The per-lambda
	// override for production is applied after the lambda loop below.
	if props != nil && props.MockServerUrl != nil && suffix != "" {
		(*sharedEnv)["NYCARES_API_BASE_URL"] = props.MockServerUrl
	}

	// --- Lambda Functions ---

	lambdaNames := []string{
		"Login",
		"FetchProjects",
		"RouteProject",
		"ComputePreProjectMessage",
		"GenerateThankYouMessage",
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

		timeout := jsii.Number(30)
		if name == "GenerateThankYouMessage" {
			timeout = jsii.Number(60)
		}

		fn := awslambda.NewFunction(stack, jsii.String(name), &awslambda.FunctionProps{
			Runtime: awslambda.Runtime_PROVIDED_AL2023(),
			Handler: jsii.String("bootstrap"),
			Code: awslambda.Code_FromAsset(
				jsii.String("../lambda-build/"+lowerName),
				lambdaAssetOptions(lowerName),
			),
			FunctionName: jsii.String(kebabName + suffix),
			Architecture: lambdaArchitecture(),
			Timeout:      awscdk.Duration_Seconds(timeout),
			Environment:  sharedEnv,
			LogGroup:     lambdaLogGroup(stack, name+"LogGroup", "/aws/lambda/"+kebabName+suffix),
		})

		lambdaFns[name] = fn
	}

	// Production mock mode (no suffix): only SendAndPinMessage routes to the mock server.
	// Login and FetchProjects continue hitting the real NYC Cares API so that real project
	// data and auth cookies are used.
	if props != nil && props.MockServerUrl != nil && suffix == "" {
		lambdaFns["SendAndPinMessage"].AddEnvironment(
			jsii.String("NYCARES_API_BASE_URL"), props.MockServerUrl, nil,
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

	// RouteProject needs DynamoDB read
	table.GrantReadData(lambdaFns["RouteProject"])

	// RecordMessage needs DynamoDB read/write
	table.GrantReadWriteData(lambdaFns["RecordMessage"])

	// SendAndPinMessage, RequestApprovalToSend, and GenerateThankYouMessage need S3 read
	bucket.GrantRead(lambdaFns["SendAndPinMessage"], nil)
	bucket.GrantRead(lambdaFns["RequestApprovalToSend"], nil)
	bucket.GrantRead(lambdaFns["GenerateThankYouMessage"], nil)

	// GenerateThankYouMessage needs Bedrock InvokeModel
	lambdaFns["GenerateThankYouMessage"].AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("bedrock:InvokeModel"),
		Resources: jsii.Strings(fmt.Sprintf("arn:aws:bedrock:%s::foundation-model/anthropic.claude-3-5-haiku-20241022-v1:0", *stack.Region())),
	}))

	// SNS publish for notification lambdas
	topic.GrantPublish(lambdaFns["RequestApprovalToSend"])
	topic.GrantPublish(lambdaFns["NotifyCompletion"])
	topic.GrantPublish(lambdaFns["DLQNotifier"])

	// --- Approval Callback Lambda (outside loop — needs API Gateway wiring) ---

	approvalCallbackFn := awslambda.NewFunction(stack, jsii.String("ApprovalCallback"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("../lambda-build/approvalcallback"), lambdaAssetOptions("approvalcallback")),
		FunctionName: jsii.String("approval-callback" + suffix),
		Architecture: lambdaArchitecture(),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
		Environment:  sharedEnv,
		LogGroup:     lambdaLogGroup(stack, "ApprovalCallbackLogGroup", "/aws/lambda/approval-callback"+suffix),
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
		RestApiName: jsii.String("approval-callback-api" + suffix),
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

	// --- SES Forwarder Lambda ---

	sesForwarderFn := awslambda.NewFunction(stack, jsii.String("SESForwarder"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("../lambda-build/sesforwarder"), lambdaAssetOptions("sesforwarder")),
		FunctionName: jsii.String("ses-forwarder" + suffix),
		Architecture: lambdaArchitecture(),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
		Environment:  sharedEnv,
		LogGroup:     lambdaLogGroup(stack, "SESForwarderLogGroup", "/aws/lambda/ses-forwarder"+suffix),
	})

	sesForwarderFn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ses:SendEmail"),
		Resources: jsii.Strings("*"),
	}))
	sesForwarderFn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ssm:GetParametersByPath"),
		Resources: jsii.Strings(ssmArn),
	}))

	topic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(sesForwarderFn, nil))

	// --- Step Functions State Machine ---

	// Build substitution map: workflow JSON uses ${LoginLambdaArn}, etc.
	definitionSubs := make(map[string]*string)
	for _, name := range lambdaNames {
		definitionSubs[name+"LambdaArn"] = lambdaFns[name].FunctionArn()
	}

	stateMachineLogGroup := awslogs.NewLogGroup(stack, jsii.String("ProjectNotifierStateMachineLogGroup"), &awslogs.LogGroupProps{
		LogGroupName:  jsii.String("/aws/states/project-notifier-workflow" + suffix),
		Retention:     awslogs.RetentionDays_THREE_MONTHS,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	stateMachine := awsstepfunctions.NewStateMachine(stack, jsii.String("ProjectNotifierStateMachine"), &awsstepfunctions.StateMachineProps{
		StateMachineName:        jsii.String("project-notifier-workflow" + suffix),
		DefinitionBody:          awsstepfunctions.DefinitionBody_FromFile(jsii.String("DailyProjectNotificationWorkflow.json"), nil),
		DefinitionSubstitutions: &definitionSubs,
		Logs: &awsstepfunctions.LogOptions{
			Destination:          stateMachineLogGroup,
			Level:                awsstepfunctions.LogLevel_ALL,
			IncludeExecutionData: jsii.Bool(true),
		},
	})

	// Grant the state machine permission to invoke all workflow lambdas
	for _, name := range lambdaNames {
		lambdaFns[name].GrantInvoke(stateMachine)
	}

	// --- Daily trigger at noon EST (17:00 UTC) ---

	awsevents.NewRule(stack, jsii.String("DailyNoonTrigger"), &awsevents.RuleProps{
		RuleName:    jsii.String("nycares-daily-noon-trigger" + suffix),
		Description: jsii.String("Triggers project-notifier-workflow daily at noon EST"),
		Schedule: awsevents.Schedule_Cron(&awsevents.CronOptions{
			Hour:   jsii.String("17"),
			Minute: jsii.String("0"),
		}),
		Targets: &[]awsevents.IRuleTarget{
			awseventstargets.NewSfnStateMachine(stateMachine, nil),
		},
	})

	// Grant the GHA deployer role permissions needed to seed resources and run integration tests
	if suffix != "" {
		if ghaRoleArn := os.Getenv("GHA_ROLE_ARN"); ghaRoleArn != "" {
			ghaRole := awsiam.Role_FromRoleArn(stack, jsii.String("GHARole"), jsii.String(ghaRoleArn), nil)
			bucket.GrantReadWrite(ghaRole, nil)
			table.GrantReadWriteData(ghaRole)
			ghaRole.AddToPrincipalPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: jsii.Strings(
					"states:StartExecution",
					"states:DescribeExecution",
					"states:GetExecutionHistory",
					"states:SendTaskSuccess",
					"states:SendTaskFailure",
				),
				Resources: jsii.Strings(*stateMachine.StateMachineArn(), "*"),
			}))
		}
	}

	return stack
}
