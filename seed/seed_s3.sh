#!/bin/bash
set -e

BUCKET="${S3_BUCKET_NAME:-nycares-project-welcomer-messages}"

aws_cmd() {
  if [ -n "${AWS_ENDPOINT_URL:-}" ]; then
    aws --endpoint-url "$AWS_ENDPOINT_URL" "$@"
  else
    aws "$@"
  fi
}

echo "Creating bucket $BUCKET if it doesn't exist..."
aws_cmd s3 mb "s3://$BUCKET" 2>/dev/null || echo "Bucket already exists"

FILES_DIR="${FILES_DIR:-seed/s3Items}"

echo "Syncing $FILES_DIR to s3://$BUCKET..."
aws_cmd s3 sync "$FILES_DIR" "s3://$BUCKET" --exclude "*.txt"

echo "All files uploaded!"
