# Browser MCP Server Example

This is a simple example of an MCP (Model Context Protocol) server that provides browser functionality.

## Features

- Open URLs in the default browser

## Building

To build the browser MCP server:

```bash
cd ollama_agents/examples/mcp_browser_server
go build -o browser browser.go
```

## Usage

### Adding to Ollama Agents

To add the browser MCP server to Ollama Agents, use the `/mcp-add` command:

```
/mcp-add browser go run ollama_agents/examples/mcp_browser_server/browser.go
```

Or if you've built the binary:

```
/mcp-add browser ollama_agents/examples/mcp_browser_server/browser
```

### Starting the Server

To start the browser MCP server:

```
/mcp-start browser
```

### Listing Available Tools

To list the tools provided by the browser MCP server:

```
/mcp-tools browser
```

### Using the Browser in Conversations

Once the browser MCP server is running, you can use it in conversations with the following syntax:

```
<use_mcp_tool>
<server_name>browser</server_name>
<tool_name>open_url</tool_name>
<arguments>
{
  "url": "https://example.com"
}
</arguments>
</use_mcp_tool>
```

This will open the URL in the default browser and return a confirmation message.

### Available Tools

- `open_url`: Open a URL in the default browser
  - Parameters: `url` (URL to open)
