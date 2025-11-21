package awsconfig

import (
	"context"
	"errors"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	sdkConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// LoadAWSConfigFromConfig inspects a config struct for AWS credentials and region.
// Either all three fields exist and are non-empty, or none exist. Partial -> error.
func LoadAWSConfigFromConfig[T any](ctx context.Context, cfg T) (aws.Config, error) {
	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	credsVal := v.FieldByName("AWS")
	var accessKey, secretKey, region string

	if credsVal.IsValid() && credsVal.Kind() == reflect.Struct {
		credsStruct := credsVal.FieldByName("Credentials")
		if credsStruct.IsValid() && credsStruct.Kind() == reflect.Struct {
			accessKeyF := credsStruct.FieldByName("AccessKeyID")
			secretKeyF := credsStruct.FieldByName("SecretAccessKey")
			regionF := credsStruct.FieldByName("Region")

			if accessKeyF.IsValid() && accessKeyF.Kind() == reflect.String {
				accessKey = accessKeyF.String()
			}
			if secretKeyF.IsValid() && secretKeyF.Kind() == reflect.String {
				secretKey = secretKeyF.String()
			}
			if regionF.IsValid() && regionF.Kind() == reflect.String {
				region = regionF.String()
			}
		}
	}

	// Check for all-or-nothing
	allEmpty := accessKey == "" && secretKey == "" && region == ""
	allFilled := accessKey != "" && secretKey != "" && region != ""

	if !allEmpty && !allFilled {
		return aws.Config{}, errors.New("AWS credentials must be either all defined or all empty")
	}

	opts := []func(*sdkConfig.LoadOptions) error{}
	if allFilled {
		opts = append(opts, sdkConfig.WithRegion(region))
		opts = append(opts, sdkConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		))
	}

	return sdkConfig.LoadDefaultConfig(ctx, opts...)
}
