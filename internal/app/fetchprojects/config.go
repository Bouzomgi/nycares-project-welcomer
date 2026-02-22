package fetchprojects

type Config struct {
	Api struct {
		BaseUrl string `mapstructure:"base_url,omitempty"`
	} `mapstructure:"api"`
	Account struct {
		Username   string `mapstructure:"username"`
		Password   string `mapstructure:"password"`
		InternalId string `mapstructure:"internalId"`
	} `mapstructure:"account"`
}
