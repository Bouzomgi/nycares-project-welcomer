package main

type Config struct {
	AWS struct {
		Credentials struct {
			AccessKeyID     string `yaml:"accessKeyId"`
			SecretAccessKey string `yaml:"secretAccessKey"`
			Region          string `yaml:"region"`
		} `yaml:"credentials"`
		Dynamo struct {
			TableName string `yaml:"tableName"`
			Region    string `yaml:"region"`
			Endpoint  string `yaml:"endpoint,omitempty"` // optional for local
		} `yaml:"dynamo"`
		S3 struct {
			BucketName string `yaml:"bucketName"`
			Endpoint   string `yaml:"endpoint,omitempty"` // optional for local
		} `yaml:"s3"`
	} `yaml:"aws"`
}
