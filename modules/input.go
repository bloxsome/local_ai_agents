// Package modules provides functionality for handling user input
package modules

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/peterh/liner"
)

// InputResult represents the result of getting user input
type InputResult struct {
	Text string
	Exit bool
	Skip bool
}

// InputHandler handles user input
type InputHandler struct {
	config      interface{}
	logger      *Logger
	liner       *liner.State
	historyFile string
}

// NewInputHandler creates a new InputHandler
func NewInputHandler(config interface{}) *InputHandler {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)

	// Create history file path in user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		GetLogger().Error("Could not get user home directory: %v", err)
		homeDir = "."
	}
	historyFile := filepath.Join(homeDir, ".local_ai_agents_history")

	// Load history file if it exists
	if f, err := os.Open(historyFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	// Set up tab completion for slash commands
	line.SetCompleter(func(line string) (c []string) {
		if strings.HasPrefix(line, "/") {
			for cmd := range GetSlashCommands() {
				if strings.HasPrefix(cmd, line) {
					c = append(c, cmd)
				}
			}
		}
		return
	})

	return &InputHandler{
		config:      config,
		logger:      GetLogger(),
		liner:       line,
		historyFile: historyFile,
	}
}

// Close closes the input handler and saves history
func (h *InputHandler) Close() {
	if h.liner != nil {
		// Save history to file
		if f, err := os.Create(h.historyFile); err == nil {
			h.liner.WriteHistory(f)
			f.Close()
		} else {
			h.logger.Error("Could not write history file: %v", err)
		}

		h.liner.Close()
	}
}

// GetUserInput gets input from the user
func (h *InputHandler) GetUserInput() *InputResult {
	// Get username from config
	var username string
	if cfg, ok := h.config.(map[string]interface{}); ok {
		if u, ok := cfg["UserName"].(string); ok {
			username = u
		} else {
			username = "User"
		}
	} else {
		username = "User"
	}

	prompt := fmt.Sprintf("%s> ", username)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nOperation cancelled by user.")
		h.Close() // Ensure history is saved
		os.Exit(0)
	}()

	// Get input from user
	input, err := h.liner.Prompt(prompt)
	if err != nil {
		if err == liner.ErrPromptAborted {
			h.logger.Info("Input interrupted by user (Ctrl+C)")
			fmt.Println("\nOperation cancelled by user.")
			return &InputResult{Exit: true}
		}
		h.logger.Error("Error getting input: %v", err)
		fmt.Println("\nError getting input.")
		return &InputResult{Exit: true}
	}

	// Add to history if non-empty
	if input != "" {
		h.liner.AppendHistory(input)
	}

	// Handle slash commands
	if strings.HasPrefix(input, "/") {
		h.logger.Info("Slash command received: %s", input)

		// Handle slash commands through the commands package
		result := HandleSlashCommand(input)

		if result == "EXIT" {
			h.logger.Info("Exit command received")
			return &InputResult{Exit: true}
		}

		// For all other slash commands
		return &InputResult{Skip: true}
	}

	h.logger.Debug("User input received: %s...", truncateString(input, 50))
	return &InputResult{Text: input}
}

// GetUserInputSimple gets input from the user using a simple scanner
// This is a fallback in case the liner library is not available
func GetUserInputSimple() *InputResult {
	logger := GetLogger()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Input> ")
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			logger.Error("Error reading input: %v", err)
			return &InputResult{Exit: true}
		}
		// EOF
		logger.Info("EOF detected, exiting program")
		return &InputResult{Exit: true}
	}

	input := scanner.Text()

	// Handle slash commands
	if strings.HasPrefix(input, "/") {
		logger.Info("Slash command received: %s", input)

		// Handle slash commands through the commands package
		result := HandleSlashCommand(input)

		if result == "EXIT" {
			logger.Info("Exit command received")
			return &InputResult{Exit: true}
		}

		// For all other slash commands
		return &InputResult{Skip: true}
	}

	logger.Debug("User input received: %s...", truncateString(input, 50))
	return &InputResult{Text: input}
}

// Helper function to truncate a string
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// RunInputTest runs a test loop for the input module
func RunInputTest() {
	fmt.Println("Type /h or /help for available commands.")
	fmt.Println("Using simple input method for testing...")

	// Use the simple input method for testing
	for {
		fmt.Print("Input> ")
		result := GetUserInputSimple()

		if result.Exit {
			break
		}

		if result.Skip {
			continue
		}

		fmt.Printf("You entered: %s\n", result.Text)
	}

	fmt.Println("Input test completed")
}

// RunAdvancedInputTest runs a test loop for the input module using the advanced input handler
func RunAdvancedInputTest() {
	cfg := map[string]interface{}{
		"UserName": "User",
	}

	handler := NewInputHandler(cfg)
	defer handler.Close()

	fmt.Println("Type /h or /help for available commands.")
	fmt.Println("Using advanced input method with history...")

	for {
		result := handler.GetUserInput()

		if result.Exit {
			break
		}

		if result.Skip {
			continue
		}

		fmt.Printf("You entered: %s\n", result.Text)
	}

	fmt.Println("Input test completed")
}
