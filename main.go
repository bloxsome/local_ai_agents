package main

import (
	"fmt"
	"os"

	"github.com/bloxsome/local_ai_agents/modules"
)

func main() {
	// Set up logging
	logFile := "local_ai_agents.log"
	if err := modules.SetupLogging(logFile, "INFO"); err != nil {
		fmt.Printf("Error setting up logging: %v\n", err)
		os.Exit(1)
	}
	logger := modules.GetLogger()
	logger.Info("Local AI Agents application started")

	// Load configuration
	cfg, err := LoadConfig()
	if err != nil {
		logger.Error("Error loading configuration: %v", err)
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Welcome to Local AI Agents! Running as %s\n", cfg.AgentName)

	// TODO: Initialize and run the agent
	// This is a placeholder for the actual agent initialization and execution
	// In a complete implementation, this would:
	// 1. Set up the agent with the loaded configuration
	// 2. Initialize MCP if enabled
	// 3. Start the main interaction loop
	// 4. Handle user input and agent responses

	logger.Info("Local AI Agents application completed")
}
