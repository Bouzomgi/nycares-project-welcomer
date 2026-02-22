#!/bin/bash
set -e

ENDPOINT_URL="${AWS_ENDPOINT_URL:-http://localhost:4566}"

BUCKET="message-storage-bucket"

echo "Creating bucket $BUCKET if it doesn't exist..."
aws --endpoint-url "$ENDPOINT_URL" s3 mb "s3://$BUCKET" 2>/dev/null || echo "Bucket already exists"

FILES_DIR="${FILES_DIR:-seed/s3Items}"

echo "Syncing $FILES_DIR to s3://$BUCKET..."
aws --endpoint-url "$ENDPOINT_URL" s3 sync "$FILES_DIR" "s3://$BUCKET" --exclude "*.txt"

echo "All files uploaded!"
