#!/bin/bash
awslocal dynamodb create-table \
  --table-name Sent_Notifications \
  --attribute-definitions \
    AttributeName=ProjectName,AttributeType=S \
    AttributeName=ProjectDate,AttributeType=S \
  --key-schema \
    AttributeName=ProjectName,KeyType=HASH \
    AttributeName=ProjectDate,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST