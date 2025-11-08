package main

type Config struct {
	Account struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"account"`
}
