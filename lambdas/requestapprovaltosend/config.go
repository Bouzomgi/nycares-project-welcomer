package main

type Config struct {
	CallbackEndpoint string `yaml:"callbackEndpoint"`
	AWS              struct {
		Credentials struct {
			AccessKeyID     string `yaml:"accessKeyId"`
			SecretAccessKey string `yaml:"secretAccessKey"`
			Region          string `yaml:"region"`
		} `yaml:"credentials"`
		SNS struct {
			TopicARN string `yaml:"topicArn"`
		} `yaml:"sns"`
	} `yaml:"aws"`
}
