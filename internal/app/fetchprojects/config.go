package fetchprojects

type Config struct {
	Api struct {
		BaseUrl string `mapstructure:"base_url,omitempty"`
	} `mapstructure:"api"`
}
