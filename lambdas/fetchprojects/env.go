package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
    Account struct {
        Username   string `mapstructure:"username"`
        Password   string `mapstructure:"password"`
        InternalId string `mapstructure:"internalId"`
    } `mapstructure:"account"`
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")  // name of config file (without extension)
    viper.SetConfigType("yaml")    // config file type
    viper.AddConfigPath(".")       // look for config in the current directory
    viper.AutomaticEnv()           // allow environment variables to override
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // 🔹 Actually read the config file from disk
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("error reading config file: %w", err)
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("error unmarshalling config: %w", err)
    }

    return &cfg, nil
}
