package s3service

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (s *S3Service) CreateS3Path(suffix string) string {
	return fmt.Sprintf("s3://%s", suffix)
}

func (s *S3Service) ObjectExists(ctx context.Context, key string) (bool, error) {
	_, err := s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err == nil {
		return true, nil
	}

	var nfe *types.NotFound
	if errors.As(err, &nfe) {
		return false, nil
	}

	return false, err
}
