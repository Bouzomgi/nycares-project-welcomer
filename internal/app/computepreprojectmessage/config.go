package computepreprojectmessage

type Config struct {
	AWS struct {
		Credentials struct {
			AccessKeyID     string `mapstructure:"accessKeyID,omitempty"`
			SecretAccessKey string `mapstructure:"secretAccessKey,omitempty"`
			Region          string `mapstructure:"region,omitempty"`
		} `mapstructure:"credentials"`
		S3 struct {
			BucketName string `mapstructure:"bucketName"`
		} `mapstructure:"s3"`
	} `mapstructure:"aws"`
}
