package config_test

import (
	"os"
	"testing"

	"github.com/kozennoki/nerine/internal/infrastructure/config"
)

func TestLoad_Success(t *testing.T) {

	// Setup environment variables
	os.Setenv("MICROCMS_API_KEY", "test-microcms-key")
	os.Setenv("MICROCMS_SERVICE_ID", "test-service-id")
	os.Setenv("NERINE_API_KEY", "test-nerine-key")
	os.Setenv("PORT", "9000")

	defer func() {
		os.Unsetenv("MICROCMS_API_KEY")
		os.Unsetenv("MICROCMS_SERVICE_ID")
		os.Unsetenv("NERINE_API_KEY")
		os.Unsetenv("PORT")
	}()

	cfg, err := config.Load()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if cfg.MicroCMSAPIKey != "test-microcms-key" {
		t.Errorf("Expected MicroCMSAPIKey to be 'test-microcms-key', got: %s", cfg.MicroCMSAPIKey)
	}

	if cfg.MicroCMSServiceID != "test-service-id" {
		t.Errorf("Expected MicroCMSServiceID to be 'test-service-id', got: %s", cfg.MicroCMSServiceID)
	}

	if cfg.NerineAPIKey != "test-nerine-key" {
		t.Errorf("Expected NerineAPIKey to be 'test-nerine-key', got: %s", cfg.NerineAPIKey)
	}

	if cfg.Port != "9000" {
		t.Errorf("Expected Port to be '9000', got: %s", cfg.Port)
	}
}

func TestLoad_DefaultPort(t *testing.T) {

	// Setup required environment variables without PORT
	os.Setenv("MICROCMS_API_KEY", "test-microcms-key")
	os.Setenv("MICROCMS_SERVICE_ID", "test-service-id")
	os.Setenv("NERINE_API_KEY", "test-nerine-key")
	os.Unsetenv("PORT")

	defer func() {
		os.Unsetenv("MICROCMS_API_KEY")
		os.Unsetenv("MICROCMS_SERVICE_ID")
		os.Unsetenv("NERINE_API_KEY")
	}()

	cfg, err := config.Load()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("Expected default Port to be '8080', got: %s", cfg.Port)
	}
}

func TestLoad_MissingMicroCMSAPIKey(t *testing.T) {

	// Setup environment variables without MICROCMS_API_KEY
	os.Unsetenv("MICROCMS_API_KEY")
	os.Setenv("MICROCMS_SERVICE_ID", "test-service-id")
	os.Setenv("NERINE_API_KEY", "test-nerine-key")

	defer func() {
		os.Unsetenv("MICROCMS_SERVICE_ID")
		os.Unsetenv("NERINE_API_KEY")
	}()

	cfg, err := config.Load()

	if err == nil {
		t.Fatal("Expected error for missing MICROCMS_API_KEY, got nil")
	}

	if cfg != nil {
		t.Error("Expected cfg to be nil when validation fails")
	}

	expectedErr := "MICROCMS_API_KEY is required"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', got: %s", expectedErr, err.Error())
	}
}

func TestLoad_MissingMicroCMSServiceID(t *testing.T) {

	// Setup environment variables without MICROCMS_SERVICE_ID
	os.Setenv("MICROCMS_API_KEY", "test-microcms-key")
	os.Unsetenv("MICROCMS_SERVICE_ID")
	os.Setenv("NERINE_API_KEY", "test-nerine-key")

	defer func() {
		os.Unsetenv("MICROCMS_API_KEY")
		os.Unsetenv("NERINE_API_KEY")
	}()

	cfg, err := config.Load()

	if err == nil {
		t.Fatal("Expected error for missing MICROCMS_SERVICE_ID, got nil")
	}

	if cfg != nil {
		t.Error("Expected cfg to be nil when validation fails")
	}

	expectedErr := "MICROCMS_SERVICE_ID is required"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', got: %s", expectedErr, err.Error())
	}
}

func TestLoad_MissingNerineAPIKey(t *testing.T) {

	// Setup environment variables without NERINE_API_KEY
	os.Setenv("MICROCMS_API_KEY", "test-microcms-key")
	os.Setenv("MICROCMS_SERVICE_ID", "test-service-id")
	os.Unsetenv("NERINE_API_KEY")

	defer func() {
		os.Unsetenv("MICROCMS_API_KEY")
		os.Unsetenv("MICROCMS_SERVICE_ID")
	}()

	cfg, err := config.Load()

	if err == nil {
		t.Fatal("Expected error for missing NERINE_API_KEY, got nil")
	}

	if cfg != nil {
		t.Error("Expected cfg to be nil when validation fails")
	}

	expectedErr := "NERINE_API_KEY is required"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', got: %s", expectedErr, err.Error())
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Environment variable exists",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "Environment variable does not exist",
			key:          "NON_EXISTENT_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "Environment variable is empty string",
			key:          "EMPTY_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := config.GetEnvOrDefault(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestConfig_validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		config      config.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid config",
			config: config.Config{
				Port:              "8080",
				MicroCMSAPIKey:    "test-key",
				MicroCMSServiceID: "test-service",
				NerineAPIKey:      "test-nerine",
			},
			expectError: false,
		},
		{
			name: "Missing MicroCMSAPIKey",
			config: config.Config{
				Port:              "8080",
				MicroCMSAPIKey:    "",
				MicroCMSServiceID: "test-service",
				NerineAPIKey:      "test-nerine",
			},
			expectError: true,
			errorMsg:    "MICROCMS_API_KEY is required",
		},
		{
			name: "Missing MicroCMSServiceID",
			config: config.Config{
				Port:              "8080",
				MicroCMSAPIKey:    "test-key",
				MicroCMSServiceID: "",
				NerineAPIKey:      "test-nerine",
			},
			expectError: true,
			errorMsg:    "MICROCMS_SERVICE_ID is required",
		},
		{
			name: "Missing NerineAPIKey",
			config: config.Config{
				Port:              "8080",
				MicroCMSAPIKey:    "test-key",
				MicroCMSServiceID: "test-service",
				NerineAPIKey:      "",
			},
			expectError: true,
			errorMsg:    "NERINE_API_KEY is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.config.ExportValidate()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got: %s", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}
