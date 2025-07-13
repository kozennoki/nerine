package config

// Export private functions for testing
var GetEnvOrDefault = getEnvOrDefault

// Export private methods for testing
func (c *Config) ExportValidate() error {
	return c.validate()
}
