// Package main provides configuration settings for the Local AI Agents application
package main

import (
	"os"
	"path/filepath"
	"strconv"
)

// Configuration holds all the configuration settings for the application
type Configuration struct {
	// User and Agent configuration
	UserName  string
	AgentName string

	// Model configuration
	DefaultModel   string
	EmbeddingModel string

	// Memory configuration
	MemoryLength int
	ChunkSize    int
	ChunkOverlap int
	ChunkLength  int

	// Path configuration
	ProjectRoot     string
	DataDir         string
	EmbeddingsDir   string
	ChatHistoryFile string

	// Edge Database configuration
	DBDir       string
	DBFile      string
	DBPath      string
	DBVersion   string
	TableName   string
	IndexPrefix string

	// Search configuration
	DefaultTopK                int
	DefaultSimilarityThreshold float64

	// Logging configuration
	LogLevel string
	LogFile  string

	// Assistant command configurations
	TerminalApp    []string
	DefaultBrowser string

	// MCP configuration
	MCPEnabled bool
	MCPDir     string
}

// LoadConfig loads the configuration from environment variables and defaults
func LoadConfig() (*Configuration, error) {
	// Get the project root directory
	projectRoot, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	// Create data directories
	dataDir := filepath.Join(projectRoot, "data", "json_history")
	embeddingsDir := filepath.Join(dataDir, "embeddings")
	dbDir := filepath.Join(projectRoot, "data", "edgebase")

	// Ensure directories exist
	dirs := []string{dataDir, embeddingsDir, dbDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	// Get home directory for chat history file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Create log directory
	logDir := filepath.Join(projectRoot, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// Create MCP directory
	mcpDir := filepath.Join(homeDir, "Documents", "Cline", "MCP")
	if err := os.MkdirAll(mcpDir, 0755); err != nil {
		return nil, err
	}

	// Load configuration from environment variables with defaults
	config := &Configuration{
		// User and Agent configuration
		UserName:  getEnv("AI_USER_NAME", "MikeBee"),
		AgentName: getEnv("AI_AGENT_NAME", "Otto"),

		// Model configuration
		DefaultModel:   "llama3.3:latest",
		EmbeddingModel: getEnv("AI_EMBEDDING_MODEL", "nomic-embed-text"),

		// Memory configuration
		MemoryLength: getEnvInt("AI_MEMORY_LENGTH", 15),
		ChunkSize:    getEnvInt("AI_CHUNK_SIZE", 5000),
		ChunkOverlap: getEnvInt("AI_CHUNK_OVERLAP", 200),
		ChunkLength:  getEnvInt("AI_CHUNK_LENGTH", 10),

		// Path configuration
		ProjectRoot:     projectRoot,
		DataDir:         dataDir,
		EmbeddingsDir:   embeddingsDir,
		ChatHistoryFile: filepath.Join(homeDir, ".local_ai_agents_chat_history.json"),

		// Edge Database configuration
		DBDir:       dbDir,
		DBFile:      "knowledge_edges.db",
		DBPath:      filepath.Join(dbDir, "knowledge_edges.db"),
		DBVersion:   "1.0",
		TableName:   "edges",
		IndexPrefix: "idx_",

		// Search configuration
		DefaultTopK:                getEnvInt("AI_DEFAULT_TOP_K", 5),
		DefaultSimilarityThreshold: getEnvFloat("AI_DEFAULT_SIMILARITY_THRESHOLD", 0.0),

		// Logging configuration
		LogLevel: getEnv("AI_LOG_LEVEL", "WARNING"),
		LogFile:  filepath.Join(logDir, "local_ai_agents.log"),

		// Assistant command configurations
		TerminalApp:    []string{"open", "-a", "Terminal"},
		DefaultBrowser: "default",

		// MCP configuration
		MCPEnabled: getEnvBool("AI_MCP_ENABLED", true),
		MCPDir:     getEnv("AI_MCP_DIR", mcpDir),
	}

	return config, nil
}

// Helper functions to get environment variables with defaults

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvFloat(key string, defaultValue float64) float64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
