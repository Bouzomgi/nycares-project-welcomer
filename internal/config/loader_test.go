package config

import (
	"os"
	"strings"
	"testing"
)

// testConfig is a minimal config struct used in loader tests.
type testConfig struct {
	Username string `mapstructure:"username"`
	Region   string `mapstructure:"region,omitempty"`
}

// nestedTestConfig is used to test nested struct validation.
type nestedTestConfig struct {
	Inner testConfig `mapstructure:"inner"`
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		wantErr     bool
		errContains string
	}{
		{
			name:    "all required fields set passes",
			input:   testConfig{Username: "alice", Region: ""},
			wantErr: false,
		},
		{
			name:        "missing required field returns error with field name",
			input:       testConfig{Username: "", Region: ""},
			wantErr:     true,
			errContains: "Username",
		},
		{
			name:    "omitempty field passes when empty",
			input:   testConfig{Username: "alice", Region: ""},
			wantErr: false,
		},
		{
			name:    "omitempty field passes when set",
			input:   testConfig{Username: "bob", Region: "us-east-1"},
			wantErr: false,
		},
		{
			name:        "nested struct with missing required field returns error",
			input:       nestedTestConfig{Inner: testConfig{Username: "", Region: ""}},
			wantErr:     true,
			errContains: "Username",
		},
		{
			name:    "nested struct with all required fields passes",
			input:   nestedTestConfig{Inner: testConfig{Username: "alice", Region: ""}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStruct(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Ensure we're running in Lambda mode so Viper uses env vars instead of config.yaml.
	origFnName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	origSSMPath := os.Getenv("NYCARES_SSM_PATH")
	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test-fn")
	os.Setenv("NYCARES_SSM_PATH", "") // skip SSM loading

	os.Setenv("NYCARES_USERNAME", "testuser")
	os.Setenv("NYCARES_REGION", "us-west-2")

	t.Cleanup(func() {
		os.Setenv("AWS_LAMBDA_FUNCTION_NAME", origFnName)
		os.Setenv("NYCARES_SSM_PATH", origSSMPath)
		os.Unsetenv("NYCARES_USERNAME")
		os.Unsetenv("NYCARES_REGION")
	})

	cfg, err := LoadConfig[testConfig]()
	if err != nil {
		t.Fatalf("LoadConfig returned unexpected error: %v", err)
	}

	if cfg.Username != "testuser" {
		t.Errorf("Username = %q, want %q", cfg.Username, "testuser")
	}
	if cfg.Region != "us-west-2" {
		t.Errorf("Region = %q, want %q", cfg.Region, "us-west-2")
	}
}

func TestLoadConfigFromEnv_MissingRequired(t *testing.T) {
	origFnName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	origSSMPath := os.Getenv("NYCARES_SSM_PATH")
	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test-fn")
	os.Setenv("NYCARES_SSM_PATH", "")

	// Deliberately do NOT set NYCARES_USERNAME.
	os.Unsetenv("NYCARES_USERNAME")
	os.Unsetenv("NYCARES_REGION")

	t.Cleanup(func() {
		os.Setenv("AWS_LAMBDA_FUNCTION_NAME", origFnName)
		os.Setenv("NYCARES_SSM_PATH", origSSMPath)
	})

	_, err := LoadConfig[testConfig]()
	if err == nil {
		t.Fatal("expected error for missing required field, got nil")
	}
	if !strings.Contains(err.Error(), "Username") {
		t.Errorf("error %q should mention missing field Username", err.Error())
	}
}
