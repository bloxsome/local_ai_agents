// Package main provides a browser MCP server example
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Request represents an MCP request
type Request struct {
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      int                    `json:"id"`
}

// Response represents an MCP response
type Response struct {
	JSONRPC string                 `json:"jsonrpc"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Error   *ErrorObject           `json:"error,omitempty"`
	ID      int                    `json:"id"`
}

// ErrorObject represents an MCP error
type ErrorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// TextContent represents text content in an MCP response
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Log to stderr for debugging
	logf := func(format string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}

	logf("Browser MCP server starting")

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		logf("Received request: %s", line)

		var request Request
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			logf("Failed to parse request: %v", err)
			continue
		}

		var response Response
		response.JSONRPC = "2.0"
		response.ID = request.ID

		switch request.Method {
		case "listTools":
			response.Result = map[string]interface{}{
				"tools": []map[string]interface{}{
					{
						"name":        "open_url",
						"description": "Open a URL in the default browser",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"url": map[string]interface{}{
									"type":        "string",
									"description": "URL to open",
								},
							},
							"required": []string{"url"},
						},
					},
				},
			}

		case "listResources":
			response.Result = map[string]interface{}{
				"resources": []map[string]interface{}{},
			}

		case "listResourceTemplates":
			response.Result = map[string]interface{}{
				"resourceTemplates": []map[string]interface{}{},
			}

		case "readResource":
			uri, ok := request.Params["uri"].(string)
			if !ok {
				response.Error = &ErrorObject{
					Code:    -32602,
					Message: "Invalid params: uri is required",
				}
				break
			}

			response.Error = &ErrorObject{
				Code:    -32602,
				Message: fmt.Sprintf("Resource not found: %s", uri),
			}

		case "callTool":
			name, ok := request.Params["name"].(string)
			if !ok {
				response.Error = &ErrorObject{
					Code:    -32602,
					Message: "Invalid params: name is required",
				}
				break
			}

			args, ok := request.Params["arguments"].(map[string]interface{})
			if !ok {
				response.Error = &ErrorObject{
					Code:    -32602,
					Message: "Invalid params: arguments is required",
				}
				break
			}

			switch name {
			case "open_url":
				url, ok := args["url"].(string)
				if !ok {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: "Invalid params: url is required",
					}
					break
				}

				err := openURL(url)
				if err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to open URL: %v", err),
					}
					break
				}

				response.Result = map[string]interface{}{
					"content": []TextContent{
						{
							Type: "text",
							Text: fmt.Sprintf("Opened URL: %s", url),
						},
					},
				}

			default:
				response.Error = &ErrorObject{
					Code:    -32601,
					Message: fmt.Sprintf("Method not found: %s", name),
				}
			}

		default:
			response.Error = &ErrorObject{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", request.Method),
			}
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			logf("Failed to marshal response: %v", err)
			continue
		}

		fmt.Println(string(responseJSON))
		logf("Sent response: %s", string(responseJSON))
	}

	if err := scanner.Err(); err != nil {
		logf("Scanner error: %v", err)
		os.Exit(1)
	}
}

// openURL opens a URL in the default browser
func openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}
