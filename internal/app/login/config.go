package login

// Config represents the unified configuration structure for all lambdas
type Config struct {
	Account struct {
		Username   string `mapstructure:"username"`
		Password   string `mapstructure:"password"`
		InternalId string `mapstructure:"internalId"`
	} `mapstructure:"account"`
}
