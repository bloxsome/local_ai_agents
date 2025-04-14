# MCP Generator Server

An MCP server that allows AI models to create their own MCP servers.

## Overview

The MCP Generator server provides tools for AI models to:

1. Generate new MCP servers based on specifications
2. Validate MCP server specifications
3. Build MCP servers from source

This enables AI models to extend their capabilities by creating custom MCP servers that can interact with external systems, APIs, and services.

## Building

To build the MCP Generator server:

```bash
cd examples/mcp_server/mcp_generator
go build -o mcp_generator
```

## Usage

### Adding to Local_AI_Agents

To add the MCP Generator server to Local_AI_Agents, use the `/mcp-add` command:

```
/mcp-add mcp_generator go run examples/mcp_server/mcp_generator/main.go
```

Or if you've built the binary:

```
/mcp-add mcp_generator examples/mcp_server/mcp_generator/mcp_generator
```

### Starting the Server

To start the MCP Generator server:

```
/mcp-start mcp_generator
```

### Available Tools

- `generate_mcp_server`: Generate a new MCP server based on a specification
  - Parameters:
    - `spec`: Server specification object
    - `outputDir`: Output directory for the generated server
  
- `validate_mcp_server`: Validate an MCP server specification
  - Parameters:
    - `spec`: Server specification to validate
  
- `build_mcp_server`: Build an MCP server from source
  - Parameters:
    - `serverPath`: Path to the server source directory

### Available Resources

- `mcp-generator://templates/server`: MCP Server Template
- `mcp-generator://templates/readme`: MCP Server README Template
- `mcp-generator://examples/calculator`: Calculator MCP Server Example

## Server Specification Format

The server specification is a JSON object with the following structure:

```json
{
  "name": "server_name",
  "description": "Server description",
  "tools": [
    {
      "name": "tool_name",
      "description": "Tool description",
      "inputSchema": {
        "type": "object",
        "properties": {
          "param1": {
            "type": "string",
            "description": "Parameter description"
          }
        },
        "required": ["param1"]
      },
      "logic": "// Go code for tool implementation"
    }
  ],
  "resources": [
    {
      "uri": "protocol://path",
      "name": "Resource name",
      "mimeType": "text/plain",
      "description": "Resource description",
      "content": "Resource content"
    }
  ]
}
```

## Testing with LLMs

The repository includes a test script that demonstrates how LLMs like llama3.3 can use the MCP Generator to create their own MCP servers:

```bash
# Install required Python package
pip install requests

# Run the test script
python examples/mcp_generator_test.py
```

This script:
1. Starts the MCP Generator server
2. Queries llama3.3 to create a specification for a greeting MCP server
3. Uses the MCP Generator to validate, generate, and build the server
4. Outputs the generated server to `./generated_mcp_server` (configurable)

You can customize the test with command-line options:
```bash
python examples/mcp_generator_test.py --model llama3.3:latest --output-dir ./my_server
```

## Example: Creating a Weather MCP Server

Here's an example of how an AI model can use the MCP Generator to create a weather MCP server:

```
<use_mcp_tool>
<server_name>mcp_generator</server_name>
<tool_name>generate_mcp_server</tool_name>
<arguments>
{
  "spec": {
    "name": "weather",
    "description": "A simple weather MCP server",
    "tools": [
      {
        "name": "get_weather",
        "description": "Get current weather for a city",
        "inputSchema": {
          "type": "object",
          "properties": {
            "city": {
              "type": "string",
              "description": "City name"
            }
          },
          "required": ["city"]
        },
        "logic": "city, ok := args[\"city\"].(string)\nif !ok {\n  response.Error = &ErrorObject{\n    Code: -32602,\n    Message: \"Invalid params: city must be a string\",\n  }\n  break\n}\n\n// In a real implementation, this would call a weather API\nresponse.Result = map[string]interface{}{\n  \"content\": []TextContent{\n    {\n      Type: \"text\",\n      Text: fmt.Sprintf(\"Weather for %s: Sunny, 22°C\", city),\n    },\n  },\n}"
      }
    ],
    "resources": [
      {
        "uri": "weather://cities/list",
        "name": "List of available cities",
        "mimeType": "application/json",
        "description": "List of cities with available weather data",
        "content": "[\"New York\", \"London\", \"Tokyo\", \"Sydney\", \"Paris\"]"
      }
    ]
  },
  "outputDir": "/path/to/output/directory"
}
</arguments>
</use_mcp_tool>
```

This will generate a fully functional weather MCP server in the specified output directory.
