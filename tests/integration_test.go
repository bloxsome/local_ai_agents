package tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bloxsome/local_ai_agents/internal/config"
	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
	"github.com/bloxsome/local_ai_agents/internal/modules/ollama"
)

// TestSuite represents a test suite for Ollama integration
type TestSuite struct {
	client     *ollama.OllamaClient
	config     *config.Configuration
	logger     *logging.Logger
	ollamaURL  string
	testModel  string
	isOllamaUp bool
}

// SetupSuite sets up the test suite
func SetupSuite(t *testing.T) *TestSuite {
	// Set up logging
	logFile := "integration_test.log"
	if err := logging.SetupLogging(logFile, "DEBUG"); err != nil {
		t.Fatalf("Error setting up logging: %v", err)
	}
	logger := logging.GetLogger()
	logger.Info("Integration test suite started")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}

	// Set up Ollama client
	ollamaURL := "http://localhost:11434"
	client := ollama.NewOllamaClient(ollamaURL, 10*time.Second)

	// Check if Ollama is running
	isOllamaUp := isOllamaRunning(ollamaURL)
	if !isOllamaUp {
		logger.Warning("Ollama is not running at %s - some tests will be skipped", ollamaURL)
	} else {
		logger.Info("Ollama is running at %s", ollamaURL)
	}

	return &TestSuite{
		client:     client,
		config:     cfg,
		logger:     logger,
		ollamaURL:  ollamaURL,
		testModel:  "llama3.1:latest", // Use a model that should be available
		isOllamaUp: isOllamaUp,
	}
}

// TearDownSuite tears down the test suite
func (s *TestSuite) TearDownSuite() {
	// Clean up log file
	logFile := "integration_test.log"
	if err := os.Remove(logFile); err != nil {
		s.logger.Warning("Failed to remove log file: %v", err)
	}
}

// Helper function to check if Ollama is running
func isOllamaRunning(url string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound
}

// TestOllamaIsRunning tests if Ollama is running
func TestOllamaIsRunning(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite()

	if !suite.isOllamaUp {
		t.Skip("Ollama is not running - skipping test")
	}

	// If we got here, Ollama is running
	t.Log("Ollama is running at", suite.ollamaURL)
}

// TestOllamaPrompt tests sending a prompt to Ollama
func TestOllamaPrompt(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite()

	if !suite.isOllamaUp {
		t.Skip("Ollama is not running - skipping test")
	}

	// Test sending a prompt
	prompt := "Hello, this is a test from Local AI Agents. Please respond with a short greeting."
	username := "TestUser"

	suite.logger.Info("Sending test prompt to Ollama")
	response, err := suite.client.ProcessPrompt(prompt, suite.testModel, username)

	if err != nil {
		t.Fatalf("Error sending prompt to Ollama: %v", err)
	}

	// Check if we got a non-empty response
	if response == "" {
		t.Errorf("Received empty response from Ollama")
	} else {
		suite.logger.Info("Response: %s", response)
		fmt.Printf("Ollama response: %s\n", response)
	}
}

// TestOllamaErrorHandling tests error handling when Ollama is not available
func TestOllamaErrorHandling(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite()

	// Create a client with an invalid URL
	invalidClient := ollama.NewOllamaClient("http://localhost:99999", 2*time.Second)

	// Test sending a prompt to the invalid URL
	prompt := "This should fail because the URL is invalid."
	username := "TestUser"

	suite.logger.Info("Sending test prompt to invalid Ollama URL")
	_, err := invalidClient.ProcessPrompt(prompt, suite.testModel, username)

	// We expect an error
	if err == nil {
		t.Errorf("Expected error when connecting to invalid URL, but got none")
	} else {
		suite.logger.Info("Received expected error: %v", err)
	}
}

// TestOllamaWithConfig tests using Ollama with the local_ai_agents configuration
func TestOllamaWithConfig(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite()

	if !suite.isOllamaUp {
		t.Skip("Ollama is not running - skipping test")
	}

	// Test that the configuration is loaded correctly
	if suite.config == nil {
		t.Fatalf("Configuration is nil")
	}

	// Check that the agent name is set
	if suite.config.AgentName == "" {
		t.Errorf("Agent name is empty")
	} else {
		t.Logf("Agent name: %s", suite.config.AgentName)
	}

	// Check that the default model is set
	if suite.config.DefaultModel == "" {
		t.Errorf("Default model is empty")
	} else {
		t.Logf("Default model: %s", suite.config.DefaultModel)
	}

	// Test sending a prompt using the configuration
	prompt := "Hello, this is a test using the configuration."
	username := suite.config.UserName

	suite.logger.Info("Sending test prompt to Ollama using configuration")
	response, err := suite.client.ProcessPrompt(prompt, suite.config.DefaultModel, username)

	if err != nil {
		t.Fatalf("Error sending prompt to Ollama: %v", err)
	}

	// Check if we got a non-empty response
	if response == "" {
		t.Errorf("Received empty response from Ollama")
	} else {
		suite.logger.Info("Response: %s", response)
		fmt.Printf("Ollama response: %s\n", response)
	}
}
