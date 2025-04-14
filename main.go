package main

import (
	"fmt"
	"os"

	"github.com/bloxsome/local_ai_agents/internal/config"
	"github.com/bloxsome/local_ai_agents/internal/modules/context"
	"github.com/bloxsome/local_ai_agents/internal/modules/input"
	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
	"github.com/bloxsome/local_ai_agents/internal/modules/mcp"
	"github.com/bloxsome/local_ai_agents/internal/modules/ollama"
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

	// Initialize context module
	logger.Info("Initializing context module...")
	if err := context.Initialize(cfg); err != nil {
		logger.Error("Error initializing context module: %v", err)
		fmt.Printf("Error initializing context module: %v\n", err)
		os.Exit(1)
	} else {
		logger.Info("Context module initialized successfully")
		defer context.Shutdown()
	}

	// Initialize and run the agent
	logger.Info("Starting agent interaction loop...")
	fmt.Println("Type your messages to interact with the agent.")
	fmt.Println("Type /help to see available commands.")

	// Create input handler
	inputHandler := input.NewInputHandler(cfg)
	defer inputHandler.Close()

	// Main interaction loop
	for {
		// Get user input
		result := inputHandler.GetUserInput()

		// Check if we should exit
		if result.Exit {
			break
		}

		// Skip if this was a command that was already handled
		if result.Skip {
			continue
		}

		// Add user message to context
		context.AddMessage("user", result.Text)

		// Process the input with Ollama
		logger.Info("Processing user input with Ollama")

		// Format prompt with context
		contextStr := context.FormatContextForPrompt()
		prompt := contextStr + "\nUser: " + result.Text + "\n\nAssistant: "

		// Generate response
		response, err := ollama.Generate(prompt, cfg.DefaultModel, cfg.UserName)
		if err != nil {
			logger.Error("Error processing prompt: %v", err)
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Add assistant message to context
		context.AddMessage("assistant", response)

		// Save history periodically
		if err := context.SaveHistory(); err != nil {
			logger.Warning("Failed to save history: %v", err)
		}
	}

	logger.Info("Local AI Agents application completed")
}
