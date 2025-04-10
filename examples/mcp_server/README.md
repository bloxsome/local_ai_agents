# Calculator MCP Server Example

This is a simple example of an MCP (Model Context Protocol) server that provides calculator functionality.

## Features

- Basic arithmetic operations (add, subtract, multiply, divide)
- Mathematical constants (pi, e)

## Building

To build the calculator MCP server:

```bash
cd ollama_agents/examples/mcp_server
go build -o calculator calculator.go
```

## Usage

### Adding to Ollama Agents

To add the calculator MCP server to Ollama Agents, use the `/mcp-add` command:

```
/mcp-add calculator go run ollama_agents/examples/mcp_server/calculator.go
```

Or if you've built the binary:

```
/mcp-add calculator ollama_agents/examples/mcp_server/calculator
```

### Starting the Server

To start the calculator MCP server:

```
/mcp-start calculator
```

### Listing Available Tools

To list the tools provided by the calculator MCP server:

```
/mcp-tools calculator
```

### Using the Calculator in Conversations

Once the calculator MCP server is running, you can use it in conversations with the following syntax:

```
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

This will be replaced with the result: `8.000000`

### Available Tools

- `add`: Add two numbers
  - Parameters: `a` (first number), `b` (second number)
- `subtract`: Subtract two numbers
  - Parameters: `a` (first number), `b` (second number)
- `multiply`: Multiply two numbers
  - Parameters: `a` (first number), `b` (second number)
- `divide`: Divide two numbers
  - Parameters: `a` (first number), `b` (second number, non-zero)

### Available Resources

- `calculator://constants/pi`: The mathematical constant π (pi)
- `calculator://constants/e`: The mathematical constant e

To access a resource:

```
<access_mcp_resource>
<server_name>calculator</server_name>
<uri>calculator://constants/pi</uri>
</access_mcp_resource>
```

This will be replaced with the value of pi: `3.14159265358979323846`
