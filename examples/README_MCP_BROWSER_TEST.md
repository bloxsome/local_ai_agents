# Testing MCP Browser Support in Ollama

This directory contains files for testing Model Context Protocol (MCP) browser support in Ollama.

## Overview

The Model Context Protocol (MCP) allows AI models to interact with external systems through a standardized interface. This test demonstrates how Ollama can use MCP to open URLs in a web browser.

## Files

- `mcp_browser_server/`: Directory containing the MCP browser server implementation
  - `browser.go`: Go implementation of the MCP browser server
  - `browser`: Compiled binary of the MCP browser server
  - `README.md`: Documentation for the MCP browser server
- `mcp_settings.json`: MCP settings file that configures the browser server
- `mcp_browser_test.py`: Python script that tests the MCP browser support
- `test_page.html`: HTML page that is opened by the test

## Prerequisites

- Ollama installed and running
- Go (for building the MCP browser server)
- Python 3 with the `requests` library

## Running the Test

1. Build the MCP browser server:
   ```bash
   cd ollama_agents/examples/mcp_browser_server
   go build -o browser browser.go
   ```

2. Run the test script:
   ```bash
   cd /Users/C283410/Documents/Repos/Ollama_Agents
   python3 ollama_agents/examples/mcp_browser_test.py
   ```

3. The test will:
   - Copy the MCP settings to the correct location
   - Send a prompt to Ollama that uses the MCP browser tool
   - Open the test page in your default browser
   - Clean up the MCP settings

## How It Works

1. The test script sends a prompt to Ollama that includes an MCP tool usage:
   ```
   <use_mcp_tool>
   <server_name>browser</server_name>
   <tool_name>open_url</tool_name>
   <arguments>
   {
     "url": "file:///path/to/test_page.html"
   }
   </arguments>
   </use_mcp_tool>
   ```

2. Ollama processes the prompt and detects the MCP tool usage.

3. The MCP module in Ollama forwards the tool usage to the browser MCP server.

4. The browser MCP server opens the URL in the default browser.

5. The test page is displayed, confirming that the MCP browser support is working.

## Extending the Test

You can extend this test by:

- Adding more tools to the MCP browser server (e.g., get page title, take screenshot)
- Creating more complex prompts that use multiple MCP tools
- Integrating with other MCP servers for more functionality
