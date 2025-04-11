// This file is deprecated. Configuration has been moved to internal/config/config.go
// This file is kept for backward compatibility and will be removed in a future version.
package main

import (
	"github.com/bloxsome/local_ai_agents/internal/config"
)

// Configuration is an alias for config.Configuration
type Configuration = config.Configuration

// LoadConfig loads the configuration from environment variables and defaults
// This function is kept for backward compatibility
func LoadConfig() (*Configuration, error) {
	return config.LoadConfig()
}
