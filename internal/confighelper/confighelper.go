package confighelper

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Generic config loader that works with any struct type.
func LoadConfig[T any]() (*T, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // config file type
	viper.AddConfigPath(".")      // look for config in the current directory
	viper.AutomaticEnv()          // allow environment variables to override
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg T
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &cfg, nil
}
