package main

type Config struct {
	AWS struct {
		Credentials struct {
			AccessKeyID     string `yaml:"accessKeyId"`
			SecretAccessKey string `yaml:"secretAccessKey"`
			Region          string `yaml:"region"`
		} `yaml:"credentials"`
		S3 struct {
			BucketName string `yaml:"bucketName"`
		} `yaml:"S3"`
	} `yaml:"aws"`
}
