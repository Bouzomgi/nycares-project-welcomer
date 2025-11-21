package requestapproval

type Config struct {
	// TODO: fix callbackEndpoint, review NOPE
	CallbackEndpoint string `mapstructure:"callbackEndpoint"`
	AWS              struct {
		SNS struct {
			TopicArn string `mapstructure:"topicArn"`
		} `mapstructure:"sns"`
	} `mapstructure:"aws"`
}
