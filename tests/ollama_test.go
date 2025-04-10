package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
	"github.com/bloxsome/local_ai_agents/internal/modules/ollama"
)

// TestOllamaConnection tests the connection to the Ollama API
func TestOllamaConnection(t *testing.T) {
	// Set up logging
	logFile := "test_ollama.log"
	if err := logging.SetupLogging(logFile, "DEBUG"); err != nil {
		t.Fatalf("Error setting up logging: %v", err)
	}
	logger := logging.GetLogger()
	logger.Info("Ollama connection test started")

	// Create Ollama client with a timeout
	client := ollama.NewOllamaClient("http://localhost:11434", 30*time.Second)

	// Test if Ollama is running by sending a simple prompt
	prompt := "Hello, this is a test from Local AI Agents. Please respond with a short greeting."
	model := "llama3.3:latest" // Use a model that should be available
	username := "TestUser"

	logger.Info("Sending test prompt to Ollama")
	response, err := client.ProcessPrompt(prompt, model, username)

	if err != nil {
		t.Fatalf("Error connecting to Ollama: %v", err)
	}

	logger.Info("Received response from Ollama")

	// Check if we got a non-empty response
	if response == "" {
		t.Errorf("Received empty response from Ollama")
	} else {
		logger.Info("Response: %s", response)
		fmt.Printf("Ollama response: %s\n", response)
	}

	// Clean up log file
	if err := os.Remove(logFile); err != nil {
		logger.Warning("Failed to remove log file: %v", err)
	}
}

// TestOllamaModelList tests listing available models from Ollama
func TestOllamaModelList(t *testing.T) {
	// This is a placeholder for a future test that would list available models
	// Ollama API doesn't directly expose this functionality through our client yet
	// This would require extending the client to support the /api/tags endpoint

	t.Skip("Model listing functionality not implemented yet")
}

// TestOllamaWithMCP tests the integration between Ollama and MCP
func TestOllamaWithMCP(t *testing.T) {
	// This is a placeholder for a future test that would test MCP integration
	// It would require setting up MCP servers and testing the integration

	t.Skip("MCP integration test not implemented yet")
}
