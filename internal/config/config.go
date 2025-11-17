package config

// Config represents the unified configuration structure for all lambdas
type Config struct {
	// AWS configuration
	AWS struct {
		Credentials struct {
			AccessKeyID     string `mapstructure:"accessKeyID"`
			SecretAccessKey string `mapstructure:"secretAccessKey"`
			Region          string `mapstructure:"region"`
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

	// Environment information
	Env string `mapstructure:"env"`
}
