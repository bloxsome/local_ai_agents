#!/bin/bash

# Script to run Ollama integration tests for Local AI Agents

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Local AI Agents - Ollama Integration Tests${NC}"
echo "========================================"

# Check if Ollama is installed
if ! command -v ollama &> /dev/null; then
    echo -e "${RED}Error: Ollama is not installed or not in PATH${NC}"
    echo "Please install Ollama from https://ollama.com/"
    exit 1
fi

# Check if Ollama is running
echo "Checking if Ollama is running..."
if ! curl -s http://localhost:11434 &> /dev/null; then
    echo -e "${YELLOW}Warning: Ollama is not running${NC}"
    echo "Starting Ollama..."
    
    # Try to start Ollama
    if command -v ollama &> /dev/null; then
        ollama serve &
        OLLAMA_PID=$!
        echo "Waiting for Ollama to start..."
        sleep 5
    else
        echo -e "${RED}Error: Could not start Ollama${NC}"
        echo "Please start Ollama manually and try again"
        exit 1
    fi
else
    echo -e "${GREEN}Ollama is running${NC}"
fi

# Check if the test model is available
echo "Checking if test model is available..."
MODEL="llama3.1:latest"
if ! curl -s "http://localhost:11434/api/tags" | grep -q "$MODEL"; then
    echo -e "${YELLOW}Warning: Test model $MODEL is not available${NC}"
    echo "Pulling model $MODEL..."
    ollama pull $MODEL
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}Error: Failed to pull model $MODEL${NC}"
        echo "Tests will use whatever model is available"
    else
        echo -e "${GREEN}Model $MODEL pulled successfully${NC}"
    fi
else
    echo -e "${GREEN}Model $MODEL is available${NC}"
fi

# Run the tests
echo "Running tests..."
cd ..
go test -v ./tests/...

# Check test results
if [ $? -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
else
    echo -e "${RED}Some tests failed${NC}"
fi

# Clean up
if [ ! -z "$OLLAMA_PID" ]; then
    echo "Stopping Ollama..."
    kill $OLLAMA_PID
fi

echo "Tests completed"
