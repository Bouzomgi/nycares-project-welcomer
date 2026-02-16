#!/bin/bash
set -e

# Bucket name
BUCKET="cdk-hnb659fds-assets-000000000000-us-east-1"

# Ensure the bucket exists
echo "Creating bucket $BUCKET if it doesn't exist..."
awslocal s3 mb s3://$BUCKET || echo "Bucket already exists"

# Path to local files to upload
FILES=(
  "seed/reminder.txt"
  "seed/welcome.txt"
)

# Upload files
for FILE in "${FILES[@]}"; do
  BASENAME=$(basename "$FILE")
  echo "Uploading $FILE to s3://$BUCKET/$BASENAME..."
  awslocal s3 cp "$FILE" "s3://$BUCKET/$BASENAME"
done

echo "All files uploaded!"
