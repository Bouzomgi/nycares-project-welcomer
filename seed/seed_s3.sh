#!/bin/bash
set -e

ENDPOINT_URL="${AWS_ENDPOINT_URL:-http://localhost:4566}"

BUCKET="nycares-message-templates"

echo "Creating bucket $BUCKET if it doesn't exist..."
aws --endpoint-url "$ENDPOINT_URL" s3 mb "s3://$BUCKET" 2>/dev/null || echo "Bucket already exists"

FILES_DIR="${FILES_DIR:-seed/s3Items}"

for FILE in "$FILES_DIR"/*.txt; do
  BASENAME=$(basename "$FILE")
  echo "Uploading $FILE to s3://$BUCKET/$BASENAME..."
  aws --endpoint-url "$ENDPOINT_URL" s3 cp "$FILE" "s3://$BUCKET/$BASENAME"
done

echo "All files uploaded!"
