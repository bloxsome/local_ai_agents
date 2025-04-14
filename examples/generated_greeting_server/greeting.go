// Package main provides an MCP server for greeting
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
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

// Helper functions
func getString(args map[string]interface{}, key string) (string, bool) {
	v, ok := args[key]
	if !ok {
		return "", false
	}

	s, ok := v.(string)
	return s, ok
}

func getResource(uri string) (string, bool) {
	// In a real implementation, this would be more sophisticated
	// For now, we'll just hardcode the greeting templates
	if uri == "greeting://templates/list" {
		template := "Hello, %s!\nHi %s!\nGreetings %s!"
		return template, true
	}
	return "", false
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Log to stderr for debugging
	logf := func(format string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}

	logf("greeting MCP server starting")

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
						"name":        "greet",
						"description": "Generate a personalized greeting message",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name": map[string]interface{}{
									"type":        "string",
									"description": "The name to include in the greeting",
								},
							},
							"required": []string{"name"},
						},
					},
				},
			}

		case "listResources":
			response.Result = map[string]interface{}{
				"resources": []map[string]interface{}{
					{
						"uri":         "greeting://templates/list",
						"name":        "Greeting templates",
						"mimeType":    "text/plain",
						"description": "A list of greeting message templates",
					},
				},
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

			var found bool

			if uri == "greeting://templates/list" {
				response.Result = map[string]interface{}{
					"contents": []map[string]interface{}{
						{
							"uri":      uri,
							"mimeType": "text/plain",
							"text":     "Hello, %s!\nHi %s!\nGreetings %s!",
						},
					},
				}
				found = true
				break
			}

			if !found {
				response.Error = &ErrorObject{
					Code:    -32602,
					Message: fmt.Sprintf("Resource not found: %s", uri),
				}
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
			case "greet":
				// Tool implementation
				name, nameOk := getString(args, "name")

				if !nameOk {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: "Invalid params: name must be a string",
					}
					break
				}

				greetingTemplate, greetingTemplateOk := getResource("greeting://templates/list")
				if !greetingTemplateOk {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: "Failed to retrieve greeting template",
					}
					break
				}

				greetingMessage := fmt.Sprintf(greetingTemplate, name, name, name)

				response.Result = map[string]interface{}{
					"content": []TextContent{
						{
							Type: "text",
							Text: greetingMessage,
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
