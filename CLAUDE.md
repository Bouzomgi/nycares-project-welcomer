# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

NYC Cares Project Welcomer — a serverless notification system that sends welcome/reminder messages to New York Cares project attendees.

This is a Go/CDK project deployed to AWS. Primary languages: Go for Lambda functions, TypeScript for CDK infrastructure, YAML for CI/CD (GitHub Actions). Always consider deployment implications when making infrastructure changes.

**Module:** `github.com/Bouzomgi/nycares-project-welcomer`

## Build & Run Commands

### Build all lambdas

```bash
docker compose up --build
```

This compiles all 10 Lambda functions in parallel via Docker (`CGO_ENABLED=0 GOOS=linux GOARCH=arm64`). Output goes to `lambda-build/`.

### Run tests

```bash
go test ./...
```

### Deploy infrastructure (CDK)

```bash
cd infra && cdk deploy
```

### Run mock server (local dev)

The mock server at [`internal/mockserver/`](internal/mockserver/CLAUDE.md) simulates the NYC Cares API. Used alongside LocalStack for local AWS services (S3, DynamoDB, SNS).

## Architecture

### Workflow

A Step Functions state machine orchestrates 8 Lambda functions per project:

1. **Login** → authenticate with NYC Cares API
2. **FetchProjects** → get upcoming projects
3. **ComputeMessageToSend** → decide welcome vs reminder (7+ days = welcome, 2+ days = reminder)
4. **RequestApprovalToSend** → SNS notification with `waitForTaskToken` callback
5. **SendAndPinMessage** → post message to project channel
6. **RecordMessage** → update DynamoDB tracking
7. **NotifyCompletion** → SNS success notification
8. **DLQNotifier** → error handler (catch blocks route here)

Two additional Lambdas support the workflow outside the state machine:

- **ApprovalCallback** → API Gateway endpoint that receives the human approval decision and resumes the state machine via `SendTaskSuccess`/`SendTaskFailure`
- **SESForwarder** → subscribed to the SNS notifications topic; forwards messages as HTML email via SES

### Layered Structure

- **`internal/domain/`** — Domain models: `Auth`, `Project`, `ProjectNotification`
- **`internal/app/`** — Use cases per workflow step (one package per Lambda)
- **`internal/platform/`** — AWS service integrations (S3, DynamoDB, SNS, HTTP with cookie jar)
- **`internal/models/`** — Lambda I/O models for Step Functions serialization
- **`internal/config/`** — YAML config locally, env vars (`NYCARES_` prefix) in Lambda
- **`lambda/`** — 10 entry points, each with `main.go` (init) + `handler.go` (logic)
- **`infra/`** — AWS CDK stack in Go + Step Functions workflow JSON

### Data Flow

Auth cookies, project metadata, message type, and task tokens are passed through the workflow via Lambda I/O models serialized as JSON between steps.

### DynamoDB Schema (`nycares-project-welcomer-notifications` table)

- **Key:** `ProjectName` + `ProjectDate` (composite)
- **Fields:** `HasSentWelcome`, `HasSentReminder`, `ShouldStopNotify`, `LastUpdated`, `ProjectId`

## Configuration

Config loads from `config.yaml` locally or environment variables in Lambda (auto-detected). See `config.template.yaml` for structure. Key sections: `aws` (dynamo, s3, sns endpoints), `account` (credentials).

## Development Workflow

Always run `gofmt -w .` from both the repo root and `infra/` before every commit.

## Key Dependencies

- `aws-lambda-go` — Lambda runtime
- `aws-sdk-go-v2` — S3, DynamoDB, SNS, SES, SSM, Step Functions clients
- `aws-cdk-go` — Infrastructure as code
- `spf13/viper` — Config management
- `gorilla/mux` — Mock server routing
