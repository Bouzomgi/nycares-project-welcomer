package s3service

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	s3Client   s3.Client
	bucketName string
}

func NewS3Service(cfg aws.Config, bucketName string) *S3Service {
	return &S3Service{
		s3Client:   *s3.NewFromConfig(cfg),
		bucketName: bucketName,
	}
}
