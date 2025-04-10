# MCP Module Documentation

The MCP (Model Context Protocol) module provides support for connecting to and using MCP servers in Ollama Agents. This allows the AI to access external tools and resources through a standardized protocol.

## Overview

MCP is a protocol that enables AI models to interact with external systems through a standardized interface. It allows models to:

1. Use tools provided by MCP servers to perform actions
2. Access resources provided by MCP servers to get information
3. Discover available tools and resources

The MCP module in Ollama Agents provides the infrastructure to connect to MCP servers, manage their lifecycle, and use their tools and resources in conversations.

## Configuration

MCP support can be configured through environment variables:

- `AI_MCP_ENABLED`: Enable or disable MCP support (default: `true`)
- `AI_MCP_DIR`: Directory for MCP server files (default: `~/Documents/Cline/MCP`)

## MCP Settings

MCP server settings are stored in a JSON file at one of the following locations (in order of precedence):

1. `~/.config/ollama_agents/mcp_settings.json`
2. `~/Library/Application Support/Ollama Agents/mcp_settings.json`
3. `<project_root>/config/mcp_settings.json`

The settings file has the following structure:

```json
{
  "mcpServers": {
    "server1": {
      "name": "server1",
      "command": "node",
      "args": ["/path/to/server.js"],
      "env": {
        "API_KEY": "your-api-key"
      },
      "disabled": false,
      "autoApprove": []
    },
    "server2": {
      "name": "server2",
      "command": "python",
      "args": ["/path/to/server.py"],
      "env": {},
      "disabled": true,
      "autoApprove": []
    }
  }
}
```

## Commands

The MCP module provides the following slash commands:

- `/mcp`: Show MCP help
- `/mcp-list`: List all MCP servers
- `/mcp-start <server>`: Start an MCP server
- `/mcp-stop <server>`: Stop an MCP server
- `/mcp-enable <server>`: Enable an MCP server
- `/mcp-disable <server>`: Disable an MCP server
- `/mcp-add <name> <command> [args...]`: Add a new MCP server
- `/mcp-remove <server>`: Remove an MCP server
- `/mcp-tools <server>`: List tools provided by an MCP server
- `/mcp-resources <server>`: List resources provided by an MCP server

## Using MCP in Conversations

### Using MCP Tools

To use an MCP tool in a conversation, use the following syntax:

```
<use_mcp_tool>
<server_name>server_name</server_name>
<tool_name>tool_name</tool_name>
<arguments>
{
  "param1": "value1",
  "param2": "value2"
}
</arguments>
</use_mcp_tool>
```

This will be replaced with the result of the tool execution.

### Accessing MCP Resources

To access an MCP resource in a conversation, use the following syntax:

```
<access_mcp_resource>
<server_name>server_name</server_name>
<uri>resource_uri</uri>
</access_mcp_resource>
```

This will be replaced with the content of the resource.

## Creating MCP Servers

MCP servers can be implemented in any programming language that can read from stdin and write to stdout. The server should implement the JSON-RPC 2.0 protocol with the following methods:

- `listTools`: List the tools provided by the server
- `callTool`: Call a tool with arguments
- `listResources`: List the resources provided by the server
- `listResourceTemplates`: List the resource templates provided by the server
- `readResource`: Read a resource by URI

See the [example MCP server](../examples/mcp_server/calculator.go) for a simple implementation.

## Example

Here's an example of using the calculator MCP server:

1. Add the server:
   ```
   /mcp-add calculator go run ollama_agents/examples/mcp_server/calculator.go
   ```

2. Start the server:
   ```
   /mcp-start calculator
   ```

3. Use the calculator in a conversation:
   ```
   What is 5 + 3?
   <use_mcp_tool>
   <server_name>calculator</server_name>
   <tool_name>add</tool_name>
   <arguments>
   {
     "a": 5,
     "b": 3
   }
   </arguments>
   </use_mcp_tool>
   ```

   The AI will replace this with: "What is 5 + 3? 8.000000"

## Implementation Details

The MCP module consists of the following components:

- `mcp.go`: Core MCP functionality for managing servers
- `integration.go`: Integration with the Ollama client for processing MCP commands in prompts

The module uses the following process to handle MCP commands:

1. The Ollama client intercepts prompts and passes them to the MCP module
2. The MCP module processes any MCP commands in the prompt
3. The processed prompt is sent to the Ollama API
4. The response is returned to the user

## Security Considerations

MCP servers can execute arbitrary code, so they should be treated with caution. Only add MCP servers from trusted sources. The MCP module provides the following security features:

- Servers can be disabled to prevent them from being used
- Servers can be stopped when not in use
- Environment variables can be used to pass sensitive information to servers
