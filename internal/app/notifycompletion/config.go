package notifycompletion

type Config struct {
	AWS struct {
		SNS struct {
			TopicArn string `mapstructure:"topicArn"`
		} `mapstructure:"sns"`
	} `mapstructure:"aws"`
	Mock struct {
		SendMessage bool `mapstructure:"sendMessage,omitempty"`
	} `mapstructure:"mock"`
}
