package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bloxsome/local_ai_agents/internal/config"
	"github.com/bloxsome/local_ai_agents/internal/modules/context"
	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
)

// TestContextSuite represents a test suite for context module
type TestContextSuite struct {
	config *config.Configuration
	logger *logging.Logger
}

// SetupContextSuite sets up the test suite for context tests
func SetupContextSuite(t *testing.T) *TestContextSuite {
	// Set up logging
	logFile := "context_test.log"
	if err := logging.SetupLogging(logFile, "DEBUG"); err != nil {
		t.Fatalf("Error setting up logging: %v", err)
	}
	logger := logging.GetLogger()
	logger.Info("Context test suite started")

	// Create a test configuration
	tempDir, err := os.MkdirTemp("", "context_test")
	if err != nil {
		t.Fatalf("Error creating temp directory: %v", err)
	}

	cfg := &config.Configuration{
		UserName:        "TestUser",
		AgentName:       "TestAgent",
		DefaultModel:    "llama3.3:latest",
		MemoryLength:    10,
		ProjectRoot:     tempDir,
		ChatHistoryFile: filepath.Join(tempDir, "test_chat_history.json"),
	}

	return &TestContextSuite{
		config: cfg,
		logger: logger,
	}
}

// TearDownContextSuite tears down the test suite
func (s *TestContextSuite) TearDownContextSuite() {
	// Clean up log file
	logFile := "context_test.log"
	if err := os.Remove(logFile); err != nil {
		s.logger.Warning("Failed to remove log file: %v", err)
	}

	// Clean up temp directory
	if err := os.RemoveAll(s.config.ProjectRoot); err != nil {
		s.logger.Warning("Failed to remove temp directory: %v", err)
	}

	// Shutdown context module if it was initialized
	if context.DefaultManager != nil {
		context.Shutdown()
	}
}

// TestContextInitialization tests initializing the context module
func TestContextInitialization(t *testing.T) {
	suite := SetupContextSuite(t)
	defer suite.TearDownContextSuite()

	// Initialize context module
	err := context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Check that the default manager was created
	if context.DefaultManager == nil {
		t.Fatalf("Default manager is nil after initialization")
	}
}

// TestContextAddMessage tests adding messages to the context
func TestContextAddMessage(t *testing.T) {
	suite := SetupContextSuite(t)
	defer suite.TearDownContextSuite()

	// Initialize context module
	err := context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Add a user message
	context.AddMessage("user", "Hello, this is a test message")

	// Add an assistant message
	context.AddMessage("assistant", "Hello! I'm responding to your test message")

	// Get the history
	history := context.GetHistory()

	// Check that we have 2 messages
	if len(history) != 2 {
		t.Errorf("Expected 2 messages in history, got %d", len(history))
	}

	// Check the roles and content
	if len(history) >= 2 {
		if history[0].Role != "user" {
			t.Errorf("Expected first message role to be 'user', got '%s'", history[0].Role)
		}
		if history[0].Content != "Hello, this is a test message" {
			t.Errorf("Expected first message content to be 'Hello, this is a test message', got '%s'", history[0].Content)
		}
		if history[1].Role != "assistant" {
			t.Errorf("Expected second message role to be 'assistant', got '%s'", history[1].Role)
		}
		if history[1].Content != "Hello! I'm responding to your test message" {
			t.Errorf("Expected second message content to be 'Hello! I'm responding to your test message', got '%s'", history[1].Content)
		}
	}
}

// TestContextExtraction tests context extraction from messages
func TestContextExtraction(t *testing.T) {
	suite := SetupContextSuite(t)
	defer suite.TearDownContextSuite()

	// Initialize context module
	err := context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Add a user message with multiple sentences
	context.AddMessage("user", "My name is John. I am a software engineer. I like to code in Go.")

	// Get the context
	ctx := context.GetContext()

	// Check that we have extracted context
	if len(ctx) == 0 {
		t.Errorf("Expected non-empty context, got empty context")
	}

	// Check that the context contains the sentences
	foundName := false
	foundProfession := false
	foundLanguage := false

	for _, item := range ctx {
		if item == "My name is John" {
			foundName = true
		}
		if item == "I am a software engineer" {
			foundProfession = true
		}
		if item == "I like to code in Go" {
			foundLanguage = true
		}
	}

	if !foundName {
		t.Errorf("Context does not contain 'My name is John'")
	}
	if !foundProfession {
		t.Errorf("Context does not contain 'I am a software engineer'")
	}
	if !foundLanguage {
		t.Errorf("Context does not contain 'I like to code in Go'")
	}
}

// TestContextClear tests clearing the context
func TestContextClear(t *testing.T) {
	suite := SetupContextSuite(t)
	defer suite.TearDownContextSuite()

	// Initialize context module
	err := context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Add a user message
	context.AddMessage("user", "My name is John. I am a software engineer.")

	// Get the context
	ctx := context.GetContext()

	// Check that we have extracted context
	if len(ctx) == 0 {
		t.Errorf("Expected non-empty context, got empty context")
	}

	// Clear the context
	context.ClearContext()

	// Get the context again
	ctx = context.GetContext()

	// Check that the context is empty
	if len(ctx) != 0 {
		t.Errorf("Expected empty context after clearing, got %d items", len(ctx))
	}
}

// TestBulletPoints tests bullet point management
func TestBulletPoints(t *testing.T) {
	suite := SetupContextSuite(t)
	defer suite.TearDownContextSuite()

	// Initialize context module
	err := context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Add bullet points
	context.AddBullet("First bullet point")
	context.AddBullet("Second bullet point")
	context.AddBullet("Third bullet point")

	// Get the bullet points
	bullets := context.GetBullets()

	// Check that we have 3 bullet points
	if len(bullets) != 3 {
		t.Errorf("Expected 3 bullet points, got %d", len(bullets))
	}

	// Check the bullet point content
	if len(bullets) >= 3 {
		if bullets[0] != "First bullet point" {
			t.Errorf("Expected first bullet point to be 'First bullet point', got '%s'", bullets[0])
		}
		if bullets[1] != "Second bullet point" {
			t.Errorf("Expected second bullet point to be 'Second bullet point', got '%s'", bullets[1])
		}
		if bullets[2] != "Third bullet point" {
			t.Errorf("Expected third bullet point to be 'Third bullet point', got '%s'", bullets[2])
		}
	}

	// Clear the bullet points
	context.ClearBullets()

	// Get the bullet points again
	bullets = context.GetBullets()

	// Check that the bullet points are empty
	if len(bullets) != 0 {
		t.Errorf("Expected empty bullet points after clearing, got %d items", len(bullets))
	}
}

// TestKnowledgeTree tests knowledge tree management
func TestKnowledgeTree(t *testing.T) {
	suite := SetupContextSuite(t)
	defer suite.TearDownContextSuite()

	// Initialize context module
	err := context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Add knowledge
	context.AddKnowledge("Go", "Go is a programming language created by Google.")
	context.AddKnowledge("Go", "Go is known for its simplicity and concurrency features.")
	context.AddKnowledge("Python", "Python is a high-level programming language.")

	// Get the knowledge tree
	tree := context.GetKnowledgeTree()

	// Check that we have 2 concepts
	if len(tree) != 2 {
		t.Errorf("Expected 2 concepts in knowledge tree, got %d", len(tree))
	}

	// Check the knowledge content
	if goKnowledge, ok := tree["Go"]; ok {
		if len(goKnowledge) != 2 {
			t.Errorf("Expected 2 knowledge items for Go, got %d", len(goKnowledge))
		}
	} else {
		t.Errorf("Expected 'Go' concept in knowledge tree, but it was not found")
	}

	if pythonKnowledge, ok := tree["Python"]; ok {
		if len(pythonKnowledge) != 1 {
			t.Errorf("Expected 1 knowledge item for Python, got %d", len(pythonKnowledge))
		}
	} else {
		t.Errorf("Expected 'Python' concept in knowledge tree, but it was not found")
	}
}

// TestUserProfile tests user profile management
func TestUserProfile(t *testing.T) {
	suite := SetupContextSuite(t)
	defer suite.TearDownContextSuite()

	// Initialize context module
	err := context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Update user profile
	context.UpdateUserProfile("name", "John Doe")
	context.UpdateUserProfile("occupation", "Software Engineer")
	context.UpdateUserProfile("interests", "Programming, Reading")

	// Get the user profile
	profile := context.GetUserProfile()

	// Check that we have 3 profile items
	if len(profile) != 3 {
		t.Errorf("Expected 3 profile items, got %d", len(profile))
	}

	// Check the profile content
	if name, ok := profile["name"]; ok {
		if name != "John Doe" {
			t.Errorf("Expected name to be 'John Doe', got '%s'", name)
		}
	} else {
		t.Errorf("Expected 'name' in user profile, but it was not found")
	}

	if occupation, ok := profile["occupation"]; ok {
		if occupation != "Software Engineer" {
			t.Errorf("Expected occupation to be 'Software Engineer', got '%s'", occupation)
		}
	} else {
		t.Errorf("Expected 'occupation' in user profile, but it was not found")
	}

	if interests, ok := profile["interests"]; ok {
		if interests != "Programming, Reading" {
			t.Errorf("Expected interests to be 'Programming, Reading', got '%s'", interests)
		}
	} else {
		t.Errorf("Expected 'interests' in user profile, but it was not found")
	}
}

// TestExplainConcept tests the explain concept functionality
func TestExplainConcept(t *testing.T) {
	suite := SetupContextSuite(t)
	defer suite.TearDownContextSuite()

	// Initialize context module
	err := context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Add knowledge
	context.AddKnowledge("Go", "Go is a programming language created by Google.")
	context.AddKnowledge("Go", "Go is known for its simplicity and concurrency features.")

	// Get explanation for a concept that exists
	explanation := context.ExplainConcept("Go")

	// Check that the explanation contains the knowledge
	if !contains(explanation, "Go is a programming language created by Google") {
		t.Errorf("Expected explanation to contain 'Go is a programming language created by Google', got '%s'", explanation)
	}
	if !contains(explanation, "Go is known for its simplicity and concurrency features") {
		t.Errorf("Expected explanation to contain 'Go is known for its simplicity and concurrency features', got '%s'", explanation)
	}

	// Get explanation for a concept that doesn't exist
	explanation = context.ExplainConcept("Java")

	// Check that the explanation indicates no knowledge
	if !contains(explanation, "don't have specific knowledge about Java") {
		t.Errorf("Expected explanation to indicate no knowledge about Java, got '%s'", explanation)
	}
}

// TestFactCheck tests the fact check functionality
func TestFactCheck(t *testing.T) {
	suite := SetupContextSuite(t)
	defer suite.TearDownContextSuite()

	// Initialize context module
	err := context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Perform a fact check
	result := context.FactCheck("Go was created by Google")

	// Check that the result contains the statement
	if !contains(result, "Go was created by Google") {
		t.Errorf("Expected fact check result to contain the statement, got '%s'", result)
	}
}

// TestFormatContextForPrompt tests formatting context for inclusion in a prompt
func TestFormatContextForPrompt(t *testing.T) {
	suite := SetupContextSuite(t)
	defer suite.TearDownContextSuite()

	// Initialize context module
	err := context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Add context
	context.AddMessage("user", "My name is John. I am a software engineer.")

	// Add bullet points
	context.AddBullet("First bullet point")
	context.AddBullet("Second bullet point")

	// Update user profile
	context.UpdateUserProfile("name", "John Doe")
	context.UpdateUserProfile("occupation", "Software Engineer")

	// Format context for prompt
	formattedContext := context.FormatContextForPrompt()

	// Check that the formatted context contains the context
	if !contains(formattedContext, "Context:") {
		t.Errorf("Expected formatted context to contain 'Context:', got '%s'", formattedContext)
	}
	if !contains(formattedContext, "My name is John") {
		t.Errorf("Expected formatted context to contain 'My name is John', got '%s'", formattedContext)
	}

	// Check that the formatted context contains the bullet points
	if !contains(formattedContext, "Key Points:") {
		t.Errorf("Expected formatted context to contain 'Key Points:', got '%s'", formattedContext)
	}
	if !contains(formattedContext, "First bullet point") {
		t.Errorf("Expected formatted context to contain 'First bullet point', got '%s'", formattedContext)
	}
	if !contains(formattedContext, "Second bullet point") {
		t.Errorf("Expected formatted context to contain 'Second bullet point', got '%s'", formattedContext)
	}

	// Check that the formatted context contains the user profile
	if !contains(formattedContext, "User Profile:") {
		t.Errorf("Expected formatted context to contain 'User Profile:', got '%s'", formattedContext)
	}
	if !contains(formattedContext, "name: John Doe") {
		t.Errorf("Expected formatted context to contain 'name: John Doe', got '%s'", formattedContext)
	}
	if !contains(formattedContext, "occupation: Software Engineer") {
		t.Errorf("Expected formatted context to contain 'occupation: Software Engineer', got '%s'", formattedContext)
	}
}

// TestSaveLoadHistory tests saving and loading conversation history
func TestSaveLoadHistory(t *testing.T) {
	suite := SetupContextSuite(t)
	defer suite.TearDownContextSuite()

	// Initialize context module
	err := context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Add messages
	context.AddMessage("user", "Hello, this is a test message")
	context.AddMessage("assistant", "Hello! I'm responding to your test message")

	// Save history
	err = context.SaveHistory()
	if err != nil {
		t.Fatalf("Error saving history: %v", err)
	}

	// Check that the history file exists
	if _, err := os.Stat(suite.config.ChatHistoryFile); os.IsNotExist(err) {
		t.Errorf("Expected chat history file to exist, but it does not")
	}

	// Shutdown context module
	context.Shutdown()

	// Initialize context module again
	err = context.Initialize(suite.config)
	if err != nil {
		t.Fatalf("Error initializing context module: %v", err)
	}

	// Get the history
	history := context.GetHistory()

	// Check that we have 2 messages
	if len(history) != 2 {
		t.Errorf("Expected 2 messages in history after loading, got %d", len(history))
	}

	// Check the roles and content
	if len(history) >= 2 {
		if history[0].Role != "user" {
			t.Errorf("Expected first message role to be 'user', got '%s'", history[0].Role)
		}
		if history[0].Content != "Hello, this is a test message" {
			t.Errorf("Expected first message content to be 'Hello, this is a test message', got '%s'", history[0].Content)
		}
		if history[1].Role != "assistant" {
			t.Errorf("Expected second message role to be 'assistant', got '%s'", history[1].Role)
		}
		if history[1].Content != "Hello! I'm responding to your test message" {
			t.Errorf("Expected second message content to be 'Hello! I'm responding to your test message', got '%s'", history[1].Content)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return s != "" && strings.Contains(s, substr)
}
