package main

type Config struct {
	Account struct {
		InternalId string `mapstructure:"internalId"`
	} `mapstructure:"account"`
}
