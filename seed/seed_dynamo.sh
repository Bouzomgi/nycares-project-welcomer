#!/bin/bash
set -e

ENDPOINT_URL="${AWS_ENDPOINT_URL:-http://localhost:4566}"

TABLE="nycares-project-welcomer-notifications"

echo "Creating DynamoDB table $TABLE..."
aws --endpoint-url "$ENDPOINT_URL" dynamodb create-table \
  --table-name "$TABLE" \
  --attribute-definitions \
    AttributeName=ProjectName,AttributeType=S \
    AttributeName=ProjectDate,AttributeType=S \
  --key-schema \
    AttributeName=ProjectName,KeyType=HASH \
    AttributeName=ProjectDate,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST

echo "Table created!"
