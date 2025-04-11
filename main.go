package main

import (
	"fmt"
	"os"

	"github.com/bloxsome/local_ai_agents/internal/config"
	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
	"github.com/bloxsome/local_ai_agents/internal/modules/mcp"
)

func main() {
	// Set up logging
	logFile := "local_ai_agents.log"
	if err := logging.SetupLogging(logFile, "INFO"); err != nil {
		fmt.Printf("Error setting up logging: %v\n", err)
		os.Exit(1)
	}
	logger := logging.GetLogger()
	logger.Info("Local AI Agents application started")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Error loading configuration: %v", err)
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Welcome to Local AI Agents! Running as %s\n", cfg.AgentName)

	// Initialize MCP if enabled
	if cfg.MCPEnabled {
		logger.Info("Initializing MCP module...")
		if err := mcp.Initialize(cfg); err != nil {
			logger.Error("Error initializing MCP module: %v", err)
			fmt.Printf("Error initializing MCP module: %v\n", err)
		} else {
			logger.Info("MCP module initialized successfully")
			defer mcp.Shutdown()
		}
	}

	// TODO: Initialize and run the agent
	// This is a placeholder for the actual agent initialization and execution
	// In a complete implementation, this would:
	// 1. Set up the agent with the loaded configuration
	// 2. Start the main interaction loop
	// 3. Handle user input and agent responses

	logger.Info("Local AI Agents application completed")
}
