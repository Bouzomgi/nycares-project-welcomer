package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/computemessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/fetchprojects"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/login"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/notifycompletion"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/requestapproval"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/sendandpinmessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	dynamoservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/dynamo"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http"
	s3service "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/s3"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func buildLoginHandler() (*login.LoginHandler, error) {
	cfg, err := config.LoadConfig[login.Config]()
	if err != nil {
		return nil, err
	}

	httpSvc, err := httpservice.NewHttpService(endpoints.BaseUrl)
	if err != nil {
		return nil, err
	}

	usecase := login.NewLoginUseCase(httpSvc)
	return login.NewLoginHandler(usecase, cfg), nil
}

func buildFetchProjectsHandler() (*fetchprojects.FetchProjectsHandler, error) {
	cfg, err := config.LoadConfig[fetchprojects.Config]()
	if err != nil {
		return nil, err
	}

	httpSvc, err := httpservice.NewHttpService(endpoints.BaseUrl)
	if err != nil {
		return nil, err
	}

	usecase := fetchprojects.NewFetchProjectsUseCase(httpSvc)
	return fetchprojects.NewFetchProjectsHandler(usecase, cfg), nil
}

func buildComputeMessageHandler() (*computemessage.ComputeMessageHandler, error) {
	cfg, err := config.LoadConfig[computemessage.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	dynamoClient := dynamodb.NewFromConfig(awsCfg)
	dynamoSvc := dynamoservice.NewDynamoService(dynamoClient, cfg.AWS.Dynamo.TableName)

	usecase := computemessage.NewComputeMessageUseCase(dynamoSvc)
	return computemessage.NewComputeMessageHandler(usecase, cfg), nil
}

func buildRequestApprovalHandler() (*requestapproval.RequestApprovalHandler, error) {
	cfg, err := config.LoadConfig[requestapproval.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	snsClient := sns.NewFromConfig(awsCfg)

	snsSvc := snsservice.NewSNSSerice(snsClient, cfg.AWS.SNS.TopicArn)

	usecase := requestapproval.NewRequestApprovalUseCase(snsSvc)
	return requestapproval.NewRequestApprovalHandler(usecase, cfg), nil
}

func buildSendAndPinMessageHandler() (*sendandpinmessage.SendAndPinMessageHandler, error) {
	cfg, err := config.LoadConfig[sendandpinmessage.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	httpSvc, err := httpservice.NewHttpService(endpoints.BaseUrl)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg)
	s3Svc := s3service.NewS3Service(s3Client, cfg.AWS.S3.BucketName)

	usecase := sendandpinmessage.NewSendAndPinMessageUseCase(s3Svc, httpSvc)
	return sendandpinmessage.NewSendAndPinMessageHandler(usecase, cfg), nil
}

func buildNotifyCompletionHandler() (*notifycompletion.NotifyCompletionHandler, error) {
	cfg, err := config.LoadConfig[notifycompletion.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	snsClient := sns.NewFromConfig(awsCfg)

	snsSvc := snsservice.NewSNSSerice(snsClient, cfg.AWS.SNS.TopicArn)

	usecase := notifycompletion.NewNotifyCompletionUseCase(snsSvc)
	return notifycompletion.NewNotifyCompletionHandler(usecase, cfg), nil
}

////////////////

func main() {

	loginHandler, err := buildLoginHandler()
	if err != nil {
		panic(err)
	}

	loginOut, err := loginHandler.Handle(context.Background())
	if err != nil {
		panic(err)
	}

	data, _ := json.MarshalIndent(loginOut, "", "  ")

	fmt.Println("//////// Login Output ////////")
	fmt.Println(string(data))

	///// Fetch Projects

	fetchProjectsHandler, err := buildFetchProjectsHandler()
	if err != nil {
		panic(err)
	}

	fetchProjectsOut, err := fetchProjectsHandler.Handle(context.Background(), loginOut)

	if err != nil {
		panic(err)
	}

	data, _ = json.MarshalIndent(fetchProjectsOut, "", "  ")

	fmt.Print("\n\n")
	fmt.Println("//////// FetchProjects Output ////////")
	fmt.Println(string(data))

	///// Compute Message

	if len(fetchProjectsOut.Projects) == 0 {
		panic("could not fetch any projects")
	}

	computeMessageInput := models.ComputeMessageInput{
		Auth:    fetchProjectsOut.Auth,
		Project: fetchProjectsOut.Projects[0],
	}

	computeMessageHandler, err := buildComputeMessageHandler()
	if err != nil {
		panic(err)
	}

	computeMessageOut, err := computeMessageHandler.Handle(context.Background(), computeMessageInput)

	if err != nil {
		panic(err)
	}

	data, _ = json.MarshalIndent(computeMessageOut, "", "  ")

	fmt.Print("\n\n")
	fmt.Println("//////// ComputeMessage Output ////////")
	fmt.Println(string(data))

	///// Request Approval

	requestApprovalHandler, err := buildRequestApprovalHandler()
	if err != nil {
		panic(err)
	}

	requestApprovalIn := models.RequestApprovalInput{
		TaskToken:     "dummy-task-token",
		Auth:          computeMessageOut.Auth,
		Project:       computeMessageOut.Project,
		MessageToSend: computeMessageOut.MessageToSend,
	}

	requestApprovalOut, err := requestApprovalHandler.Handle(context.Background(), requestApprovalIn)

	if err != nil {
		panic(err)
	}

	data, _ = json.MarshalIndent(requestApprovalOut, "", "  ")

	fmt.Print("\n\n")
	fmt.Println("//////// FetchProjects Output ////////")
	fmt.Println(string(data))

	///// Send and Pin Message

	sendAndPinMessageHandler, err := buildSendAndPinMessageHandler()
	if err != nil {
		panic(err)
	}

	sendAndPinMessageOut, err := sendAndPinMessageHandler.Handle(context.Background(), requestApprovalOut)

	if err != nil {
		panic(err)
	}

	data, _ = json.MarshalIndent(sendAndPinMessageOut, "", "  ")

	fmt.Print("\n\n")
	fmt.Println("//////// SendAndPinMessage Output ////////")
	fmt.Println(string(data))

	///// Notify Completion

	notifyCompletionHandler, err := buildNotifyCompletionHandler()
	if err != nil {
		panic(err)
	}

	err = notifyCompletionHandler.Handle(context.Background(), sendAndPinMessageOut)

	if err != nil {
		panic(err)
	}

	fmt.Print("\n\n")
	fmt.Println("//////// Notify Completion ////////")

	return
}
