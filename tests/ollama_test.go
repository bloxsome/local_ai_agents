package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/bloxsome/local_ai_agents/internal/config"
	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
	"github.com/bloxsome/local_ai_agents/internal/modules/mcp"
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
	// Set up logging
	logFile := "test_ollama_models.log"
	if err := logging.SetupLogging(logFile, "DEBUG"); err != nil {
		t.Fatalf("Error setting up logging: %v", err)
	}
	logger := logging.GetLogger()
	logger.Info("Ollama model list test started")

	// Create HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send request to Ollama API to list models
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		t.Fatalf("Error connecting to Ollama API: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ollama API returned non-200 status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	// Parse response
	var response struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Error parsing response: %v", err)
	}

	// Check if we got any models
	if len(response.Models) == 0 {
		t.Errorf("No models found")
	} else {
		logger.Info("Found %d models", len(response.Models))

		// Print models
		fmt.Println("Available Ollama models:")
		for _, model := range response.Models {
			fmt.Printf("- %s\n", model.Name)
			logger.Info("Model: %s", model.Name)
		}

		// Check if our test model is available
		testModel := "llama3.3:latest"
		found := false
		for _, model := range response.Models {
			if model.Name == testModel {
				found = true
				break
			}
		}

		if !found {
			t.Logf("Test model %s not found, but test can continue", testModel)
		} else {
			t.Logf("Test model %s is available", testModel)
		}
	}

	// Clean up log file
	if err := os.Remove(logFile); err != nil {
		logger.Warning("Failed to remove log file: %v", err)
	}
}

// TestOllamaWithMCP tests the integration between Ollama and MCP
func TestOllamaWithMCP(t *testing.T) {
	// Set up logging
	logFile := "test_ollama_mcp.log"
	if err := logging.SetupLogging(logFile, "DEBUG"); err != nil {
		t.Fatalf("Error setting up logging: %v", err)
	}
	logger := logging.GetLogger()
	logger.Info("Ollama MCP integration test started")

	// Start the weather MCP server
	logger.Info("Starting weather MCP server")
	weatherServerPath := filepath.Join(".", "mcp_server", "weather")
	cmd := exec.Command(weatherServerPath)

	// Set up pipes for stdin, stdout, stderr
	_, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}

	_, err = cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	_, err = cmd.StderrPipe()
	if err != nil {
		t.Fatalf("Failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start weather MCP server: %v", err)
	}

	// Defer stopping the server
	defer func() {
		if err := cmd.Process.Kill(); err != nil {
			logger.Error("Failed to kill weather MCP server: %v", err)
		}
		cmd.Wait()
	}()

	// Wait a moment for the server to start
	time.Sleep(1 * time.Second)

	// Create Ollama client with a timeout
	client := ollama.NewOllamaClient("http://localhost:11434", 30*time.Second)

	// Create a prompt with MCP command syntax to get weather for a city
	prompt := `Hello, this is a test from Local AI Agents.
	
<use_mcp_tool>
<server_name>weather</server_name>
<tool_name>get_weather</tool_name>
<arguments>
{
  "city": "London"
}
</arguments>
</use_mcp_tool>

Please respond with a short greeting.`

	// Initialize MCP module with a mock configuration
	mockConfig := &config.Configuration{
		MCPEnabled: true,
		MCPDir:     ".",
	}

	// Initialize MCP module
	if err := mcp.Initialize(mockConfig); err != nil {
		logger.Warning("Failed to initialize MCP module: %v", err)
	}

	// Add weather server to MCP manager
	if mcp.DefaultManager != nil {
		if err := mcp.DefaultManager.AddServer("weather", weatherServerPath, []string{}, nil); err != nil {
			logger.Warning("Failed to add weather server: %v", err)
		} else {
			// Start weather server
			if err := mcp.DefaultManager.StartServer("weather"); err != nil {
				logger.Warning("Failed to start weather server: %v", err)
			}
		}
	}

	// Process the prompt with MCP
	processedPrompt, err := mcp.ProcessMCPInPrompt(prompt)
	if err != nil {
		t.Fatalf("Error processing MCP in prompt: %v", err)
	}

	// Log the processed prompt
	logger.Info("Processed prompt: %s", processedPrompt)

	// Now test sending the processed prompt to Ollama
	model := "llama3.3:latest"
	username := "TestUser"

	logger.Info("Sending processed prompt to Ollama")
	response, err := client.ProcessPrompt(processedPrompt, model, username)

	if err != nil {
		t.Fatalf("Error connecting to Ollama: %v", err)
	}

	logger.Info("Received response from Ollama")

	// Check if we got a non-empty response
	if response == "" {
		t.Errorf("Received empty response from Ollama")
	} else {
		logger.Info("Response: %s", response)
		fmt.Printf("Ollama response to MCP prompt: %s\n", response)
	}

	// Shutdown MCP module
	mcp.Shutdown()

	// Clean up log file
	if err := os.Remove(logFile); err != nil {
		logger.Warning("Failed to remove log file: %v", err)
	}
}
