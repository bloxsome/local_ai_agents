// Package modules provides functionality for handling slash commands
package modules

import (
	"fmt"
	"strings"
)

// CommandHandler is a function that handles a slash command
type CommandHandler func(string) string

// slashCommands is a map of command names to their handler functions
var slashCommands = map[string]CommandHandler{
	"/h":    helpCommand,
	"/help": helpCommand,
	"/e":    exitCommand,
	"/exit": exitCommand,
	"/q":    exitCommand,
	"/quit": exitCommand,
}

// GetSlashCommands returns the map of available slash commands
func GetSlashCommands() map[string]CommandHandler {
	return slashCommands
}

// HandleSlashCommand processes a slash command and returns a result string
func HandleSlashCommand(command string) string {
	logger := GetLogger()

	cmdParts := strings.Fields(command)
	if len(cmdParts) == 0 {
		return "CONTINUE"
	}

	cmd := cmdParts[0]
	handler, exists := slashCommands[cmd]

	if exists {
		logger.Info("Executing command: %s", cmd)
		return handler(command)
	}

	// Handle unknown command
	logger.Warning("Unknown command received: %s", command)
	fmt.Printf("Unknown command: %s\n", command)
	return "CONTINUE"
}

// helpCommand displays help information
func helpCommand(command string) string {
	helpText := `
Available commands:
/h or /help - Show this help message
/e or /exit - Exit the program
/q or /quit - Exit the program
`
	fmt.Println(helpText)
	return "CONTINUE"
}

// exitCommand exits the program
func exitCommand(command string) string {
	fmt.Println("Exiting program.")
	return "EXIT"
}

// RegisterCommand adds a new command handler to the slash commands map
func RegisterCommand(name string, handler CommandHandler) {
	slashCommands[name] = handler
}
