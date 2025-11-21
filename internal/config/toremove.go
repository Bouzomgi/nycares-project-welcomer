package config

// Config represents the unified configuration structure for all lambdas
type Nope struct {
	// AWS configuration
	AWS struct {
		Credentials struct {
			AccessKeyID     string `mapstructure:"accessKeyID,omitempty"`
			SecretAccessKey string `mapstructure:"secretAccessKey,omitempty"`
			Region          string `mapstructure:"region,omitempty"`
		} `mapstructure:"credentials"`
		Dynamo struct {
			TableName string `mapstructure:"tableName"`
			Region    string `mapstructure:"region"`
			Endpoint  string `mapstructure:"endpoint,omitempty"`
		} `mapstructure:"dynamo"`
		S3 struct {
			BucketName string `mapstructure:"bucketName"`
			Endpoint   string `mapstructure:"endpoint,omitempty"`
		} `mapstructure:"s3"`
		SNS struct {
			TopicArn string `mapstructure:"topicArn"`
		} `mapstructure:"sns"`
		SF struct {
			CallbackEndpoint string `mapstructure:"callbackEndpoint"`
		} `mapstructure:"sf"`
	} `mapstructure:"aws"`

	// Auth configuration
	Account struct {
		Username   string `mapstructure:"username"`
		Password   string `mapstructure:"password"`
		InternalId string `mapstructure:"internalId"`
	} `mapstructure:"account"`
}
