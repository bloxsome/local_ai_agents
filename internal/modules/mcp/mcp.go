// Package mcp provides functionality for Model Context Protocol (MCP) support
package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bloxsome/local_ai_agents/internal/config"
	"github.com/bloxsome/local_ai_agents/internal/errors"
	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
)

// MCPServer represents a Model Context Protocol server
type MCPServer struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	Env         map[string]string `json:"env"`
	Disabled    bool              `json:"disabled"`
	AutoApprove []string          `json:"autoApprove"`

	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	running bool
	mutex   sync.Mutex
	logger  *logging.Logger
}

// MCPManager manages MCP servers
type MCPManager struct {
	config  *config.Configuration
	servers map[string]*MCPServer
	logger  *logging.Logger
}

// MCPSettings represents the MCP settings configuration file
type MCPSettings struct {
	MCPServers map[string]*MCPServer `json:"mcpServers"`
}

// NewMCPManager creates a new MCP manager
func NewMCPManager(config *config.Configuration) *MCPManager {
	return &MCPManager{
		config:  config,
		servers: make(map[string]*MCPServer),
		logger:  logging.GetLogger(),
	}
}

// LoadSettings loads MCP settings from the configuration file
func (m *MCPManager) LoadSettings() error {
	m.logger.Info("Loading MCP settings")

	// Get the MCP settings file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Check for settings in multiple locations
	settingsPaths := []string{
		filepath.Join(homeDir, ".config", "local_ai_agents", "mcp_settings.json"),
		filepath.Join(homeDir, "Library", "Application Support", "Local AI Agents", "mcp_settings.json"),
		filepath.Join(m.config.ProjectRoot, "config", "mcp_settings.json"),
	}

	var settingsFile string
	for _, path := range settingsPaths {
		if _, err := os.Stat(path); err == nil {
			settingsFile = path
			break
		}
	}

	// If no settings file found, create a default one
	if settingsFile == "" {
		m.logger.Info("No MCP settings file found, creating default")

		// Create directory if it doesn't exist
		settingsDir := filepath.Join(homeDir, ".config", "local_ai_agents")
		if err := os.MkdirAll(settingsDir, 0755); err != nil {
			return fmt.Errorf("failed to create settings directory: %w", err)
		}

		settingsFile = filepath.Join(settingsDir, "mcp_settings.json")

		// Create default settings
		defaultSettings := MCPSettings{
			MCPServers: make(map[string]*MCPServer),
		}

		// Write default settings to file
		data, err := json.MarshalIndent(defaultSettings, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal default settings: %w", err)
		}

		if err := os.WriteFile(settingsFile, data, 0644); err != nil {
			return fmt.Errorf("failed to write default settings file: %w", err)
		}
	}

	// Read settings file
	data, err := os.ReadFile(settingsFile)
	if err != nil {
		return fmt.Errorf("failed to read MCP settings file: %w", err)
	}

	// Parse settings
	var settings MCPSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("failed to parse MCP settings: %w", err)
	}

	// Store servers
	m.servers = settings.MCPServers

	m.logger.Info("Loaded %d MCP servers from settings", len(m.servers))
	return nil
}

// SaveSettings saves MCP settings to the configuration file
func (m *MCPManager) SaveSettings() error {
	m.logger.Info("Saving MCP settings")

	// Get the MCP settings file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	settingsDir := filepath.Join(homeDir, ".config", "local_ai_agents")
	if err := os.MkdirAll(settingsDir, 0755); err != nil {
		return fmt.Errorf("failed to create settings directory: %w", err)
	}

	settingsFile := filepath.Join(settingsDir, "mcp_settings.json")

	// Create settings object
	settings := MCPSettings{
		MCPServers: m.servers,
	}

	// Write settings to file
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	m.logger.Info("Saved MCP settings")
	return nil
}

// StartServer starts an MCP server
func (m *MCPManager) StartServer(name string) error {
	server, exists := m.servers[name]
	if !exists {
		return fmt.Errorf("MCP server '%s' not found", name)
	}

	if server.Disabled {
		return fmt.Errorf("MCP server '%s' is disabled", name)
	}

	server.mutex.Lock()
	defer server.mutex.Unlock()

	if server.running {
		return nil // Server already running
	}

	m.logger.Info("Starting MCP server: %s", name)

	// Create command
	cmd := exec.Command(server.Command, server.Args...)

	// Set environment variables
	env := os.Environ()
	for k, v := range server.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	// Set up pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	// Store command and pipes
	server.cmd = cmd
	server.stdin = stdin
	server.stdout = stdout
	server.stderr = stderr
	server.running = true

	// Start goroutine to handle server output
	go func() {
		defer func() {
			server.mutex.Lock()
			server.running = false
			server.mutex.Unlock()
		}()

		// Wait for command to complete
		err := cmd.Wait()
		if err != nil {
			m.logger.Error("MCP server '%s' exited with error: %v", name, err)
		} else {
			m.logger.Info("MCP server '%s' exited", name)
		}
	}()

	// Start goroutine to log stderr
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				if err != io.EOF {
					m.logger.Error("Error reading from MCP server '%s' stderr: %v", name, err)
				}
				break
			}

			if n > 0 {
				m.logger.Debug("MCP server '%s' stderr: %s", name, strings.TrimSpace(string(buf[:n])))
			}
		}
	}()

	m.logger.Info("MCP server '%s' started", name)
	return nil
}

// StopServer stops an MCP server
func (m *MCPManager) StopServer(name string) error {
	server, exists := m.servers[name]
	if !exists {
		return fmt.Errorf("MCP server '%s' not found", name)
	}

	server.mutex.Lock()
	defer server.mutex.Unlock()

	if !server.running {
		return nil // Server not running
	}

	m.logger.Info("Stopping MCP server: %s", name)

	// Close stdin to signal server to shut down
	if server.stdin != nil {
		server.stdin.Close()
	}

	// Wait for a short time for the server to shut down gracefully
	done := make(chan struct{})
	go func() {
		time.Sleep(2 * time.Second)
		close(done)
	}()

	select {
	case <-done:
		// If server still running after timeout, kill it
		if server.cmd != nil && server.cmd.Process != nil {
			m.logger.Warning("MCP server '%s' did not shut down gracefully, killing", name)
			server.cmd.Process.Kill()
		}
	}

	server.running = false
	m.logger.Info("MCP server '%s' stopped", name)
	return nil
}

// StartAllServers starts all enabled MCP servers
func (m *MCPManager) StartAllServers() {
	m.logger.Info("Starting all enabled MCP servers")

	for name, server := range m.servers {
		if !server.Disabled {
			if err := m.StartServer(name); err != nil {
				m.logger.Error("Failed to start MCP server '%s': %v", name, err)
			}
		}
	}
}

// StopAllServers stops all running MCP servers
func (m *MCPManager) StopAllServers() {
	m.logger.Info("Stopping all MCP servers")

	for name := range m.servers {
		if err := m.StopServer(name); err != nil {
			m.logger.Error("Failed to stop MCP server '%s': %v", name, err)
		}
	}
}

// AddServer adds a new MCP server
func (m *MCPManager) AddServer(name, command string, args []string, env map[string]string) error {
	m.logger.Info("Adding MCP server: %s", name)

	if _, exists := m.servers[name]; exists {
		return fmt.Errorf("MCP server '%s' already exists", name)
	}

	server := &MCPServer{
		Name:        name,
		Command:     command,
		Args:        args,
		Env:         env,
		Disabled:    false,
		AutoApprove: []string{},
		logger:      m.logger,
	}

	m.servers[name] = server

	// Save settings
	if err := m.SaveSettings(); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	return nil
}

// RemoveServer removes an MCP server
func (m *MCPManager) RemoveServer(name string) error {
	m.logger.Info("Removing MCP server: %s", name)

	server, exists := m.servers[name]
	if !exists {
		return fmt.Errorf("MCP server '%s' not found", name)
	}

	// Stop server if running
	if server.running {
		if err := m.StopServer(name); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
	}

	// Remove server
	delete(m.servers, name)

	// Save settings
	if err := m.SaveSettings(); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	return nil
}

// ListServers returns a list of all MCP servers
func (m *MCPManager) ListServers() []*MCPServer {
	servers := make([]*MCPServer, 0, len(m.servers))

	for _, server := range m.servers {
		servers = append(servers, server)
	}

	return servers
}

// GetServer returns an MCP server by name
func (m *MCPManager) GetServer(name string) (*MCPServer, error) {
	server, exists := m.servers[name]
	if !exists {
		return nil, fmt.Errorf("MCP server '%s' not found", name)
	}

	return server, nil
}

// IsServerRunning checks if an MCP server is running
func (m *MCPManager) IsServerRunning(name string) (bool, error) {
	server, exists := m.servers[name]
	if !exists {
		return false, fmt.Errorf("MCP server '%s' not found", name)
	}

	server.mutex.Lock()
	defer server.mutex.Unlock()

	return server.running, nil
}

// EnableServer enables an MCP server
func (m *MCPManager) EnableServer(name string) error {
	server, exists := m.servers[name]
	if !exists {
		return fmt.Errorf("MCP server '%s' not found", name)
	}

	server.Disabled = false

	// Save settings
	if err := m.SaveSettings(); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	return nil
}

// DisableServer disables an MCP server
func (m *MCPManager) DisableServer(name string) error {
	server, exists := m.servers[name]
	if !exists {
		return fmt.Errorf("MCP server '%s' not found", name)
	}

	// Stop server if running
	if server.running {
		if err := m.StopServer(name); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
	}

	server.Disabled = true

	// Save settings
	if err := m.SaveSettings(); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	return nil
}

// UseTool uses an MCP tool
func (m *MCPManager) UseTool(serverName, toolName string, arguments map[string]interface{}) (string, error) {
	server, exists := m.servers[serverName]
	if !exists {
		return "", fmt.Errorf("MCP server '%s' not found", serverName)
	}

	if server.Disabled {
		return "", fmt.Errorf("MCP server '%s' is disabled", serverName)
	}

	if !server.running {
		return "", fmt.Errorf("MCP server '%s' is not running", serverName)
	}

	// Create request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "callTool",
		"params": map[string]interface{}{
			"name":      toolName,
			"arguments": arguments,
		},
		"id": 1,
	}

	// Send request
	response, err := m.sendRequest(server, request)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	// Check for error
	if errObj, ok := response["error"].(map[string]interface{}); ok {
		return "", fmt.Errorf("MCP error: %v", errObj["message"])
	}

	// Get result
	result, ok := response["result"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	// Get content
	content, ok := result["content"].([]interface{})
	if !ok {
		return "", fmt.Errorf("invalid content format")
	}

	// Build response string
	var responseStr strings.Builder
	for _, item := range content {
		if textItem, ok := item.(map[string]interface{}); ok {
			if text, ok := textItem["text"].(string); ok {
				responseStr.WriteString(text)
			}
		}
	}

	return responseStr.String(), nil
}

// AccessResource accesses an MCP resource
func (m *MCPManager) AccessResource(serverName, uri string) (string, error) {
	server, exists := m.servers[serverName]
	if !exists {
		return "", fmt.Errorf("MCP server '%s' not found", serverName)
	}

	if server.Disabled {
		return "", fmt.Errorf("MCP server '%s' is disabled", serverName)
	}

	if !server.running {
		return "", fmt.Errorf("MCP server '%s' is not running", serverName)
	}

	// Create request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "readResource",
		"params": map[string]interface{}{
			"uri": uri,
		},
		"id": 1,
	}

	// Send request
	response, err := m.sendRequest(server, request)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	// Check for error
	if errObj, ok := response["error"].(map[string]interface{}); ok {
		return "", fmt.Errorf("MCP error: %v", errObj["message"])
	}

	// Get result
	result, ok := response["result"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	// Get contents
	contents, ok := result["contents"].([]interface{})
	if !ok || len(contents) == 0 {
		return "", fmt.Errorf("invalid contents format")
	}

	// Get first content item
	content, ok := contents[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid content item format")
	}

	// Get text
	text, ok := content["text"].(string)
	if !ok {
		return "", fmt.Errorf("invalid text format")
	}

	return text, nil
}

// ListTools lists the tools provided by an MCP server
func (m *MCPManager) ListTools(serverName string) ([]map[string]interface{}, error) {
	server, exists := m.servers[serverName]
	if !exists {
		return nil, fmt.Errorf("MCP server '%s' not found", serverName)
	}

	if server.Disabled {
		return nil, fmt.Errorf("MCP server '%s' is disabled", serverName)
	}

	if !server.running {
		return nil, fmt.Errorf("MCP server '%s' is not running", serverName)
	}

	// Create request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "listTools",
		"params":  map[string]interface{}{},
		"id":      1,
	}

	// Send request
	response, err := m.sendRequest(server, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Check for error
	if errObj, ok := response["error"].(map[string]interface{}); ok {
		return nil, fmt.Errorf("MCP error: %v", errObj["message"])
	}

	// Get result
	result, ok := response["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	// Get tools
	toolsRaw, ok := result["tools"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid tools format")
	}

	// Convert tools to map
	tools := make([]map[string]interface{}, 0, len(toolsRaw))
	for _, toolRaw := range toolsRaw {
		if tool, ok := toolRaw.(map[string]interface{}); ok {
			tools = append(tools, tool)
		}
	}

	return tools, nil
}

// ListResources lists the resources provided by an MCP server
func (m *MCPManager) ListResources(serverName string) ([]map[string]interface{}, error) {
	server, exists := m.servers[serverName]
	if !exists {
		return nil, fmt.Errorf("MCP server '%s' not found", serverName)
	}

	if server.Disabled {
		return nil, fmt.Errorf("MCP server '%s' is disabled", serverName)
	}

	if !server.running {
		return nil, fmt.Errorf("MCP server '%s' is not running", serverName)
	}

	// Create request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "listResources",
		"params":  map[string]interface{}{},
		"id":      1,
	}

	// Send request
	response, err := m.sendRequest(server, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Check for error
	if errObj, ok := response["error"].(map[string]interface{}); ok {
		return nil, fmt.Errorf("MCP error: %v", errObj["message"])
	}

	// Get result
	result, ok := response["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	// Get resources
	resourcesRaw, ok := result["resources"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid resources format")
	}

	// Convert resources to map
	resources := make([]map[string]interface{}, 0, len(resourcesRaw))
	for _, resourceRaw := range resourcesRaw {
		if resource, ok := resourceRaw.(map[string]interface{}); ok {
			resources = append(resources, resource)
		}
	}

	return resources, nil
}

// ListResourceTemplates lists the resource templates provided by an MCP server
func (m *MCPManager) ListResourceTemplates(serverName string) ([]map[string]interface{}, error) {
	server, exists := m.servers[serverName]
	if !exists {
		return nil, fmt.Errorf("MCP server '%s' not found", serverName)
	}

	if server.Disabled {
		return nil, fmt.Errorf("MCP server '%s' is disabled", serverName)
	}

	if !server.running {
		return nil, fmt.Errorf("MCP server '%s' is not running", serverName)
	}

	// Create request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "listResourceTemplates",
		"params":  map[string]interface{}{},
		"id":      1,
	}

	// Send request
	response, err := m.sendRequest(server, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Check for error
	if errObj, ok := response["error"].(map[string]interface{}); ok {
		return nil, fmt.Errorf("MCP error: %v", errObj["message"])
	}

	// Get result
	result, ok := response["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	// Get resource templates
	templatesRaw, ok := result["resourceTemplates"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid resource templates format")
	}

	// Convert resource templates to map
	templates := make([]map[string]interface{}, 0, len(templatesRaw))
	for _, templateRaw := range templatesRaw {
		if template, ok := templateRaw.(map[string]interface{}); ok {
			templates = append(templates, template)
		}
	}

	return templates, nil
}

// sendRequest sends a request to an MCP server and returns the response
func (m *MCPManager) sendRequest(server *MCPServer, request map[string]interface{}) (map[string]interface{}, error) {
	// Marshal request
	requestData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send request
	server.mutex.Lock()
	if !server.running {
		server.mutex.Unlock()
		return nil, errors.NewAPIConnectionError("server not running")
	}

	if _, err := server.stdin.Write(requestData); err != nil {
		server.mutex.Unlock()
		return nil, fmt.Errorf("failed to write request: %w", err)
	}

	if _, err := server.stdin.Write([]byte("\n")); err != nil {
		server.mutex.Unlock()
		return nil, fmt.Errorf("failed to write newline: %w", err)
	}

	server.mutex.Unlock()

	// Read response
	buf := make([]byte, 4096)
	n, err := server.stdout.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(buf[:n], &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

// DefaultManager is the default MCP manager instance
var DefaultManager *MCPManager

// Initialize initializes the MCP module
func Initialize(config *config.Configuration) error {
	logger := logging.GetLogger()
	logger.Info("Initializing MCP module")

	// Create MCP manager
	DefaultManager = NewMCPManager(config)

	// Load settings
	if err := DefaultManager.LoadSettings(); err != nil {
		logger.Error("Failed to load MCP settings: %v", err)
		return err
	}

	// Start all enabled servers
	DefaultManager.StartAllServers()

	logger.Info("MCP module initialized")
	return nil
}

// Shutdown shuts down the MCP module
func Shutdown() {
	if DefaultManager != nil {
		logging.GetLogger().Info("Shutting down MCP module")
		DefaultManager.StopAllServers()
	}
}
