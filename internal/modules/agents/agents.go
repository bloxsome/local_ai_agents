// Package agents provides functionality for managing different agent personalities
package agents

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bloxsome/local_ai_agents/internal/config"
	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
)

// AgentPersonality represents an agent personality
type AgentPersonality struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	PromptTemplate   string            `json:"promptTemplate"`
	SystemPrompt     string            `json:"systemPrompt"`
	Traits           map[string]string `json:"traits"`
	PreferredModel   string            `json:"preferredModel"`
	ResponseTemplate string            `json:"responseTemplate"`
}

// AgentManager manages agent personalities
type AgentManager struct {
	config        *config.Configuration
	logger        *logging.Logger
	personalities map[string]*AgentPersonality
	activeAgent   string
	mutex         sync.RWMutex
}

// NewAgentManager creates a new AgentManager
func NewAgentManager(config *config.Configuration) *AgentManager {
	return &AgentManager{
		config:        config,
		logger:        logging.GetLogger(),
		personalities: make(map[string]*AgentPersonality),
		activeAgent:   "",
		mutex:         sync.RWMutex{},
	}
}

// DefaultManager is the default AgentManager instance
var DefaultManager *AgentManager

// Initialize initializes the agents module
func Initialize(config *config.Configuration) error {
	logger := logging.GetLogger()
	logger.Info("Initializing agents module")

	// Create agent manager
	DefaultManager = NewAgentManager(config)

	// Load agent personalities
	if err := DefaultManager.LoadPersonalities(); err != nil {
		logger.Warning("Failed to load agent personalities: %v", err)
	}

	// Set default agent
	if DefaultManager.activeAgent == "" {
		DefaultManager.SetActiveAgent(config.AgentName)
	}

	logger.Info("Agents module initialized")
	return nil
}

// Shutdown shuts down the agents module
func Shutdown() {
	if DefaultManager != nil {
		logging.GetLogger().Info("Shutting down agents module")
		if err := DefaultManager.SavePersonalities(); err != nil {
			logging.GetLogger().Warning("Failed to save agent personalities: %v", err)
		}
	}
}

// LoadPersonalities loads agent personalities from a file
func (m *AgentManager) LoadPersonalities() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Get the personalities file path
	personalitiesDir := filepath.Join(m.config.ProjectRoot, "data", "agents")
	if err := os.MkdirAll(personalitiesDir, 0755); err != nil {
		return fmt.Errorf("failed to create personalities directory: %w", err)
	}

	personalitiesFile := filepath.Join(personalitiesDir, "personalities.json")

	// Check if the personalities file exists
	if _, err := os.Stat(personalitiesFile); os.IsNotExist(err) {
		m.logger.Info("Personalities file does not exist, creating default personalities")

		// Create default personalities
		m.personalities = map[string]*AgentPersonality{
			"Otto": {
				Name:           "Otto",
				Description:    "A helpful and friendly AI assistant",
				PromptTemplate: "{{context}}User: {{input}}\n\nOtto: ",
				SystemPrompt:   "You are Otto, a helpful and friendly AI assistant. You provide clear and concise answers to user questions.",
				Traits: map[string]string{
					"Helpfulness":  "High",
					"Friendliness": "High",
					"Formality":    "Medium",
					"Creativity":   "Medium",
				},
				PreferredModel:   "llama3.3:latest",
				ResponseTemplate: "{{response}}",
			},
			"Cline": {
				Name:           "Cline",
				Description:    "A technical and precise coding assistant",
				PromptTemplate: "{{context}}User: {{input}}\n\nCline: ",
				SystemPrompt:   "You are Cline, a technical and precise coding assistant. You provide detailed and accurate code examples and explanations.",
				Traits: map[string]string{
					"Technical":  "High",
					"Precision":  "High",
					"Formality":  "High",
					"Creativity": "Low",
				},
				PreferredModel:   "llama3.3:latest",
				ResponseTemplate: "{{response}}",
			},
			"Luna": {
				Name:           "Luna",
				Description:    "A creative and imaginative storyteller",
				PromptTemplate: "{{context}}User: {{input}}\n\nLuna: ",
				SystemPrompt:   "You are Luna, a creative and imaginative storyteller. You help users craft engaging stories and creative content.",
				Traits: map[string]string{
					"Creativity":  "Very High",
					"Imagination": "Very High",
					"Formality":   "Low",
					"Technical":   "Low",
				},
				PreferredModel:   "llama3.3:latest",
				ResponseTemplate: "{{response}}",
			},
		}

		// Set default active agent
		m.activeAgent = "Otto"

		// Save default personalities
		return m.SavePersonalities()
	}

	// Read the personalities file
	data, err := os.ReadFile(personalitiesFile)
	if err != nil {
		return fmt.Errorf("failed to read personalities file: %w", err)
	}

	// Parse the personalities
	var personalities map[string]*AgentPersonality
	if err := json.Unmarshal(data, &personalities); err != nil {
		return fmt.Errorf("failed to parse personalities: %w", err)
	}

	m.personalities = personalities
	m.logger.Info("Loaded %d agent personalities", len(m.personalities))

	// Read active agent from a separate file
	activeAgentFile := filepath.Join(personalitiesDir, "active_agent.txt")
	if activeAgentData, err := os.ReadFile(activeAgentFile); err == nil {
		m.activeAgent = string(activeAgentData)
		m.logger.Info("Active agent: %s", m.activeAgent)
	} else {
		// Set default active agent if not found
		if len(m.personalities) > 0 {
			for name := range m.personalities {
				m.activeAgent = name
				break
			}
		}
	}

	return nil
}

// SavePersonalities saves agent personalities to a file
func (m *AgentManager) SavePersonalities() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Get the personalities file path
	personalitiesDir := filepath.Join(m.config.ProjectRoot, "data", "agents")
	if err := os.MkdirAll(personalitiesDir, 0755); err != nil {
		return fmt.Errorf("failed to create personalities directory: %w", err)
	}

	personalitiesFile := filepath.Join(personalitiesDir, "personalities.json")

	// Marshal the personalities
	data, err := json.MarshalIndent(m.personalities, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal personalities: %w", err)
	}

	// Write the personalities file
	if err := os.WriteFile(personalitiesFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write personalities file: %w", err)
	}

	// Save active agent to a separate file
	activeAgentFile := filepath.Join(personalitiesDir, "active_agent.txt")
	if err := os.WriteFile(activeAgentFile, []byte(m.activeAgent), 0644); err != nil {
		return fmt.Errorf("failed to write active agent file: %w", err)
	}

	m.logger.Info("Saved %d agent personalities", len(m.personalities))
	return nil
}

// GetPersonalities returns all agent personalities
func (m *AgentManager) GetPersonalities() map[string]*AgentPersonality {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy of the personalities to prevent modification
	personalitiesCopy := make(map[string]*AgentPersonality)
	for name, personality := range m.personalities {
		personalityCopy := *personality
		personalitiesCopy[name] = &personalityCopy
	}

	return personalitiesCopy
}

// GetPersonality returns an agent personality by name
func (m *AgentManager) GetPersonality(name string) (*AgentPersonality, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	personality, exists := m.personalities[name]
	if !exists {
		return nil, fmt.Errorf("agent personality '%s' not found", name)
	}

	// Return a copy of the personality to prevent modification
	personalityCopy := *personality
	return &personalityCopy, nil
}

// AddPersonality adds a new agent personality
func (m *AgentManager) AddPersonality(personality *AgentPersonality) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.personalities[personality.Name]; exists {
		return fmt.Errorf("agent personality '%s' already exists", personality.Name)
	}

	// Add the personality
	m.personalities[personality.Name] = personality

	// Save personalities
	return m.SavePersonalities()
}

// UpdatePersonality updates an existing agent personality
func (m *AgentManager) UpdatePersonality(personality *AgentPersonality) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.personalities[personality.Name]; !exists {
		return fmt.Errorf("agent personality '%s' not found", personality.Name)
	}

	// Update the personality
	m.personalities[personality.Name] = personality

	// Save personalities
	return m.SavePersonalities()
}

// RemovePersonality removes an agent personality
func (m *AgentManager) RemovePersonality(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.personalities[name]; !exists {
		return fmt.Errorf("agent personality '%s' not found", name)
	}

	// Check if this is the active agent
	if m.activeAgent == name {
		return fmt.Errorf("cannot remove active agent personality '%s'", name)
	}

	// Remove the personality
	delete(m.personalities, name)

	// Save personalities
	return m.SavePersonalities()
}

// GetActiveAgent returns the active agent personality
func (m *AgentManager) GetActiveAgent() (*AgentPersonality, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.activeAgent == "" {
		return nil, fmt.Errorf("no active agent personality set")
	}

	personality, exists := m.personalities[m.activeAgent]
	if !exists {
		return nil, fmt.Errorf("active agent personality '%s' not found", m.activeAgent)
	}

	// Return a copy of the personality to prevent modification
	personalityCopy := *personality
	return &personalityCopy, nil
}

// SetActiveAgent sets the active agent personality
func (m *AgentManager) SetActiveAgent(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.personalities[name]; !exists {
		return fmt.Errorf("agent personality '%s' not found", name)
	}

	m.activeAgent = name
	m.logger.Info("Active agent set to: %s", name)

	// Save active agent
	personalitiesDir := filepath.Join(m.config.ProjectRoot, "data", "agents")
	activeAgentFile := filepath.Join(personalitiesDir, "active_agent.txt")
	if err := os.WriteFile(activeAgentFile, []byte(m.activeAgent), 0644); err != nil {
		return fmt.Errorf("failed to write active agent file: %w", err)
	}

	return nil
}

// FormatPrompt formats a prompt using the active agent's template
func (m *AgentManager) FormatPrompt(context, input string) (string, error) {
	activeAgent, err := m.GetActiveAgent()
	if err != nil {
		return "", err
	}

	// Replace placeholders in the prompt template
	prompt := activeAgent.PromptTemplate
	prompt = strings.Replace(prompt, "{{context}}", context, -1)
	prompt = strings.Replace(prompt, "{{input}}", input, -1)

	return prompt, nil
}

// FormatResponse formats a response using the active agent's template
func (m *AgentManager) FormatResponse(response string) (string, error) {
	activeAgent, err := m.GetActiveAgent()
	if err != nil {
		return "", err
	}

	// Replace placeholders in the response template
	formattedResponse := activeAgent.ResponseTemplate
	formattedResponse = strings.Replace(formattedResponse, "{{response}}", response, -1)

	return formattedResponse, nil
}

// GetPreferredModel returns the preferred model for the active agent
func (m *AgentManager) GetPreferredModel() (string, error) {
	activeAgent, err := m.GetActiveAgent()
	if err != nil {
		return "", err
	}

	return activeAgent.PreferredModel, nil
}

// GetSystemPrompt returns the system prompt for the active agent
func (m *AgentManager) GetSystemPrompt() (string, error) {
	activeAgent, err := m.GetActiveAgent()
	if err != nil {
		return "", err
	}

	return activeAgent.SystemPrompt, nil
}

// GetPersonalities returns all agent personalities using the default manager
func GetPersonalities() map[string]*AgentPersonality {
	if DefaultManager != nil {
		return DefaultManager.GetPersonalities()
	}
	return map[string]*AgentPersonality{}
}

// GetPersonality returns an agent personality by name using the default manager
func GetPersonality(name string) (*AgentPersonality, error) {
	if DefaultManager != nil {
		return DefaultManager.GetPersonality(name)
	}
	return nil, fmt.Errorf("agents module not initialized")
}

// AddPersonality adds a new agent personality using the default manager
func AddPersonality(personality *AgentPersonality) error {
	if DefaultManager != nil {
		return DefaultManager.AddPersonality(personality)
	}
	return fmt.Errorf("agents module not initialized")
}

// UpdatePersonality updates an existing agent personality using the default manager
func UpdatePersonality(personality *AgentPersonality) error {
	if DefaultManager != nil {
		return DefaultManager.UpdatePersonality(personality)
	}
	return fmt.Errorf("agents module not initialized")
}

// RemovePersonality removes an agent personality using the default manager
func RemovePersonality(name string) error {
	if DefaultManager != nil {
		return DefaultManager.RemovePersonality(name)
	}
	return fmt.Errorf("agents module not initialized")
}

// GetActiveAgent returns the active agent personality using the default manager
func GetActiveAgent() (*AgentPersonality, error) {
	if DefaultManager != nil {
		return DefaultManager.GetActiveAgent()
	}
	return nil, fmt.Errorf("agents module not initialized")
}

// SetActiveAgent sets the active agent personality using the default manager
func SetActiveAgent(name string) error {
	if DefaultManager != nil {
		return DefaultManager.SetActiveAgent(name)
	}
	return fmt.Errorf("agents module not initialized")
}

// FormatPrompt formats a prompt using the active agent's template using the default manager
func FormatPrompt(context, input string) (string, error) {
	if DefaultManager != nil {
		return DefaultManager.FormatPrompt(context, input)
	}
	return "", fmt.Errorf("agents module not initialized")
}

// FormatResponse formats a response using the active agent's template using the default manager
func FormatResponse(response string) (string, error) {
	if DefaultManager != nil {
		return DefaultManager.FormatResponse(response)
	}
	return "", fmt.Errorf("agents module not initialized")
}

// GetPreferredModel returns the preferred model for the active agent using the default manager
func GetPreferredModel() (string, error) {
	if DefaultManager != nil {
		return DefaultManager.GetPreferredModel()
	}
	return "", fmt.Errorf("agents module not initialized")
}

// GetSystemPrompt returns the system prompt for the active agent using the default manager
func GetSystemPrompt() (string, error) {
	if DefaultManager != nil {
		return DefaultManager.GetSystemPrompt()
	}
	return "", fmt.Errorf("agents module not initialized")
}
