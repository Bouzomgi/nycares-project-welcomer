package main

type Config struct {
	AWS struct {
		SES struct {
			Sender    string `mapstructure:"sender"`
			Recipient string `mapstructure:"recipient"`
		} `mapstructure:"ses"`
		SNS struct {
			TopicArn string `mapstructure:"topicArn"`
		} `mapstructure:"sns"`
	} `mapstructure:"aws"`
}
