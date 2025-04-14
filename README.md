# Local AI Agents

A Go-based framework for creating AI agents using local models with Model Context Protocol (MCP) integration for connecting to external tools and resources.

## Overview

Local AI Agents provides a foundation for building AI assistants that can:
- Connect to local language models (via Ollama)
- Integrate with external tools through the Model Context Protocol (MCP)
- Allow the AI model to build tools dynamically 
- Control web browsers for automated tasks*
- Maintain conversation history and context
- Support multiple agent personalities

## Prerequisites

- Go 1.18 or higher
- Git
- Local AI model (e.g., Ollama)

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/bloxsome/local_ai_agents.git
   cd Local_AI_Agents
   ```

2. Install Ollama:
   - Visit [Ollama website](https://ollama.com/) and follow installation instructions
   - Download a model:
     ```bash
     ollama run llama3.3:latest
     ```

3. Build the application:
   ```bash
   go build
   ```

## Usage

Run the application:
```bash
./local_ai_agents
```

### Configuration

Configuration is managed through environment variables and defaults in the `config.go` file.

### Available Commands

- `/help`: Show available commands
- `/search <query>`: Perform an interactive web search
- `/context`: Show current context
- `/clear_context`: Clear the current context
- `/bullets`: Display current bullet points
- `/knowledge_tree`: Display the knowledge tree
- `/explain <concept>`: Get an explanation of a concept
- `/fact_check <statement>`: Perform a fact check on a statement
- `/profile`: Display your user profile

## Testing

Run the integration tests:
```bash
cd tests
./run_tests.sh  # On Unix-like systems
run_tests.bat   # On Windows - Not Tested
```

See [tests/README.md](tests/README.md) for more details on testing.

## Extending the Framework

### AI-Assisted MCP Server Development

Yes, AI models (like Claude) can write MCP servers for you let try and bring this local with ollama. The MCP protocol follows a well-defined structure that makes it suitable for AI-assisted development:

1. Start with an existing example (calculator.go or weather.go) as a template
2. Define the tools and resources your server will provide
3. Implement the required MCP methods (listTools, listResources, etc.)
4. Add your business logic for handling tool calls

This approach allows you to rapidly prototype new MCP servers that extend your AI agent's capabilities with custom functionality or API integrations.

#### MCP Generator Server

The repository includes an MCP Generator server that allows AI models to create their own MCP servers:

- Located in `examples/mcp_server/mcp_generator/`
- Provides tools for generating, validating, and building MCP servers
- Allows AI models to create custom MCP servers based on specifications
- Includes templates and examples for rapid development

See the [MCP Generator README](examples/mcp_server/mcp_generator/README.md) for details on how to use this server.

#### Testing with llama3.3

You can test if llama3.3 is able to use the MCP Generator with the provided test script:

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

### Model Context Protocol (MCP)

The MCP integration allows you to extend your AI agent's capabilities by connecting to external tools and resources:

1. Create custom MCP servers for specialized functionality:
   - See fully functioning example in `examples/mcp_server/` (calculator server)
   - Reference implementation in `internal/modules/mcp/`

2. Implement new MCP tools:
   - Browser control: `examples/mcp_browser_server/`
   - API integrations: `tests/mcp_server/weather.go` (complete weather server example)

3. Configure MCP settings:
   - Example configuration: `examples/mcp_settings.json`
   - Initialization example: `examples/mcp_initialization_example.go`

### Adding New Modules

The framework uses a modular architecture:

1. Create a new module in `internal/modules/`
2. Implement the required interfaces
3. Register your module in `main.go`

See existing modules for reference:
- Command handling: `internal/modules/commands/`
- Input processing: `internal/modules/input/`
- Logging: `internal/modules/logging/`

### Customizing Agent Behavior

Modify agent behavior by:
1. Updating prompt templates in the configuration
2. Adding new command handlers
3. Implementing custom memory/context management

## Documentation

Additional documentation can be found in the `docs/` directory:
- [MCP Module](docs/mcp_module.md)

## Contributing

Contributions are welcome. Please submit pull requests or open issues to improve the project.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
