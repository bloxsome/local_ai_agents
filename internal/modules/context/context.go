// Package context provides functionality for managing conversation context and history
package context

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bloxsome/local_ai_agents/internal/config"
	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
)

// ContextManager manages conversation context and history
type ContextManager struct {
	config         *config.Configuration
	logger         *logging.Logger
	history        []Message
	context        []string
	bullets        []string
	knowledgeTree  map[string][]string
	userProfile    map[string]string
	mutex          sync.RWMutex
	historyChanged bool
}

// Message represents a message in the conversation history
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// NewContextManager creates a new ContextManager
func NewContextManager(config *config.Configuration) *ContextManager {
	return &ContextManager{
		config:        config,
		logger:        logging.GetLogger(),
		history:       []Message{},
		context:       []string{},
		bullets:       []string{},
		knowledgeTree: make(map[string][]string),
		userProfile:   make(map[string]string),
		mutex:         sync.RWMutex{},
	}
}

// DefaultManager is the default ContextManager instance
var DefaultManager *ContextManager

// Initialize initializes the context module
func Initialize(config *config.Configuration) error {
	logger := logging.GetLogger()
	logger.Info("Initializing context module")

	// Create context manager
	DefaultManager = NewContextManager(config)

	// Load history from file
	if err := DefaultManager.LoadHistory(); err != nil {
		logger.Warning("Failed to load history: %v", err)
	}

	logger.Info("Context module initialized")
	return nil
}

// Shutdown shuts down the context module
func Shutdown() {
	if DefaultManager != nil {
		logging.GetLogger().Info("Shutting down context module")
		if err := DefaultManager.SaveHistory(); err != nil {
			logging.GetLogger().Warning("Failed to save history: %v", err)
		}
	}
}

// AddMessage adds a message to the conversation history
func (m *ContextManager) AddMessage(role, content string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	message := Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}

	m.history = append(m.history, message)
	m.historyChanged = true

	// Trim history if it exceeds the configured length
	if len(m.history) > m.config.MemoryLength {
		m.history = m.history[len(m.history)-m.config.MemoryLength:]
	}

	// Extract context from user messages
	if role == "user" {
		m.extractContext(content)
	}
}

// GetHistory returns the conversation history
func (m *ContextManager) GetHistory() []Message {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy of the history to prevent modification
	historyCopy := make([]Message, len(m.history))
	copy(historyCopy, m.history)

	return historyCopy
}

// GetContext returns the current context
func (m *ContextManager) GetContext() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy of the context to prevent modification
	contextCopy := make([]string, len(m.context))
	copy(contextCopy, m.context)

	return contextCopy
}

// ClearContext clears the current context
func (m *ContextManager) ClearContext() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.context = []string{}
}

// GetBullets returns the current bullet points
func (m *ContextManager) GetBullets() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy of the bullets to prevent modification
	bulletsCopy := make([]string, len(m.bullets))
	copy(bulletsCopy, m.bullets)

	return bulletsCopy
}

// AddBullet adds a bullet point
func (m *ContextManager) AddBullet(bullet string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.bullets = append(m.bullets, bullet)
}

// ClearBullets clears the current bullet points
func (m *ContextManager) ClearBullets() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.bullets = []string{}
}

// GetKnowledgeTree returns the knowledge tree
func (m *ContextManager) GetKnowledgeTree() map[string][]string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy of the knowledge tree to prevent modification
	treeCopy := make(map[string][]string)
	for k, v := range m.knowledgeTree {
		vCopy := make([]string, len(v))
		copy(vCopy, v)
		treeCopy[k] = vCopy
	}

	return treeCopy
}

// AddKnowledge adds knowledge to the knowledge tree
func (m *ContextManager) AddKnowledge(concept, knowledge string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.knowledgeTree[concept]; !ok {
		m.knowledgeTree[concept] = []string{}
	}

	m.knowledgeTree[concept] = append(m.knowledgeTree[concept], knowledge)
}

// GetUserProfile returns the user profile
func (m *ContextManager) GetUserProfile() map[string]string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy of the user profile to prevent modification
	profileCopy := make(map[string]string)
	for k, v := range m.userProfile {
		profileCopy[k] = v
	}

	return profileCopy
}

// UpdateUserProfile updates the user profile
func (m *ContextManager) UpdateUserProfile(key, value string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.userProfile[key] = value
}

// LoadHistory loads the conversation history from a file
func (m *ContextManager) LoadHistory() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if the history file exists
	if _, err := os.Stat(m.config.ChatHistoryFile); os.IsNotExist(err) {
		m.logger.Info("Chat history file does not exist, starting with empty history")
		return nil
	}

	// Read the history file
	data, err := os.ReadFile(m.config.ChatHistoryFile)
	if err != nil {
		return fmt.Errorf("failed to read chat history file: %w", err)
	}

	// Parse the history
	var history []Message
	if err := json.Unmarshal(data, &history); err != nil {
		return fmt.Errorf("failed to parse chat history: %w", err)
	}

	m.history = history
	m.logger.Info("Loaded %d messages from chat history", len(m.history))

	return nil
}

// SaveHistory saves the conversation history to a file
func (m *ContextManager) SaveHistory() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Only save if the history has changed
	if !m.historyChanged {
		return nil
	}

	// Marshal the history
	data, err := json.MarshalIndent(m.history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal chat history: %w", err)
	}

	// Write the history file
	if err := os.WriteFile(m.config.ChatHistoryFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write chat history file: %w", err)
	}

	m.historyChanged = false
	m.logger.Info("Saved %d messages to chat history", len(m.history))

	return nil
}

// extractContext extracts context from a message
func (m *ContextManager) extractContext(content string) {
	// This is a simple implementation that extracts sentences as context
	// In a real implementation, this would use NLP techniques to extract key information
	sentences := strings.Split(content, ".")
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" && len(sentence) > 10 {
			m.context = append(m.context, sentence)
		}
	}

	// Trim context if it gets too long
	if len(m.context) > m.config.MemoryLength*2 {
		m.context = m.context[len(m.context)-(m.config.MemoryLength*2):]
	}
}

// ExplainConcept provides an explanation of a concept
func (m *ContextManager) ExplainConcept(concept string) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Check if we have knowledge about this concept
	if knowledge, ok := m.knowledgeTree[concept]; ok && len(knowledge) > 0 {
		return fmt.Sprintf("Explanation of %s:\n%s", concept, strings.Join(knowledge, "\n"))
	}

	return fmt.Sprintf("I don't have specific knowledge about %s in my context. Please provide more information.", concept)
}

// FactCheck performs a fact check on a statement
func (m *ContextManager) FactCheck(statement string) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// This is a placeholder implementation
	// In a real implementation, this would use external sources or the knowledge tree to fact check
	return fmt.Sprintf("Fact check for: %s\n\nI don't have enough information in my context to verify this statement. Please provide more information or sources.", statement)
}

// FormatContextForPrompt formats the context for inclusion in a prompt
func (m *ContextManager) FormatContextForPrompt() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var sb strings.Builder

	// Add context
	if len(m.context) > 0 {
		sb.WriteString("Context:\n")
		for _, ctx := range m.context {
			sb.WriteString("- ")
			sb.WriteString(ctx)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Add bullets
	if len(m.bullets) > 0 {
		sb.WriteString("Key Points:\n")
		for _, bullet := range m.bullets {
			sb.WriteString("• ")
			sb.WriteString(bullet)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Add user profile information
	if len(m.userProfile) > 0 {
		sb.WriteString("User Profile:\n")
		for k, v := range m.userProfile {
			sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// AddMessage adds a message to the conversation history using the default manager
func AddMessage(role, content string) {
	if DefaultManager != nil {
		DefaultManager.AddMessage(role, content)
	}
}

// GetHistory returns the conversation history using the default manager
func GetHistory() []Message {
	if DefaultManager != nil {
		return DefaultManager.GetHistory()
	}
	return []Message{}
}

// GetContext returns the current context using the default manager
func GetContext() []string {
	if DefaultManager != nil {
		return DefaultManager.GetContext()
	}
	return []string{}
}

// ClearContext clears the current context using the default manager
func ClearContext() {
	if DefaultManager != nil {
		DefaultManager.ClearContext()
	}
}

// GetBullets returns the current bullet points using the default manager
func GetBullets() []string {
	if DefaultManager != nil {
		return DefaultManager.GetBullets()
	}
	return []string{}
}

// AddBullet adds a bullet point using the default manager
func AddBullet(bullet string) {
	if DefaultManager != nil {
		DefaultManager.AddBullet(bullet)
	}
}

// ClearBullets clears the current bullet points using the default manager
func ClearBullets() {
	if DefaultManager != nil {
		DefaultManager.ClearBullets()
	}
}

// GetKnowledgeTree returns the knowledge tree using the default manager
func GetKnowledgeTree() map[string][]string {
	if DefaultManager != nil {
		return DefaultManager.GetKnowledgeTree()
	}
	return map[string][]string{}
}

// AddKnowledge adds knowledge to the knowledge tree using the default manager
func AddKnowledge(concept, knowledge string) {
	if DefaultManager != nil {
		DefaultManager.AddKnowledge(concept, knowledge)
	}
}

// GetUserProfile returns the user profile using the default manager
func GetUserProfile() map[string]string {
	if DefaultManager != nil {
		return DefaultManager.GetUserProfile()
	}
	return map[string]string{}
}

// UpdateUserProfile updates the user profile using the default manager
func UpdateUserProfile(key, value string) {
	if DefaultManager != nil {
		DefaultManager.UpdateUserProfile(key, value)
	}
}

// ExplainConcept provides an explanation of a concept using the default manager
func ExplainConcept(concept string) string {
	if DefaultManager != nil {
		return DefaultManager.ExplainConcept(concept)
	}
	return fmt.Sprintf("Context module not initialized. Cannot explain %s.", concept)
}

// FactCheck performs a fact check on a statement using the default manager
func FactCheck(statement string) string {
	if DefaultManager != nil {
		return DefaultManager.FactCheck(statement)
	}
	return fmt.Sprintf("Context module not initialized. Cannot fact check: %s", statement)
}

// FormatContextForPrompt formats the context for inclusion in a prompt using the default manager
func FormatContextForPrompt() string {
	if DefaultManager != nil {
		return DefaultManager.FormatContextForPrompt()
	}
	return ""
}

// SaveHistory saves the conversation history to a file using the default manager
func SaveHistory() error {
	if DefaultManager != nil {
		return DefaultManager.SaveHistory()
	}
	return fmt.Errorf("context module not initialized")
}
