package generatethankyoumessage

type Config struct {
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
		Bedrock struct {
			Endpoint string `mapstructure:"endpoint,omitempty"`
		} `mapstructure:"bedrock"`
	} `mapstructure:"aws"`
	Mock struct {
		GenerateThankYou bool `mapstructure:"generateThankYou,omitempty"`
	} `mapstructure:"mock"`
}
