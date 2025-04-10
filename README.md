# 🤖 Local_AI_Agents: Your Advanced AI Assistant Builder with MCP Integration 🚀

Welcome to Local_AI_Agents! This repository allows you to create sophisticated AI agents using local models, featuring powerful Model Context Protocol (MCP) integration for connecting to external tools and resources. It's like having a high-tech AI laboratory with seamless connectivity to the tools you need! 🧠✨

## 🌟 What's New?

- 🔌 Model Context Protocol (MCP): Seamlessly connect your AI to external tools and resources
- 🧩 MCP Server Integration: Create and manage custom MCP servers for extended functionality
- 🔄 Dynamic Tool Discovery: Automatically detect and use tools provided by MCP servers
- 🌐 Browser Integration: Control web browsers through MCP for web-based tasks
- 🧠 Enhanced Debug Agent with detailed cognitive processing visualization
- 🎭 Multi-agent system with easy switching between agents
- 🔀 Interactive follow-up question handling
- 🎨 Rich, colorful command-line interface with progress tracking
- 🛠️ Modular design with improved error handling and logging
- 🚀 High performance and better concurrency with Go
- 🔍 Improved memory search and context management

## 🚀 Key Features

1. 🔌 Model Context Protocol (MCP): Connect your AI to external tools and resources
2. 🧩 Extensible MCP Architecture: Create custom MCP servers for specialized functionality
3. 🌐 Browser Control: Automate web tasks through MCP browser integration
4. 🔍 API Integration: Connect to external APIs through custom MCP servers
5. 📚 Modular Architecture: Each function is in a separate module for easy customization
6. 💬 Interactive CLI: Clean, efficient command-line interface
7. 🔐 Secure Configuration: Customize your AI's personality and behavior through configuration
8. 🧪 Comprehensive Testing: Includes Ollama integration tests for quality assurance
9. 📜 Advanced Chat History: Built-in history management and analysis
10. 🎭 Multi-Agent System: Interact with multiple AI personalities in one session
11. 🤖 Debug Mode: Visualize the agent's thought process and decision-making in real-time
12. 🧐 Fact-Checking: Verify information and assess source credibility
13. 👤 User Profiling: Adapt responses based on user expertise and interests
14. 🚀 High Performance: Optimized Go implementation for speed and efficiency
15. 🛠️ Ollama Integration: Seamless connection to local Ollama models

## 💡 Why Model Context Protocol (MCP)?

Our implementation of the Model Context Protocol (MCP) offers powerful advantages for AI agents:

1. 🔌 Extensibility: Easily extend your AI's capabilities by connecting to external tools and services
2. 🧩 Modularity: Create specialized MCP servers for different tasks and domains
3. 🌐 Web Integration: Control browsers to perform complex web-based tasks
4. 🔄 Dynamic Discovery: Automatically detect and use tools provided by MCP servers
5. 🔒 Security: Fine-grained control over which tools require explicit approval
6. 📦 Portability: MCP servers can be shared and reused across different projects
7. 🛠️ Custom Tools: Create custom tools tailored to your specific needs

This approach allows Local AI Agents to interact with the world beyond just text, enabling more powerful and practical applications.

## 🛠️ Getting Started

### Prerequisites

- Go 1.18 or higher
- Git
- Local AI model (e.g., Ollama)

### Setting Up the Environment

1. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/Local_AI_Agents.git
   cd Local_AI_Agents
   ```

### Installing Ollama

1. Visit the [Ollama website](https://ollama.com/) and follow the installation instructions for your operating system.

2. Once installed, run Ollama and download a model (e.g., llama3.3:latest):
   ```bash
   ollama run llama3.3:latest
   ```

### Configuration

Configuration is managed through environment variables and defaults in the `config.go` file.

### Running the Application

```bash
go build
./local_ai_agents
```

See the [Setup Guide](docs/setup_guide.md) for detailed instructions.

### Testing with Ollama

The repository includes integration tests for Ollama:

```bash
cd tests
./run_tests.sh  # On Unix-like systems
run_tests.bat   # On Windows
```

These tests verify that Local AI Agents can connect to Ollama, send prompts, and handle responses correctly. See [tests/README.md](tests/README.md) for more details.

## 📘 Documentation

Detailed documentation can be found in the `docs/` directory:

### Getting Started
- [Setup Guide](docs/setup_guide.md)
- [Architecture Guide](docs/architecture_guide.md)
- [MCP Module](docs/mcp_module.md)

### System Architecture and Configuration
- [Architecture Guide](docs/architecture_guide.md)
- [Command Modules](docs/command_modules.md)
- [Config File](docs/config_file.md)

### Memory and History
- [JSON Memory Guide](docs/json_memory_guide.md)
- [Save History](docs/save_history.md)

### Utility and Logging
- [Logging Guide](docs/logging_guide.md)
- [Error Handling](docs/error_handling.md)

### User Guides
- [Building Agents](docs/building_agents.md)
- [Assistant User Guide](docs/assistant_user_guide.md)

## 🧠 Debug Agent Commands

- `/help`: Show available commands
- `/search <query>`: Perform an interactive web search
- `/context`: Show current context
- `/clear_context`: Clear the current context and bullet points
- `/bullets`: Display current bullet points
- `/knowledge_tree`: Display the knowledge tree
- `/explain <concept>`: Get an explanation of a concept
- `/fact_check <statement>`: Perform a fact check on a statement
- `/profile`: Display your user profile

## 🤝 Contributing

Got ideas? We love them! 💡 Submit a pull request or open an issue. Let's build the future of AI together!

## 📜 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

Built with ❤️ and 🧠 by the Local_AI_Agents team. Let's revolutionize AI connectivity and capabilities! 🚀
