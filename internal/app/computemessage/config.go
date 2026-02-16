package computemessage

type Config struct {
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
		} `mapstructure:"s3"`
	} `mapstructure:"aws"`
	CurrentDate string `mapstructure:"currentDate,omitempty"`
}
