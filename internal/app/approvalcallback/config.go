package approvalcallback

type Config struct {
	AWS struct {
		Credentials struct {
			AccessKeyID     string `mapstructure:"accessKeyID,omitempty"`
			SecretAccessKey string `mapstructure:"secretAccessKey,omitempty"`
			Region          string `mapstructure:"region,omitempty"`
		} `mapstructure:"credentials"`
		SF struct {
			ApprovalSecret string `mapstructure:"approvalSecret"`
		} `mapstructure:"sf"`
	} `mapstructure:"aws"`
}
