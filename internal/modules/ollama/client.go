// Package ollama provides a client for interacting with the Ollama API
package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bloxsome/local_ai_agents/internal/errors"
	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
	"github.com/bloxsome/local_ai_agents/internal/modules/mcp"
)

// OllamaClient is a client for interacting with the Ollama API
type OllamaClient struct {
	BaseURL string
	Timeout time.Duration
	logger  *logging.Logger
}

// GenerateResponse represents the response from the Ollama API
type GenerateResponse struct {
	Model         string `json:"model"`
	Response      string `json:"response"`
	Done          bool   `json:"done"`
	Context       []int  `json:"context,omitempty"`
	TotalDuration int64  `json:"total_duration,omitempty"`
}

// NewOllamaClient creates a new OllamaClient
func NewOllamaClient(baseURL string, timeout time.Duration) *OllamaClient {
	return &OllamaClient{
		BaseURL: baseURL,
		Timeout: timeout,
		logger:  logging.GetLogger(),
	}
}

// DefaultClient is the default OllamaClient instance
var DefaultClient = NewOllamaClient("http://localhost:11434", 60*time.Second)

// ProcessPrompt sends a prompt to the Ollama API and returns the response
func (c *OllamaClient) ProcessPrompt(prompt, model, username string) (string, error) {
	c.logger.Info("Processing prompt for user: %s, model: %s", username, model)

	// Process MCP commands in the prompt
	processedPrompt, err := mcp.ProcessMCPInPrompt(prompt)
	if err != nil {
		c.logger.Warning("Failed to process MCP commands: %v", err)
		// Continue with the original prompt
		processedPrompt = prompt
	}

	url := fmt.Sprintf("%s/api/generate", c.BaseURL)
	data := map[string]string{
		"model":  model,
		"prompt": processedPrompt,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.logger.Error("Failed to marshal JSON: %v", err)
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to create request: %v", err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: c.Timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		c.logger.Error("Failed to send request: %v", err)
		return "", errors.NewAPIConnectionError(fmt.Sprintf("failed to send request: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("API returned non-200 status code: %d", resp.StatusCode)
		return "", errors.NewAPIConnectionError(fmt.Sprintf("API returned status code %d", resp.StatusCode))
	}

	// Stream the response
	fullResponse, err := c.streamResponse(resp.Body)
	if err != nil {
		c.logger.Error("Failed to stream response: %v", err)
		return "", fmt.Errorf("failed to stream response: %w", err)
	}

	c.logger.Info("Response generated for prompt: %s...", truncateString(prompt, 50))

	// Save interaction (in a real implementation, this would save to a file or database)
	c.logger.Info("Saving interaction for user: %s", username)

	return fullResponse, nil
}

// streamResponse streams the response from the Ollama API
func (c *OllamaClient) streamResponse(body io.Reader) (string, error) {
	scanner := bufio.NewScanner(body)
	var fullResponse strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var response GenerateResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			c.logger.Warning("Failed to decode JSON from line: %s", line)
			continue
		}

		if response.Response != "" {
			fullResponse.WriteString(response.Response)
			fmt.Print(response.Response) // Print the response chunk
		}

		if response.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return fullResponse.String(), nil
}

// ProcessPrompt sends a prompt to the Ollama API using the default client
func ProcessPrompt(prompt, model, username string) (string, error) {
	return DefaultClient.ProcessPrompt(prompt, model, username)
}

// Generate is an alias for ProcessPrompt
func Generate(prompt, model, username string) (string, error) {
	return ProcessPrompt(prompt, model, username)
}

// Helper function to truncate a string
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
