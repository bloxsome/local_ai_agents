// Package main provides an example of how to initialize and use the MCP module with Ollama
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bloxsome/local_ai_agents/internal/config"
	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
	"github.com/bloxsome/local_ai_agents/internal/modules/mcp"
	"github.com/bloxsome/local_ai_agents/internal/modules/ollama"
)

func main() {
	// Set up logging
	logFile := "mcp_example.log"
	if err := logging.SetupLogging(logFile, "DEBUG"); err != nil {
		fmt.Printf("Error setting up logging: %v\n", err)
		os.Exit(1)
	}
	logger := logging.GetLogger()
	logger.Info("MCP initialization example started")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Error loading configuration: %v", err)
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize MCP module
	logger.Info("Initializing MCP module...")
	if err := mcp.Initialize(cfg); err != nil {
		logger.Error("Error initializing MCP module: %v", err)
		fmt.Printf("Error initializing MCP module: %v\n", err)
		os.Exit(1)
	}
	logger.Info("MCP module initialized")

	// Add browser MCP server
	browserServerPath := filepath.Join(cfg.ProjectRoot, "local_ai_agents", "examples", "mcp_browser_server", "browser")
	logger.Info("Adding browser MCP server: %s", browserServerPath)
	if err := mcp.DefaultManager.AddServer("browser", browserServerPath, []string{}, nil); err != nil {
		logger.Error("Error adding browser MCP server: %v", err)
		fmt.Printf("Error adding browser MCP server: %v\n", err)
		os.Exit(1)
	}
	logger.Info("Browser MCP server added")

	// Start browser MCP server
	logger.Info("Starting browser MCP server...")
	if err := mcp.DefaultManager.StartServer("browser"); err != nil {
		logger.Error("Error starting browser MCP server: %v", err)
		fmt.Printf("Error starting browser MCP server: %v\n", err)
		os.Exit(1)
	}
	logger.Info("Browser MCP server started")

	// List available MCP servers
	servers := mcp.DefaultManager.ListServers()
	fmt.Println("Available MCP servers:")
	for _, server := range servers {
		running, _ := mcp.DefaultManager.IsServerRunning(server.Name)
		status := "stopped"
		if running {
			status = "running"
		}
		fmt.Printf("- %s (%s)\n", server.Name, status)
	}

	// Create a prompt with MCP tool usage
	testPagePath := filepath.Join(cfg.ProjectRoot, "local_ai_agents", "examples", "test_page.html")
	testPageURL := fmt.Sprintf("file://%s", testPagePath)
	prompt := fmt.Sprintf(`
Open the test page in my browser.

<use_mcp_tool>
<server_name>browser</server_name>
<tool_name>open_url</tool_name>
<arguments>
{
  "url": "%s"
}
</arguments>
</use_mcp_tool>
`, testPageURL)

	// Process MCP commands in the prompt
	logger.Info("Processing MCP commands in prompt...")
	processedPrompt, err := mcp.ProcessMCPInPrompt(prompt)
	if err != nil {
		logger.Error("Error processing MCP commands: %v", err)
		fmt.Printf("Error processing MCP commands: %v\n", err)
	} else {
		logger.Info("MCP commands processed successfully")
		fmt.Println("\nOriginal prompt:")
		fmt.Println(prompt)
		fmt.Println("\nProcessed prompt:")
		fmt.Println(processedPrompt)
	}

	// Send the processed prompt to Ollama
	logger.Info("Sending processed prompt to Ollama...")
	response, err := ollama.ProcessPrompt(processedPrompt, "llama3.3", "user")
	if err != nil {
		logger.Error("Error sending prompt to Ollama: %v", err)
		fmt.Printf("Error sending prompt to Ollama: %v\n", err)
	} else {
		logger.Info("Prompt sent to Ollama successfully")
		fmt.Println("\nOllama response:")
		fmt.Println(response)
	}

	// Stop browser MCP server
	logger.Info("Stopping browser MCP server...")
	if err := mcp.DefaultManager.StopServer("browser"); err != nil {
		logger.Error("Error stopping browser MCP server: %v", err)
		fmt.Printf("Error stopping browser MCP server: %v\n", err)
	}
	logger.Info("Browser MCP server stopped")

	// Shutdown MCP module
	logger.Info("Shutting down MCP module...")
	mcp.Shutdown()
	logger.Info("MCP module shut down")

	logger.Info("MCP initialization example completed")
}
