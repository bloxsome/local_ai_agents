// Package main provides a simple MCP server example
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

	logf("Calculator MCP server starting")

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
						"name":        "add",
						"description": "Add two numbers",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"a": map[string]interface{}{
									"type":        "number",
									"description": "First number",
								},
								"b": map[string]interface{}{
									"type":        "number",
									"description": "Second number",
								},
							},
							"required": []string{"a", "b"},
						},
					},
					{
						"name":        "subtract",
						"description": "Subtract two numbers",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"a": map[string]interface{}{
									"type":        "number",
									"description": "First number",
								},
								"b": map[string]interface{}{
									"type":        "number",
									"description": "Second number",
								},
							},
							"required": []string{"a", "b"},
						},
					},
					{
						"name":        "multiply",
						"description": "Multiply two numbers",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"a": map[string]interface{}{
									"type":        "number",
									"description": "First number",
								},
								"b": map[string]interface{}{
									"type":        "number",
									"description": "Second number",
								},
							},
							"required": []string{"a", "b"},
						},
					},
					{
						"name":        "divide",
						"description": "Divide two numbers",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"a": map[string]interface{}{
									"type":        "number",
									"description": "First number",
								},
								"b": map[string]interface{}{
									"type":        "number",
									"description": "Second number (non-zero)",
								},
							},
							"required": []string{"a", "b"},
						},
					},
				},
			}

		case "listResources":
			response.Result = map[string]interface{}{
				"resources": []map[string]interface{}{
					{
						"uri":         "calculator://constants/pi",
						"name":        "Pi constant",
						"mimeType":    "text/plain",
						"description": "The mathematical constant π (pi)",
					},
					{
						"uri":         "calculator://constants/e",
						"name":        "Euler's number",
						"mimeType":    "text/plain",
						"description": "The mathematical constant e",
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

			var text string
			switch uri {
			case "calculator://constants/pi":
				text = "3.14159265358979323846"
			case "calculator://constants/e":
				text = "2.71828182845904523536"
			default:
				response.Error = &ErrorObject{
					Code:    -32602,
					Message: fmt.Sprintf("Resource not found: %s", uri),
				}
				break
			}

			response.Result = map[string]interface{}{
				"contents": []map[string]interface{}{
					{
						"uri":      uri,
						"mimeType": "text/plain",
						"text":     text,
					},
				},
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

			a, aOk := getFloat(args, "a")
			b, bOk := getFloat(args, "b")

			if !aOk || !bOk {
				response.Error = &ErrorObject{
					Code:    -32602,
					Message: "Invalid params: a and b must be numbers",
				}
				break
			}

			var result float64
			var err error

			switch name {
			case "add":
				result = a + b
			case "subtract":
				result = a - b
			case "multiply":
				result = a * b
			case "divide":
				if b == 0 {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: "Division by zero",
					}
					break
				}
				result = a / b
			default:
				response.Error = &ErrorObject{
					Code:    -32601,
					Message: fmt.Sprintf("Method not found: %s", name),
				}
				break
			}

			if err != nil {
				response.Error = &ErrorObject{
					Code:    -32603,
					Message: err.Error(),
				}
				break
			}

			resultStr := fmt.Sprintf("%.6f", result)

			response.Result = map[string]interface{}{
				"content": []TextContent{
					{
						Type: "text",
						Text: resultStr,
					},
				},
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

// getFloat gets a float value from a map
func getFloat(m map[string]interface{}, key string) (float64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}

	switch v := v.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}
