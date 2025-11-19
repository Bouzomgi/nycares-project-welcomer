package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig[T any]() (*T, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvPrefix("NYCARES")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Detect Lambda
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda: read only from env vars
	} else {
		// Local dev: read YAML config file
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".") // local config path

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("error reading config file: %w", err)
			}
			// File not found is ok, continue with env vars
		}
	}

	var cfg T
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate all fields
	if err := validateStruct(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// validateStruct recursively checks all string fields in a struct
func validateStruct(s interface{}) error {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				return fmt.Errorf("missing required field: %s", fieldType.Name)
			}
		case reflect.Struct:
			if err := validateStruct(field.Interface()); err != nil {
				return fmt.Errorf("%s.%w", fieldType.Name, err)
			}
			// Optionally handle slices/maps if needed
		}
	}
	return nil
}
