#!/bin/bash
set -e

TABLE="${DYNAMO_TABLE_NAME:-nycares-project-welcomer-notifications}"
BUCKET="${S3_BUCKET_NAME:-nycares-project-welcomer-messages}"

aws_cmd() {
  if [ -n "${AWS_ENDPOINT_URL:-}" ]; then
    aws --endpoint-url "$AWS_ENDPOINT_URL" "$@"
  else
    aws "$@"
  fi
}

echo "Creating DynamoDB table $TABLE if it doesn't exist..."
aws_cmd dynamodb create-table \
  --table-name "$TABLE" \
  --attribute-definitions AttributeName=ProjectName,AttributeType=S AttributeName=ProjectDate,AttributeType=S \
  --key-schema AttributeName=ProjectName,KeyType=HASH AttributeName=ProjectDate,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null || echo "Table already exists"

echo "Creating S3 bucket $BUCKET if it doesn't exist..."
aws_cmd s3 mb "s3://$BUCKET" 2>/dev/null || echo "Bucket already exists"

echo "LocalStack seed complete."
