package requestapproval

type Config struct {
	CallbackEndpoint string `mapstructure:"callbackEndpoint"`
	AWS              struct {
		SNS struct {
			TopicArn string `mapstructure:"topicArn"`
		} `mapstructure:"sns"`
	} `mapstructure:"aws"`
}
