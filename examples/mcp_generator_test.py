#!/usr/bin/env python3
"""
Test script to demonstrate how llama3.3 can use the MCP Generator server
to create its own MCP servers.

This script:
1. Starts the MCP Generator server
2. Simulates a conversation with llama3.3 where it uses the MCP Generator
3. Validates the generated MCP server

Requirements:
- Python 3.6+
- requests library (pip install requests)
- subprocess library (standard)
- json library (standard)
- os library (standard)
- time library (standard)
"""

import subprocess
import requests
import json
import os
import time
import argparse
import sys

# Parse command line arguments
parser = argparse.ArgumentParser(description='Test the MCP Generator with llama3.3')
parser.add_argument('--ollama-url', default='http://localhost:11434/api/generate',
                    help='URL for Ollama API (default: http://localhost:11434/api/generate)')
parser.add_argument('--model', default='llama3.3:latest',
                    help='Ollama model to use (default: llama3.3:latest)')
parser.add_argument('--output-dir', default='./generated_mcp_server',
                    help='Directory to output the generated MCP server (default: ./generated_mcp_server)')
args = parser.parse_args()

# Ensure output directory exists
os.makedirs(args.output_dir, exist_ok=True)

# Path to the MCP Generator binary
MCP_GENERATOR_PATH = os.path.join(os.path.dirname(os.path.abspath(__file__)), 
                                 "../examples/mcp_server/mcp_generator/mcp_generator")

def start_mcp_generator():
    """Start the MCP Generator server as a subprocess"""
    print("Starting MCP Generator server...")
    process = subprocess.Popen(
        [MCP_GENERATOR_PATH],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        bufsize=1
    )
    # Give it a moment to start up
    time.sleep(1)
    return process

def stop_mcp_generator(process):
    """Stop the MCP Generator server"""
    print("Stopping MCP Generator server...")
    process.terminate()
    process.wait()

def send_mcp_request(process, method, params=None):
    """Send a request to the MCP Generator server and get the response"""
    if params is None:
        params = {}
    
    request = {
        "jsonrpc": "2.0",
        "method": method,
        "params": params,
        "id": 1
    }
    
    request_json = json.dumps(request)
    print(f"Sending request: {request_json}")
    
    # Write the request to the process's stdin
    process.stdin.write(request_json + "\n")
    process.stdin.flush()
    
    # Read the response from the process's stdout
    response_json = process.stdout.readline().strip()
    print(f"Received response: {response_json}")
    
    return json.loads(response_json)

def query_ollama(prompt):
    """Query the Ollama API with a prompt"""
    print(f"Querying Ollama with prompt: {prompt[:100]}...")
    
    response = requests.post(
        args.ollama_url,
        json={
            "model": args.model,
            "prompt": prompt,
            "stream": False
        }
    )
    
    if response.status_code != 200:
        print(f"Error querying Ollama: {response.status_code} {response.text}")
        return None
    
    return response.json()["response"]

def main():
    # Start the MCP Generator server
    mcp_generator = start_mcp_generator()
    
    try:
        # List the available tools
        tools_response = send_mcp_request(mcp_generator, "listTools")
        if "error" in tools_response:
            print(f"Error listing tools: {tools_response['error']}")
            return 1
        
        # Get the calculator example
        example_response = send_mcp_request(mcp_generator, "readResource", {"uri": "mcp-generator://examples/calculator"})
        if "error" in example_response:
            print(f"Error getting calculator example: {example_response['error']}")
            return 1
        
        calculator_spec = example_response["result"]["contents"][0]["text"]
        
        # Create a prompt for llama3.3 to modify the calculator example
        prompt = f"""
You are an AI assistant that can create MCP (Model Context Protocol) servers. 
I want you to create a simple "greeting" MCP server that provides a tool to greet users.

Here's an example of a calculator MCP server specification:

```json
{calculator_spec}
```

Please create a similar specification for a greeting server that:
1. Has a name of "greeting"
2. Has a description of "A simple greeting MCP server"
3. Provides a tool called "greet" that takes a "name" parameter and returns a greeting message
4. Has a resource that provides a list of greeting templates

Format your response as a valid JSON object that follows the same structure as the calculator example.
Only include the JSON, no other text.
"""
        
        # Query llama3.3 to get the greeting server specification
        llama_response = query_ollama(prompt)
        if not llama_response:
            print("Failed to get response from Ollama")
            return 1
        
        # Extract the JSON from the response
        try:
            # Try to find JSON in the response
            json_start = llama_response.find('{')
            json_end = llama_response.rfind('}') + 1
            if json_start == -1 or json_end == 0:
                print("Could not find JSON in the response")
                print(f"Response: {llama_response}")
                return 1
            
            greeting_spec_json = llama_response[json_start:json_end]
            greeting_spec = json.loads(greeting_spec_json)
            
            print("\nGreeting Server Specification generated by llama3.3:")
            print(json.dumps(greeting_spec, indent=2))
            
            # Validate the specification
            validate_response = send_mcp_request(mcp_generator, "callTool", {
                "name": "validate_mcp_server",
                "arguments": {
                    "spec": greeting_spec
                }
            })
            
            if "error" in validate_response:
                print(f"Error validating specification: {validate_response['error']}")
                return 1
            
            validation_result = validate_response["result"]["content"][0]["text"]
            print(f"\nValidation result: {validation_result}")
            
            if "Validation successful" not in validation_result:
                print("Validation failed, fixing specification...")
                # Here you could add code to fix the specification if needed
            
            # Generate the MCP server
            generate_response = send_mcp_request(mcp_generator, "callTool", {
                "name": "generate_mcp_server",
                "arguments": {
                    "spec": greeting_spec,
                    "outputDir": args.output_dir
                }
            })
            
            if "error" in generate_response:
                print(f"Error generating MCP server: {generate_response['error']}")
                return 1
            
            generation_result = generate_response["result"]["content"][0]["text"]
            print(f"\nGeneration result: {generation_result}")
            
            # Build the generated MCP server
            build_response = send_mcp_request(mcp_generator, "callTool", {
                "name": "build_mcp_server",
                "arguments": {
                    "serverPath": args.output_dir
                }
            })
            
            if "error" in build_response:
                print(f"Error building MCP server: {build_response['error']}")
                return 1
            
            build_result = build_response["result"]["content"][0]["text"]
            print(f"\nBuild result: {build_result}")
            
            print("\nTest completed successfully!")
            print(f"Generated MCP server is available in: {args.output_dir}")
            
        except json.JSONDecodeError as e:
            print(f"Error parsing JSON from llama3.3 response: {e}")
            print(f"Response: {llama_response}")
            return 1
        
    finally:
        # Stop the MCP Generator server
        stop_mcp_generator(mcp_generator)
    
    return 0

if __name__ == "__main__":
    sys.exit(main())
