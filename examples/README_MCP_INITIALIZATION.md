# MCP Initialization Example

This directory contains an example of how to initialize and use the Model Context Protocol (MCP) module with Ollama.

## Overview

The Model Context Protocol (MCP) allows AI models to interact with external systems through a standardized interface. This example demonstrates how to:

1. Initialize the MCP module
2. Add and start an MCP server
3. Process MCP commands in prompts
4. Send processed prompts to Ollama

## Files

- `mcp_initialization_example.go`: Go example that demonstrates MCP initialization and usage
- `mcp_browser_server/`: Directory containing the MCP browser server implementation
- `test_page.html`: HTML page that is opened by the browser server

## How It Works

The example performs the following steps:

1. Sets up logging
2. Loads configuration
3. Initializes the MCP module
4. Adds the browser MCP server
5. Starts the browser MCP server
6. Lists available MCP servers
7. Creates a prompt with MCP tool usage
8. Processes MCP commands in the prompt
9. Sends the processed prompt to Ollama
10. Stops the browser MCP server
11. Shuts down the MCP module

## Building and Running

To build and run the example:

```bash
cd /Users/C283410/Documents/Repos/Ollama_Agents
go build -o mcp_example ollama_agents/examples/mcp_initialization_example.go
./mcp_example
```

## Expected Output

If the MCP module is initialized correctly, you should see:

1. The browser MCP server being added and started
2. The original prompt and the processed prompt (with the MCP command replaced by the result)
3. The browser opening the test page
4. The Ollama response

## Troubleshooting

If you encounter errors:

1. Check that the MCP module is properly initialized
2. Verify that the browser MCP server is built and accessible
3. Check the log file (`mcp_example.log`) for detailed error messages

## Integration with Ollama

To integrate MCP with Ollama in your own application:

1. Load configuration using `config.LoadConfig()`
2. Initialize the MCP module using `mcp.Initialize(cfg)`
3. Add and start your MCP servers
4. Process MCP commands in prompts using `mcp.ProcessMCPInPrompt(prompt)`
5. Send the processed prompt to Ollama using `ollama.ProcessPrompt(processedPrompt, model, username)`
6. Stop your MCP servers and shut down the MCP module when done

## Notes

- The MCP module must be explicitly initialized before use
- MCP servers must be added and started before they can be used
- MCP commands in prompts are processed by the `ProcessMCPInPrompt` function
- The processed prompt is then sent to Ollama
