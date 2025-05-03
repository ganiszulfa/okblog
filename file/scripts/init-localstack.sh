#!/bin/bash
set -e

echo "Waiting for LocalStack to be ready..."
while ! nc -z localstack 4566; do
  sleep 1
done

echo "Creating S3 bucket..."
aws --endpoint-url=http://localstack:4566 s3 mb s3://file-bucket

echo "Setting bucket policy to allow public access..."
aws --endpoint-url=http://localstack:4566 s3api put-bucket-policy \
  --bucket file-bucket \
  --policy '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "PublicReadGetObject",
        "Effect": "Allow",
        "Principal": "*",
        "Action": "s3:GetObject",
        "Resource": "arn:aws:s3:::file-bucket/*"
      }
    ]
  }'

echo "LocalStack initialization completed!" 