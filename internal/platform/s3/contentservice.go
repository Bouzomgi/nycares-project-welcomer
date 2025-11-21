package s3service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type ContentService interface {
	GetMessageContent(ctx context.Context, ref string) (string, error)
}

func (s *S3Service) GetMessageContent(ctx context.Context, messageRef string) (string, error) {

	resp, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(messageRef),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return "", fmt.Errorf("messageRef %s not found", messageRef)
		}
		return "", fmt.Errorf("failed to fetch S3 object %s: %w", messageRef, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read S3 body: %w", err)
	}

	return string(data), nil
}
