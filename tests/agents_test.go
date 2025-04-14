package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/bloxsome/local_ai_agents/internal/config"
	"github.com/bloxsome/local_ai_agents/internal/modules/agents"
	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
)

// TestAgentsSuite represents a test suite for agents module
type TestAgentsSuite struct {
	config *config.Configuration
	logger *logging.Logger
}

// SetupAgentsSuite sets up the test suite for agents tests
func SetupAgentsSuite(t *testing.T) *TestAgentsSuite {
	// Set up logging
	logFile := "agents_test.log"
	if err := logging.SetupLogging(logFile, "DEBUG"); err != nil {
		t.Fatalf("Error setting up logging: %v", err)
	}
	logger := logging.GetLogger()
	logger.Info("Agents test suite started")

	// Create a test configuration
	tempDir, err := os.MkdirTemp("", "agents_test")
	if err != nil {
		t.Fatalf("Error creating temp directory: %v", err)
	}

	cfg := &config.Configuration{
		UserName:        "TestUser",
		AgentName:       "TestAgent",
		DefaultModel:    "llama3.3:latest",
		ProjectRoot:     tempDir,
		ChatHistoryFile: tempDir + "/test_chat_history.json",
	}

	return &TestAgentsSuite{
		config: cfg,
		logger: logger,
	}
}

// TearDownAgentsSuite tears down the test suite
func (s *TestAgentsSuite) TearDownAgentsSuite() {
	// Clean up log file
	logFile := "agents_test.log"
	if err := os.Remove(logFile); err != nil {
		s.logger.Warning("Failed to remove log file: %v", err)
	}

	// Clean up temp directory
	if err := os.RemoveAll(s.config.ProjectRoot); err != nil {
		s.logger.Warning("Failed to remove temp directory: %v", err)
	}

	// Shutdown agents module if it was initialized
	if agents.DefaultManager != nil {
		agents.Shutdown()
	}
}

// TestAgentsInitialization tests initializing the agents module
func TestAgentsInitialization(t *testing.T) {
	suite := SetupAgentsSuite(t)
	defer suite.TearDownAgentsSuite()

	// Initialize agents module
	err := agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Check that the default manager was created
	if agents.DefaultManager == nil {
		t.Fatalf("Default manager is nil after initialization")
	}
}

// TestDefaultPersonalities tests that default personalities are created
func TestDefaultPersonalities(t *testing.T) {
	suite := SetupAgentsSuite(t)
	defer suite.TearDownAgentsSuite()

	// Initialize agents module
	err := agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Get personalities
	personalities := agents.GetPersonalities()

	// Check that we have at least 3 default personalities
	if len(personalities) < 3 {
		t.Errorf("Expected at least 3 default personalities, got %d", len(personalities))
	}

	// Check for specific default personalities
	if _, ok := personalities["Otto"]; !ok {
		t.Errorf("Expected 'Otto' personality, but it was not found")
	}
	if _, ok := personalities["Cline"]; !ok {
		t.Errorf("Expected 'Cline' personality, but it was not found")
	}
	if _, ok := personalities["Luna"]; !ok {
		t.Errorf("Expected 'Luna' personality, but it was not found")
	}
}

// TestGetActiveAgent tests getting the active agent
func TestGetActiveAgent(t *testing.T) {
	suite := SetupAgentsSuite(t)
	defer suite.TearDownAgentsSuite()

	// Initialize agents module
	err := agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Get active agent
	activeAgent, err := agents.GetActiveAgent()
	if err != nil {
		t.Fatalf("Error getting active agent: %v", err)
	}

	// Check that we got a non-nil agent
	if activeAgent == nil {
		t.Fatalf("Active agent is nil")
	}

	// Check that the agent has a name
	if activeAgent.Name == "" {
		t.Errorf("Active agent name is empty")
	}
}

// TestSetActiveAgent tests setting the active agent
func TestSetActiveAgent(t *testing.T) {
	suite := SetupAgentsSuite(t)
	defer suite.TearDownAgentsSuite()

	// Initialize agents module
	err := agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Get personalities
	personalities := agents.GetPersonalities()
	if len(personalities) == 0 {
		t.Fatalf("No personalities found")
	}

	// Get a personality name to switch to
	var targetName string
	for name := range personalities {
		if name != "Otto" { // Try to switch to a non-default agent
			targetName = name
			break
		}
	}
	if targetName == "" {
		targetName = "Otto" // Fall back to Otto if no other agent found
	}

	// Set active agent
	err = agents.SetActiveAgent(targetName)
	if err != nil {
		t.Fatalf("Error setting active agent: %v", err)
	}

	// Get active agent
	activeAgent, err := agents.GetActiveAgent()
	if err != nil {
		t.Fatalf("Error getting active agent after setting: %v", err)
	}

	// Check that the active agent is the one we set
	if activeAgent.Name != targetName {
		t.Errorf("Expected active agent name to be '%s', got '%s'", targetName, activeAgent.Name)
	}
}

// TestAddPersonality tests adding a new personality
func TestAddPersonality(t *testing.T) {
	suite := SetupAgentsSuite(t)
	defer suite.TearDownAgentsSuite()

	// Initialize agents module
	err := agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Create a new personality
	newPersonality := &agents.AgentPersonality{
		Name:             "TestAgent",
		Description:      "A test agent personality",
		PromptTemplate:   "{{context}}User: {{input}}\n\nTestAgent: ",
		SystemPrompt:     "You are TestAgent, a test agent personality.",
		Traits:           map[string]string{"Test": "High"},
		PreferredModel:   "llama3.3:latest",
		ResponseTemplate: "{{response}}",
	}

	// Add the personality
	err = agents.AddPersonality(newPersonality)
	if err != nil {
		t.Fatalf("Error adding personality: %v", err)
	}

	// Get personalities
	personalities := agents.GetPersonalities()

	// Check that the new personality was added
	if _, ok := personalities["TestAgent"]; !ok {
		t.Errorf("Expected 'TestAgent' personality to be added, but it was not found")
	}
}

// TestRemovePersonality tests removing a personality
func TestRemovePersonality(t *testing.T) {
	suite := SetupAgentsSuite(t)
	defer suite.TearDownAgentsSuite()

	// Initialize agents module
	err := agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Create a new personality
	newPersonality := &agents.AgentPersonality{
		Name:             "TestAgent",
		Description:      "A test agent personality",
		PromptTemplate:   "{{context}}User: {{input}}\n\nTestAgent: ",
		SystemPrompt:     "You are TestAgent, a test agent personality.",
		Traits:           map[string]string{"Test": "High"},
		PreferredModel:   "llama3.3:latest",
		ResponseTemplate: "{{response}}",
	}

	// Add the personality
	err = agents.AddPersonality(newPersonality)
	if err != nil {
		t.Fatalf("Error adding personality: %v", err)
	}

	// Make sure we're not trying to remove the active agent
	activeAgent, err := agents.GetActiveAgent()
	if err != nil {
		t.Fatalf("Error getting active agent: %v", err)
	}
	if activeAgent.Name == "TestAgent" {
		// Switch to a different agent
		err = agents.SetActiveAgent("Otto")
		if err != nil {
			t.Fatalf("Error switching active agent: %v", err)
		}
	}

	// Remove the personality
	err = agents.RemovePersonality("TestAgent")
	if err != nil {
		t.Fatalf("Error removing personality: %v", err)
	}

	// Get personalities
	personalities := agents.GetPersonalities()

	// Check that the personality was removed
	if _, ok := personalities["TestAgent"]; ok {
		t.Errorf("Expected 'TestAgent' personality to be removed, but it was found")
	}
}

// TestFormatPrompt tests formatting a prompt using the active agent's template
func TestFormatPrompt(t *testing.T) {
	suite := SetupAgentsSuite(t)
	defer suite.TearDownAgentsSuite()

	// Initialize agents module
	err := agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Set active agent to Otto
	err = agents.SetActiveAgent("Otto")
	if err != nil {
		t.Fatalf("Error setting active agent: %v", err)
	}

	// Format a prompt
	context := "Context information\n"
	input := "Hello, how are you?"
	prompt, err := agents.FormatPrompt(context, input)
	if err != nil {
		t.Fatalf("Error formatting prompt: %v", err)
	}

	// Check that the prompt contains the context and input
	if !strings.Contains(prompt, context) {
		t.Errorf("Expected prompt to contain context '%s', got '%s'", context, prompt)
	}
	if !strings.Contains(prompt, input) {
		t.Errorf("Expected prompt to contain input '%s', got '%s'", input, prompt)
	}
	if !strings.Contains(prompt, "Otto:") {
		t.Errorf("Expected prompt to contain 'Otto:', got '%s'", prompt)
	}
}

// TestFormatResponse tests formatting a response using the active agent's template
func TestFormatResponse(t *testing.T) {
	suite := SetupAgentsSuite(t)
	defer suite.TearDownAgentsSuite()

	// Initialize agents module
	err := agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Set active agent to Otto
	err = agents.SetActiveAgent("Otto")
	if err != nil {
		t.Fatalf("Error setting active agent: %v", err)
	}

	// Format a response
	response := "Hello! I'm doing well. How can I help you today?"
	formattedResponse, err := agents.FormatResponse(response)
	if err != nil {
		t.Fatalf("Error formatting response: %v", err)
	}

	// Check that the formatted response contains the response
	if !strings.Contains(formattedResponse, response) {
		t.Errorf("Expected formatted response to contain '%s', got '%s'", response, formattedResponse)
	}
}

// TestGetPreferredModel tests getting the preferred model for the active agent
func TestGetPreferredModel(t *testing.T) {
	suite := SetupAgentsSuite(t)
	defer suite.TearDownAgentsSuite()

	// Initialize agents module
	err := agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Set active agent to Otto
	err = agents.SetActiveAgent("Otto")
	if err != nil {
		t.Fatalf("Error setting active agent: %v", err)
	}

	// Get preferred model
	model, err := agents.GetPreferredModel()
	if err != nil {
		t.Fatalf("Error getting preferred model: %v", err)
	}

	// Check that we got a non-empty model
	if model == "" {
		t.Errorf("Expected non-empty preferred model, got empty string")
	}
}

// TestGetSystemPrompt tests getting the system prompt for the active agent
func TestGetSystemPrompt(t *testing.T) {
	suite := SetupAgentsSuite(t)
	defer suite.TearDownAgentsSuite()

	// Initialize agents module
	err := agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Set active agent to Otto
	err = agents.SetActiveAgent("Otto")
	if err != nil {
		t.Fatalf("Error setting active agent: %v", err)
	}

	// Get system prompt
	systemPrompt, err := agents.GetSystemPrompt()
	if err != nil {
		t.Fatalf("Error getting system prompt: %v", err)
	}

	// Check that we got a non-empty system prompt
	if systemPrompt == "" {
		t.Errorf("Expected non-empty system prompt, got empty string")
	}

	// Check that the system prompt contains the agent name
	if !strings.Contains(systemPrompt, "Otto") {
		t.Errorf("Expected system prompt to contain 'Otto', got '%s'", systemPrompt)
	}
}

// TestSaveLoadPersonalities tests saving and loading agent personalities
func TestSaveLoadPersonalities(t *testing.T) {
	suite := SetupAgentsSuite(t)
	defer suite.TearDownAgentsSuite()

	// Initialize agents module
	err := agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Create a new personality
	newPersonality := &agents.AgentPersonality{
		Name:             "TestAgent",
		Description:      "A test agent personality",
		PromptTemplate:   "{{context}}User: {{input}}\n\nTestAgent: ",
		SystemPrompt:     "You are TestAgent, a test agent personality.",
		Traits:           map[string]string{"Test": "High"},
		PreferredModel:   "llama3.3:latest",
		ResponseTemplate: "{{response}}",
	}

	// Add the personality
	err = agents.AddPersonality(newPersonality)
	if err != nil {
		t.Fatalf("Error adding personality: %v", err)
	}

	// Set active agent to the new personality
	err = agents.SetActiveAgent("TestAgent")
	if err != nil {
		t.Fatalf("Error setting active agent: %v", err)
	}

	// Shutdown agents module
	agents.Shutdown()

	// Initialize agents module again
	err = agents.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing agents module: %v", err)
	}

	// Get personalities
	personalities := agents.GetPersonalities()

	// Check that the new personality was loaded
	if _, ok := personalities["TestAgent"]; !ok {
		t.Errorf("Expected 'TestAgent' personality to be loaded, but it was not found")
	}

	// Get active agent
	activeAgent, err := agents.GetActiveAgent()
	if err != nil {
		t.Fatalf("Error getting active agent after reloading: %v", err)
	}

	// Check that the active agent is the one we set
	if activeAgent.Name != "TestAgent" {
		t.Errorf("Expected active agent name to be 'TestAgent', got '%s'", activeAgent.Name)
	}
}
