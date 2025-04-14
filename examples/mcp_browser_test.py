#!/usr/bin/env python3
"""
Test script for Ollama MCP browser support.
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
            print("\nTest completed successfully!")
        else:
            print("\nTest failed: No response from Ollama.")
    
    finally:
        # Clean up MCP settings
        cleanup_mcp_settings(settings_path)

if __name__ == "__main__":
    main()
