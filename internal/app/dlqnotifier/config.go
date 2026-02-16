package dlqnotifier

type Config struct {
	AWS struct {
		SNS struct {
			TopicArn string `mapstructure:"topicArn"`
		} `mapstructure:"sns"`
	} `mapstructure:"aws"`
}
