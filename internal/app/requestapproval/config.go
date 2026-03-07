package requestapproval

type Config struct {
	AWS struct {
		SNS struct {
			TopicArn string `mapstructure:"topicArn"`
		} `mapstructure:"sns"`
		SF struct {
			CallbackEndpoint string `mapstructure:"callbackEndpoint"`
			ApprovalSecret   string `mapstructure:"approvalSecret"`
		} `mapstructure:"sf"`
	} `mapstructure:"aws"`
	Mock struct {
		SendMessage bool `mapstructure:"sendMessage,omitempty"`
	} `mapstructure:"mock"`
}
