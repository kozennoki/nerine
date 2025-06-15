package config

import (
	"errors"
	"os"
)

type Config struct {
	Port              string
	MicroCMSAPIKey    string
	MicroCMSServiceID string
	NerineAPIKey      string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:              getEnvOrDefault("PORT", "8080"),
		MicroCMSAPIKey:    os.Getenv("MICROCMS_API_KEY"),
		MicroCMSServiceID: os.Getenv("MICROCMS_SERVICE_ID"),
		NerineAPIKey:      os.Getenv("NERINE_API_KEY"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.MicroCMSAPIKey == "" {
		return errors.New("MICROCMS_API_KEY is required")
	}
	if c.MicroCMSServiceID == "" {
		return errors.New("MICROCMS_SERVICE_ID is required")
	}
	if c.NerineAPIKey == "" {
		return errors.New("NERINE_API_KEY is required")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
