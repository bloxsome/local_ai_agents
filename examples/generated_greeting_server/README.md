# greeting MCP Server

A simple greeting MCP server

## Features


- Generate a personalized greeting message


## Building

To build the greeting MCP server:

```bash
cd path/to/greeting
go build -o greeting greeting.go
```

## Usage

### Adding to Local_AI_Agents

To add the greeting MCP server to Local_AI_Agents, use the `/mcp-add` command:

```
/mcp-add greeting go run path/to/greeting/greeting.go
```

Or if you've built the binary:

```
/mcp-add greeting path/to/greeting/greeting
```

### Starting the Server

To start the greeting MCP server:

```
/mcp-start greeting
```

### Listing Available Tools

To list the tools provided by the greeting MCP server:

```
/mcp-tools greeting
```

### Available Tools


- `greet`: Generate a personalized greeting message


### Available Resources


- `greeting://templates/list`: A list of greeting message templates

