package sendandpinmessage

type Config struct {
	Api struct {
		BaseUrl string `mapstructure:"baseUrl,omitempty"`
	} `mapstructure:"api"`
	AWS struct {
		Credentials struct {
			AccessKeyID     string `mapstructure:"accessKeyID,omitempty"`
			SecretAccessKey string `mapstructure:"secretAccessKey,omitempty"`
			Region          string `mapstructure:"region,omitempty"`
		} `mapstructure:"credentials"`
		S3 struct {
			BucketName string `mapstructure:"bucketName"`
			Endpoint   string `mapstructure:"endpoint,omitempty"`
		} `mapstructure:"s3"`
		SNS struct {
			TopicArn string `mapstructure:"topicArn"`
		} `mapstructure:"sns"`
	} `mapstructure:"aws"`
}
