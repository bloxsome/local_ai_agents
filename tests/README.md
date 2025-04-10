# Local AI Agents - Ollama Integration Tests

This directory contains tests for the integration between Local AI Agents and Ollama.

## Prerequisites

- Go 1.18 or higher
- [Ollama](https://ollama.com/) installed and available in your PATH
- An LLM model available in Ollama (default: `llama3.1:latest`)

## Running the Tests

### On Unix-like Systems (Linux, macOS)

1. Make the script executable:
   ```bash
   chmod +x run_tests.sh
   ```

2. Run the tests:
   ```bash
   ./run_tests.sh
   ```

### On Windows

1. Run the tests:
   ```cmd
   run_tests.bat
   ```

## What the Tests Do

The test suite includes:

1. **Basic Connectivity Test** (`TestOllamaConnection`): Verifies that Local AI Agents can connect to a running Ollama instance.

2. **Prompt Processing Test** (`TestOllamaPrompt`): Tests sending a prompt to Ollama and receiving a response.

3. **Error Handling Test** (`TestOllamaErrorHandling`): Tests how the system handles errors when Ollama is not available.

4. **Configuration Integration Test** (`TestOllamaWithConfig`): Tests using Ollama with the Local AI Agents configuration.

## Test Scripts

The test scripts (`run_tests.sh` and `run_tests.bat`) perform the following actions:

1. Check if Ollama is installed
2. Check if Ollama is running (and try to start it if not)
3. Check if the test model is available (and try to pull it if not)
4. Run the tests
5. Report the results

## Manual Testing

You can also run the tests manually using Go's testing framework:

```bash
cd /path/to/local_ai_agents
go test -v ./tests/...
```

## Troubleshooting

If the tests fail, check the following:

1. Is Ollama installed and in your PATH?
2. Is Ollama running? You can start it with `ollama serve`
3. Is the test model available? You can check with `ollama list` and pull it with `ollama pull llama3.1:latest`
4. Are there any network issues preventing connection to Ollama?

## Adding More Tests

To add more tests:

1. Create a new test function in one of the existing test files or create a new test file
2. Follow the Go testing conventions (function name should start with `Test`)
3. Use the `TestSuite` structure in `integration_test.go` for integration tests
