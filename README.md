# 🤖 Local_AI_Agents: Your Advanced AI Assistant Builder with Graph Knowledgebase 🚀

Welcome to Local_AI_Agents! This repository allows you to create sophisticated AI agents using local models, featuring a unique graph-based knowledgebase with a high-performance Go implementation. It's like having a high-tech AI laboratory with a built-in brain! 🧠✨

## 🌟 What's New?

- 🕸️ Graph-based Knowledgebase: A novel approach using JSON for flexible, relational knowledge storage
- 🧠 Enhanced Debug Agent with detailed cognitive processing visualization
- 🌳 Dynamic Knowledge Tree generation and management
- 🔍 Improved memory search and context management
- 🧐 Fact-checking and source credibility assessment
- 🎭 Multi-agent system with easy switching between agents
- 🔀 Interactive follow-up question handling
- 🎨 Rich, colorful command-line interface with progress tracking
- 🛠️ Modular design with improved error handling and logging
- 🚀 High performance and better concurrency with Go
- 🔌 Model Context Protocol (MCP) for external system integration

## 🚀 Key Features

1. 📊 Graph Knowledgebase: Utilizes a JSON-based graph structure for flexible and relational knowledge representation
2. 📚 Modular Architecture: Each function is in a separate module for easy customization and extension
3. 💬 Interactive CLI: Clean, efficient command-line interface
4. 🔐 Secure Configuration: Customize your AI's personality and behavior through configuration
5. 🧪 Comprehensive Testing: Because quality is our superpower!
6. 🌐 Web Search Integration: Your AI can search the web using DuckDuckGo
7. 📜 Advanced Chat History: Never forget a conversation with built-in history management and analysis
8. 🧠 Sophisticated Memory Search: Quickly retrieve and utilize relevant information from past interactions and uploaded documents
9. 🧵 Fabric Integration: Use Fabric patterns for enhanced AI interactions
10. 🎭 Multi-Agent System: Interact with multiple AI personalities in one session
11. 🤖 Debug Mode: Visualize the agent's thought process and decision-making in real-time
12. 🌳 Knowledge Tree: Dynamic generation and visualization of knowledge structures
13. 🧐 Fact-Checking: Verify information and assess source credibility
14. 👤 User Profiling: Adapt responses based on user expertise and interests
15. 🔌 MCP Support: Connect to external tools and resources through the Model Context Protocol

## 💡 Why JSON-based Graph Knowledgebase?

Our unique approach of using a JSON-based graph structure for the knowledgebase offers several advantages:

1. 🔄 Flexibility: Easily adapt and evolve the knowledge structure as your AI learns
2. 🔗 Rich Relationships: Capture complex relationships between concepts more intuitively than in traditional vector databases
3. 🚀 Performance: Efficient querying and updating of interconnected information
4. 🧩 Simplicity: No need for complex vector database setups or maintenance
5. 📦 Portability: JSON format allows for easy data transfer and backup
6. 🔍 Interpretability: Graph structure provides clear visibility into the AI's knowledge connections

This approach allows Local_AI_Agents to have a more nuanced and context-aware understanding, leading to more intelligent and adaptive responses.

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

2. Once installed, run Ollama and download a model (e.g., llama3.1:latest):
   ```bash
   ollama run llama3.1:latest
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

Built with ❤️ and 🧠 by the Local_AI_Agents team. Let's revolutionize AI knowledge representation! 🚀
