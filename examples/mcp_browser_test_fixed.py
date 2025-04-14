#!/usr/bin/env python3
"""
Test script for Ollama MCP browser support with explicit MCP initialization.
This script demonstrates how to use the MCP browser server with Ollama.
"""

import os
import subprocess
import time
import json
import requests
import sys
import shutil
from pathlib import Path

# Configuration
OLLAMA_API_URL = "http://localhost:11434/api/generate"
MODEL = "llama3.3"  # Change to your preferred model
MCP_SETTINGS_PATH = os.path.join(os.getcwd(), "ollama_agents/examples/mcp_settings.json")
TEST_PAGE_PATH = os.path.join(os.getcwd(), "ollama_agents/examples/test_page.html")
BROWSER_SERVER_PATH = os.path.join(os.getcwd(), "ollama_agents/examples/mcp_browser_server/browser")

def setup_mcp_settings():
    """Copy MCP settings to the correct location."""
    print("Setting up MCP settings...")
    
    # Get user home directory
    home_dir = str(Path.home())
    
    # Create config directory if it doesn't exist
    config_dir = os.path.join(home_dir, ".config", "ollama_agents")
    os.makedirs(config_dir, exist_ok=True)
    
    # Copy settings file
    settings_dest = os.path.join(config_dir, "mcp_settings.json")
    shutil.copy(MCP_SETTINGS_PATH, settings_dest)
    print(f"Copied MCP settings to {settings_dest}")
    
    return settings_dest

def cleanup_mcp_settings(settings_path):
    """Remove the MCP settings file."""
    print(f"Cleaning up MCP settings at {settings_path}...")
    if os.path.exists(settings_path):
        os.remove(settings_path)

def start_browser_server():
    """Start the MCP browser server."""
    print("Starting MCP browser server...")
    process = subprocess.Popen(
        [BROWSER_SERVER_PATH],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        stdin=subprocess.PIPE,
        text=True
    )
    time.sleep(1)  # Give the server a moment to start
    return process

def stop_browser_server(process):
    """Stop the MCP browser server."""
    print("Stopping MCP browser server...")
    if process:
        process.terminate()
        try:
            process.wait(timeout=5)
        except subprocess.TimeoutExpired:
            process.kill()

def send_prompt_to_ollama(prompt):
    """Send a prompt to Ollama and return the response."""
    print(f"Sending prompt to Ollama: {prompt}")
    
    data = {
        "model": MODEL,
        "prompt": prompt,
        "stream": False
    }
    
    try:
        response = requests.post(OLLAMA_API_URL, json=data)
        response.raise_for_status()
        return response.json()["response"]
    except requests.exceptions.RequestException as e:
        print(f"Error sending prompt to Ollama: {e}")
        return None

def main():
    """Main function to test Ollama MCP browser support."""
    # Setup MCP settings
    settings_path = setup_mcp_settings()
    
    # Start the MCP browser server
    browser_server = start_browser_server()
    
    try:
        # Create a prompt that uses the MCP browser tool
        test_page_url = f"file://{TEST_PAGE_PATH}"
        prompt = f"""
        Open the test page in my browser.
        
        <use_mcp_tool>
        <server_name>browser</server_name>
        <tool_name>open_url</tool_name>
        <arguments>
        {{
          "url": "{test_page_url}"
        }}
        </arguments>
        </use_mcp_tool>
        """
        
        # Send the prompt to Ollama
        response = send_prompt_to_ollama(prompt)
        
        if response:
            print("\nOllama response:")
            print(response)
            
            # Check if the response contains the expected result
            if "Error: MCP module not initialized" in response:
                print("\nTest failed: MCP module not initialized.")
                print("This suggests that the Ollama client is not properly initializing the MCP module.")
                print("You may need to modify the Ollama client code to ensure the MCP module is initialized.")
            elif "Opened URL" in response:
                print("\nTest succeeded! The MCP browser server opened the URL.")
            else:
                print("\nTest completed, but the response doesn't indicate success or failure.")
                print("The MCP command might not have been processed as expected.")
        else:
            print("\nTest failed: No response from Ollama.")
    
    finally:
        # Stop the MCP browser server
        stop_browser_server(browser_server)
        
        # Clean up MCP settings
        cleanup_mcp_settings(settings_path)

if __name__ == "__main__":
    main()
