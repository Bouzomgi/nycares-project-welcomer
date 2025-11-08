#!/bin/bash
set -e

echo "🚀 Creating test S3 bucket and DynamoDB table..."

# Localstack S3 bucket
awslocal s3 mb s3://message-bucket

# DynamoDB table with composite primary key: ProjectName (HASH), ProjectDate (RANGE)
awslocal dynamodb create-table \
  --table-name message-table \
  --attribute-definitions \
    AttributeName=ProjectName,AttributeType=S \
    AttributeName=ProjectDate,AttributeType=S \
  --key-schema \
    AttributeName=ProjectName,KeyType=HASH \
    AttributeName=ProjectDate,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST

# Insert sample items
awslocal dynamodb put-item \
  --table-name message-table \
  --item '{
    "ProjectName": {"S": "ProjectA"},
    "ProjectDate": {"S": "2025-01-01"},
    "hasSentWelcome": {"BOOL": true},
    "hasSentReminder": {"BOOL": false},
    "shouldStopNotify": {"BOOL": false},
    "lastUpdated": {"S": "2025-01-01"}
  }'

awslocal dynamodb put-item \
  --table-name message-table \
  --item '{
    "ProjectName": {"S": "ProjectB"},
    "ProjectDate": {"S": "2025-01-01"},
    "hasSentWelcome": {"BOOL": true},
    "hasSentReminder": {"BOOL": true},
    "shouldStopNotify": {"BOOL": false},
    "lastUpdated": {"S": "2025-01-01"}
  }'

awslocal dynamodb put-item \
  --table-name message-table \
  --item '{
    "ProjectName": {"S": "ProjectC"},
    "ProjectDate": {"S": "2025-01-01"},
    "hasSentWelcome": {"BOOL": true},
    "hasSentReminder": {"BOOL": false},
    "shouldStopNotify": {"BOOL": true},
    "lastUpdated": {"S": "2025-01-01"}
  }'

echo "✅ Resources created!"
