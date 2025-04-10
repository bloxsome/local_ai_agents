// Package mcp provides functionality for Model Context Protocol (MCP) support
package mcp

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/bloxsome/local_ai_agents/internal/modules/logging"
)

// ToolPattern is the regex pattern to match MCP tool usage in prompts
var ToolPattern = regexp.MustCompile(`<use_mcp_tool>\s*<server_name>(.*?)</server_name>\s*<tool_name>(.*?)</tool_name>\s*<arguments>\s*([\s\S]*?)\s*</arguments>\s*</use_mcp_tool>`)

// ResourcePattern is the regex pattern to match MCP resource access in prompts
var ResourcePattern = regexp.MustCompile(`<access_mcp_resource>\s*<server_name>(.*?)</server_name>\s*<uri>(.*?)</uri>\s*</access_mcp_resource>`)

// ProcessMCPInPrompt processes MCP tool usage and resource access in a prompt
func ProcessMCPInPrompt(prompt string) (string, error) {
	logger := logging.GetLogger()

	// Process tool usage
	toolMatches := ToolPattern.FindAllStringSubmatch(prompt, -1)
	for _, match := range toolMatches {
		if len(match) != 4 {
			continue
		}

		fullMatch := match[0]
		serverName := strings.TrimSpace(match[1])
		toolName := strings.TrimSpace(match[2])
		argumentsStr := strings.TrimSpace(match[3])

		logger.Info("Processing MCP tool usage: server=%s, tool=%s", serverName, toolName)

		// Parse arguments
		var arguments map[string]interface{}
		if err := json.Unmarshal([]byte(argumentsStr), &arguments); err != nil {
			logger.Error("Failed to parse MCP tool arguments: %v", err)
			replacement := fmt.Sprintf("Error: Failed to parse MCP tool arguments: %v", err)
			prompt = strings.Replace(prompt, fullMatch, replacement, 1)
			continue
		}

		// Use the tool
		if DefaultManager == nil {
			logger.Error("MCP module not initialized")
			replacement := "Error: MCP module not initialized"
			prompt = strings.Replace(prompt, fullMatch, replacement, 1)
			continue
		}

		result, err := DefaultManager.UseTool(serverName, toolName, arguments)
		if err != nil {
			logger.Error("Failed to use MCP tool: %v", err)
			replacement := fmt.Sprintf("Error: Failed to use MCP tool: %v", err)
			prompt = strings.Replace(prompt, fullMatch, replacement, 1)
			continue
		}

		// Replace the tool usage with the result
		prompt = strings.Replace(prompt, fullMatch, result, 1)
	}

	// Process resource access
	resourceMatches := ResourcePattern.FindAllStringSubmatch(prompt, -1)
	for _, match := range resourceMatches {
		if len(match) != 3 {
			continue
		}

		fullMatch := match[0]
		serverName := strings.TrimSpace(match[1])
		uri := strings.TrimSpace(match[2])

		logger.Info("Processing MCP resource access: server=%s, uri=%s", serverName, uri)

		// Access the resource
		if DefaultManager == nil {
			logger.Error("MCP module not initialized")
			replacement := "Error: MCP module not initialized"
			prompt = strings.Replace(prompt, fullMatch, replacement, 1)
			continue
		}

		result, err := DefaultManager.AccessResource(serverName, uri)
		if err != nil {
			logger.Error("Failed to access MCP resource: %v", err)
			replacement := fmt.Sprintf("Error: Failed to access MCP resource: %v", err)
			prompt = strings.Replace(prompt, fullMatch, replacement, 1)
			continue
		}

		// Replace the resource access with the result
		prompt = strings.Replace(prompt, fullMatch, result, 1)
	}

	return prompt, nil
}

// GetConnectedMCPServersInfo returns information about connected MCP servers
func GetConnectedMCPServersInfo() string {
	if DefaultManager == nil {
		return "(No MCP servers currently connected)"
	}

	servers := DefaultManager.ListServers()
	if len(servers) == 0 {
		return "(No MCP servers currently connected)"
	}

	var sb strings.Builder
	sb.WriteString("# Connected MCP Servers\n\n")
	sb.WriteString("When a server is connected, you can use the server's tools via the `use_mcp_tool` tool, and access the server's resources via the `access_mcp_resource` tool.\n\n")

	for _, server := range servers {
		if server.Disabled {
			continue
		}

		running, _ := DefaultManager.IsServerRunning(server.Name)
		if !running {
			continue
		}

		sb.WriteString(fmt.Sprintf("## %s\n\n", server.Name))

		// List tools
		tools, err := DefaultManager.ListTools(server.Name)
		if err == nil && len(tools) > 0 {
			sb.WriteString("### Tools\n\n")
			for _, tool := range tools {
				name := tool["name"].(string)
				description := ""
				if desc, ok := tool["description"].(string); ok {
					description = desc
				}
				sb.WriteString(fmt.Sprintf("- `%s`: %s\n", name, description))
			}
			sb.WriteString("\n")
		}

		// List resources
		resources, err := DefaultManager.ListResources(server.Name)
		if err == nil && len(resources) > 0 {
			sb.WriteString("### Resources\n\n")
			for _, resource := range resources {
				uri := resource["uri"].(string)
				name := ""
				if n, ok := resource["name"].(string); ok {
					name = n
				}
				sb.WriteString(fmt.Sprintf("- `%s`: %s\n", uri, name))
			}
			sb.WriteString("\n")
		}

		// List resource templates
		templates, err := DefaultManager.ListResourceTemplates(server.Name)
		if err == nil && len(templates) > 0 {
			sb.WriteString("### Resource Templates\n\n")
			for _, template := range templates {
				uriTemplate := template["uriTemplate"].(string)
				name := ""
				if n, ok := template["name"].(string); ok {
					name = n
				}
				sb.WriteString(fmt.Sprintf("- `%s`: %s\n", uriTemplate, name))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
