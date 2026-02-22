package login

type Config struct {
	Api struct {
		BaseUrl string `mapstructure:"base_url,omitempty"`
	} `mapstructure:"api"`
	Account struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"account"`
}
