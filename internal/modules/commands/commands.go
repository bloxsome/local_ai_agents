// Package commands provides functionality for handling slash commands
package commands

import (
	"fmt"
	"strings"

	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
	"github.com/bloxsome/local_ai_agents/internal/modules/mcp"
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

	// MCP commands
	"/mcp":           mcpHelpCommand,
	"/mcp-list":      mcpListCommand,
	"/mcp-start":     mcpStartCommand,
	"/mcp-stop":      mcpStopCommand,
	"/mcp-enable":    mcpEnableCommand,
	"/mcp-disable":   mcpDisableCommand,
	"/mcp-add":       mcpAddCommand,
	"/mcp-remove":    mcpRemoveCommand,
	"/mcp-tools":     mcpToolsCommand,
	"/mcp-resources": mcpResourcesCommand,
}

// GetSlashCommands returns the map of available slash commands
func GetSlashCommands() map[string]CommandHandler {
	return slashCommands
}

// HandleSlashCommand processes a slash command and returns a result string
func HandleSlashCommand(command string) string {
	logger := logging.GetLogger()

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

MCP Commands:
/mcp           - Show MCP help
/mcp-list      - List all MCP servers
/mcp-start     - Start an MCP server
/mcp-stop      - Stop an MCP server
/mcp-enable    - Enable an MCP server
/mcp-disable   - Disable an MCP server
/mcp-add       - Add a new MCP server
/mcp-remove    - Remove an MCP server
/mcp-tools     - List tools provided by an MCP server
/mcp-resources - List resources provided by an MCP server
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

// MCP command handlers

// mcpHelpCommand displays MCP help information
func mcpHelpCommand(command string) string {
	helpText := `
MCP Commands:
/mcp-list                  - List all MCP servers
/mcp-start <server>        - Start an MCP server
/mcp-stop <server>         - Stop an MCP server
/mcp-enable <server>       - Enable an MCP server
/mcp-disable <server>      - Disable an MCP server
/mcp-add <name> <command> [args...] - Add a new MCP server
/mcp-remove <server>       - Remove an MCP server
/mcp-tools <server>        - List tools provided by an MCP server
/mcp-resources <server>    - List resources provided by an MCP server
`
	fmt.Println(helpText)
	return "CONTINUE"
}

// mcpListCommand lists all MCP servers
func mcpListCommand(command string) string {
	if mcp.DefaultManager == nil {
		fmt.Println("MCP module not initialized")
		return "CONTINUE"
	}

	servers := mcp.DefaultManager.ListServers()
	if len(servers) == 0 {
		fmt.Println("No MCP servers configured")
		return "CONTINUE"
	}

	fmt.Println("MCP Servers:")
	for _, server := range servers {
		status := "disabled"
		if !server.Disabled {
			running, _ := mcp.DefaultManager.IsServerRunning(server.Name)
			if running {
				status = "running"
			} else {
				status = "stopped"
			}
		}
		fmt.Printf("- %s (%s)\n", server.Name, status)
	}

	return "CONTINUE"
}

// mcpStartCommand starts an MCP server
func mcpStartCommand(command string) string {
	if mcp.DefaultManager == nil {
		fmt.Println("MCP module not initialized")
		return "CONTINUE"
	}

	parts := strings.Fields(command)
	if len(parts) < 2 {
		fmt.Println("Usage: /mcp-start <server>")
		return "CONTINUE"
	}

	serverName := parts[1]
	if err := mcp.DefaultManager.StartServer(serverName); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	} else {
		fmt.Printf("Server '%s' started\n", serverName)
	}

	return "CONTINUE"
}

// mcpStopCommand stops an MCP server
func mcpStopCommand(command string) string {
	if mcp.DefaultManager == nil {
		fmt.Println("MCP module not initialized")
		return "CONTINUE"
	}

	parts := strings.Fields(command)
	if len(parts) < 2 {
		fmt.Println("Usage: /mcp-stop <server>")
		return "CONTINUE"
	}

	serverName := parts[1]
	if err := mcp.DefaultManager.StopServer(serverName); err != nil {
		fmt.Printf("Error stopping server: %v\n", err)
	} else {
		fmt.Printf("Server '%s' stopped\n", serverName)
	}

	return "CONTINUE"
}

// mcpEnableCommand enables an MCP server
func mcpEnableCommand(command string) string {
	if mcp.DefaultManager == nil {
		fmt.Println("MCP module not initialized")
		return "CONTINUE"
	}

	parts := strings.Fields(command)
	if len(parts) < 2 {
		fmt.Println("Usage: /mcp-enable <server>")
		return "CONTINUE"
	}

	serverName := parts[1]
	if err := mcp.DefaultManager.EnableServer(serverName); err != nil {
		fmt.Printf("Error enabling server: %v\n", err)
	} else {
		fmt.Printf("Server '%s' enabled\n", serverName)
	}

	return "CONTINUE"
}

// mcpDisableCommand disables an MCP server
func mcpDisableCommand(command string) string {
	if mcp.DefaultManager == nil {
		fmt.Println("MCP module not initialized")
		return "CONTINUE"
	}

	parts := strings.Fields(command)
	if len(parts) < 2 {
		fmt.Println("Usage: /mcp-disable <server>")
		return "CONTINUE"
	}

	serverName := parts[1]
	if err := mcp.DefaultManager.DisableServer(serverName); err != nil {
		fmt.Printf("Error disabling server: %v\n", err)
	} else {
		fmt.Printf("Server '%s' disabled\n", serverName)
	}

	return "CONTINUE"
}

// mcpAddCommand adds a new MCP server
func mcpAddCommand(command string) string {
	if mcp.DefaultManager == nil {
		fmt.Println("MCP module not initialized")
		return "CONTINUE"
	}

	parts := strings.Fields(command)
	if len(parts) < 3 {
		fmt.Println("Usage: /mcp-add <name> <command> [args...]")
		return "CONTINUE"
	}

	serverName := parts[1]
	serverCommand := parts[2]
	var serverArgs []string
	if len(parts) > 3 {
		serverArgs = parts[3:]
	}

	if err := mcp.DefaultManager.AddServer(serverName, serverCommand, serverArgs, nil); err != nil {
		fmt.Printf("Error adding server: %v\n", err)
	} else {
		fmt.Printf("Server '%s' added\n", serverName)
	}

	return "CONTINUE"
}

// mcpRemoveCommand removes an MCP server
func mcpRemoveCommand(command string) string {
	if mcp.DefaultManager == nil {
		fmt.Println("MCP module not initialized")
		return "CONTINUE"
	}

	parts := strings.Fields(command)
	if len(parts) < 2 {
		fmt.Println("Usage: /mcp-remove <server>")
		return "CONTINUE"
	}

	serverName := parts[1]
	if err := mcp.DefaultManager.RemoveServer(serverName); err != nil {
		fmt.Printf("Error removing server: %v\n", err)
	} else {
		fmt.Printf("Server '%s' removed\n", serverName)
	}

	return "CONTINUE"
}

// mcpToolsCommand lists tools provided by an MCP server
func mcpToolsCommand(command string) string {
	if mcp.DefaultManager == nil {
		fmt.Println("MCP module not initialized")
		return "CONTINUE"
	}

	parts := strings.Fields(command)
	if len(parts) < 2 {
		fmt.Println("Usage: /mcp-tools <server>")
		return "CONTINUE"
	}

	serverName := parts[1]
	tools, err := mcp.DefaultManager.ListTools(serverName)
	if err != nil {
		fmt.Printf("Error listing tools: %v\n", err)
		return "CONTINUE"
	}

	if len(tools) == 0 {
		fmt.Printf("No tools provided by server '%s'\n", serverName)
		return "CONTINUE"
	}

	fmt.Printf("Tools provided by server '%s':\n", serverName)
	for _, tool := range tools {
		name := tool["name"].(string)
		description := ""
		if desc, ok := tool["description"].(string); ok {
			description = desc
		}
		fmt.Printf("- %s: %s\n", name, description)
	}

	return "CONTINUE"
}

// mcpResourcesCommand lists resources provided by an MCP server
func mcpResourcesCommand(command string) string {
	if mcp.DefaultManager == nil {
		fmt.Println("MCP module not initialized")
		return "CONTINUE"
	}

	parts := strings.Fields(command)
	if len(parts) < 2 {
		fmt.Println("Usage: /mcp-resources <server>")
		return "CONTINUE"
	}

	serverName := parts[1]
	resources, err := mcp.DefaultManager.ListResources(serverName)
	if err != nil {
		fmt.Printf("Error listing resources: %v\n", err)
		return "CONTINUE"
	}

	if len(resources) == 0 {
		fmt.Printf("No resources provided by server '%s'\n", serverName)
		return "CONTINUE"
	}

	fmt.Printf("Resources provided by server '%s':\n", serverName)
	for _, resource := range resources {
		uri := resource["uri"].(string)
		name := ""
		if n, ok := resource["name"].(string); ok {
			name = n
		}
		fmt.Printf("- %s: %s\n", uri, name)
	}

	// Also list resource templates
	templates, err := mcp.DefaultManager.ListResourceTemplates(serverName)
	if err == nil && len(templates) > 0 {
		fmt.Printf("\nResource templates provided by server '%s':\n", serverName)
		for _, template := range templates {
			uriTemplate := template["uriTemplate"].(string)
			name := ""
			if n, ok := template["name"].(string); ok {
				name = n
			}
			fmt.Printf("- %s: %s\n", uriTemplate, name)
		}
	}

	return "CONTINUE"
}
