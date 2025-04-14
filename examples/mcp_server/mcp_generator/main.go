// Package main provides an MCP server generator that allows AI models to create their own MCP servers
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
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

// ServerSpec represents the specification for an MCP server
type ServerSpec struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Tools       []ToolSpec     `json:"tools"`
	Resources   []ResourceSpec `json:"resources"`
}

// ToolSpec represents the specification for an MCP tool
type ToolSpec struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
	Logic       string                 `json:"logic"`
}

// ResourceSpec represents the specification for an MCP resource
type ResourceSpec struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	MimeType    string `json:"mimeType"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

// ServerTemplate is the Go template for generating an MCP server
const ServerTemplate = `// Package main provides an MCP server for {{.Name}}
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Request represents an MCP request
type Request struct {
	JSONRPC string                 ` + "`json:\"jsonrpc\"`" + `
	Method  string                 ` + "`json:\"method\"`" + `
	Params  map[string]interface{} ` + "`json:\"params\"`" + `
	ID      int                    ` + "`json:\"id\"`" + `
}

// Response represents an MCP response
type Response struct {
	JSONRPC string                 ` + "`json:\"jsonrpc\"`" + `
	Result  map[string]interface{} ` + "`json:\"result,omitempty\"`" + `
	Error   *ErrorObject           ` + "`json:\"error,omitempty\"`" + `
	ID      int                    ` + "`json:\"id\"`" + `
}

// ErrorObject represents an MCP error
type ErrorObject struct {
	Code    int    ` + "`json:\"code\"`" + `
	Message string ` + "`json:\"message\"`" + `
}

// TextContent represents text content in an MCP response
type TextContent struct {
	Type string ` + "`json:\"type\"`" + `
	Text string ` + "`json:\"text\"`" + `
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Log to stderr for debugging
	logf := func(format string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}

	logf("{{.Name}} MCP server starting")

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
					{{range .Tools}}
					{
						"name":        "{{.Name}}",
						"description": "{{.Description}}",
						"inputSchema": {{marshal .InputSchema}},
					},
					{{end}}
				},
			}

		case "listResources":
			response.Result = map[string]interface{}{
				"resources": []map[string]interface{}{
					{{range .Resources}}
					{
						"uri":         "{{.URI}}",
						"name":        "{{.Name}}",
						"mimeType":    "{{.MimeType}}",
						"description": "{{.Description}}",
					},
					{{end}}
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
			{{range .Resources}}
			if uri == "{{.URI}}" {
				response.Result = map[string]interface{}{
					"contents": []map[string]interface{}{
						{
							"uri":      uri,
							"mimeType": "{{.MimeType}}",
							"text":     {{printf "%q" .Content}},
						},
					},
				}
				found = true
				break
			}
			{{end}}

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
			{{range .Tools}}
			case "{{.Name}}":
				// Tool implementation
				{{.Logic}}
			{{end}}
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
`

// ReadmeTemplate is the template for generating a README.md file for the MCP server
const ReadmeTemplate = `# {{.Name}} MCP Server

{{.Description}}

## Features

{{range .Tools}}
- {{.Description}}
{{end}}

## Building

To build the {{.Name}} MCP server:

` + "```bash" + `
cd path/to/{{.Name}}
go build -o {{.Name}} {{.Name}}.go
` + "```" + `

## Usage

### Adding to Local_AI_Agents

To add the {{.Name}} MCP server to Local_AI_Agents, use the ` + "`/mcp-add`" + ` command:

` + "```" + `
/mcp-add {{.Name}} go run path/to/{{.Name}}/{{.Name}}.go
` + "```" + `

Or if you've built the binary:

` + "```" + `
/mcp-add {{.Name}} path/to/{{.Name}}/{{.Name}}
` + "```" + `

### Starting the Server

To start the {{.Name}} MCP server:

` + "```" + `
/mcp-start {{.Name}}
` + "```" + `

### Listing Available Tools

To list the tools provided by the {{.Name}} MCP server:

` + "```" + `
/mcp-tools {{.Name}}
` + "```" + `

### Available Tools

{{range .Tools}}
- ` + "`{{.Name}}`" + `: {{.Description}}
{{end}}

### Available Resources

{{range .Resources}}
- ` + "`{{.URI}}`" + `: {{.Description}}
{{end}}
`

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Log to stderr for debugging
	logf := func(format string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}

	logf("MCP Generator server starting")

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
						"name":        "generate_mcp_server",
						"description": "Generate a new MCP server based on a specification",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"spec": map[string]interface{}{
									"type":        "object",
									"description": "Server specification",
								},
								"outputDir": map[string]interface{}{
									"type":        "string",
									"description": "Output directory for the generated server",
								},
							},
							"required": []string{"spec", "outputDir"},
						},
					},
					{
						"name":        "validate_mcp_server",
						"description": "Validate an MCP server specification",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"spec": map[string]interface{}{
									"type":        "object",
									"description": "Server specification to validate",
								},
							},
							"required": []string{"spec"},
						},
					},
					{
						"name":        "build_mcp_server",
						"description": "Build an MCP server from source",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"serverPath": map[string]interface{}{
									"type":        "string",
									"description": "Path to the server source directory",
								},
							},
							"required": []string{"serverPath"},
						},
					},
				},
			}

		case "listResources":
			response.Result = map[string]interface{}{
				"resources": []map[string]interface{}{
					{
						"uri":         "mcp-generator://templates/server",
						"name":        "MCP Server Template",
						"mimeType":    "text/plain",
						"description": "Template for generating MCP servers",
					},
					{
						"uri":         "mcp-generator://templates/readme",
						"name":        "MCP Server README Template",
						"mimeType":    "text/plain",
						"description": "Template for generating MCP server README files",
					},
					{
						"uri":         "mcp-generator://examples/calculator",
						"name":        "Calculator MCP Server Example",
						"mimeType":    "application/json",
						"description": "Example specification for a calculator MCP server",
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

			var content string
			var mimeType string

			switch uri {
			case "mcp-generator://templates/server":
				content = ServerTemplate
				mimeType = "text/plain"
			case "mcp-generator://templates/readme":
				content = ReadmeTemplate
				mimeType = "text/plain"
			case "mcp-generator://examples/calculator":
				calculatorSpec := ServerSpec{
					Name:        "calculator",
					Description: "A simple calculator MCP server",
					Tools: []ToolSpec{
						{
							Name:        "add",
							Description: "Add two numbers",
							InputSchema: map[string]interface{}{
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
							Logic: `
								a, aOk := getFloat(args, "a")
								b, bOk := getFloat(args, "b")
								
								if !aOk || !bOk {
									response.Error = &ErrorObject{
										Code:    -32602,
										Message: "Invalid params: a and b must be numbers",
									}
									break
								}
								
								result := a + b
								resultStr := fmt.Sprintf("%.6f", result)
								
								response.Result = map[string]interface{}{
									"content": []TextContent{
										{
											Type: "text",
											Text: resultStr,
										},
									},
								}
							`,
						},
					},
					Resources: []ResourceSpec{
						{
							URI:         "calculator://constants/pi",
							Name:        "Pi constant",
							MimeType:    "text/plain",
							Description: "The mathematical constant π (pi)",
							Content:     "3.14159265358979323846",
						},
					},
				}
				specJSON, _ := json.MarshalIndent(calculatorSpec, "", "  ")
				content = string(specJSON)
				mimeType = "application/json"
			default:
				response.Error = &ErrorObject{
					Code:    -32602,
					Message: fmt.Sprintf("Resource not found: %s", uri),
				}
				break
			}

			if content != "" {
				response.Result = map[string]interface{}{
					"contents": []map[string]interface{}{
						{
							"uri":      uri,
							"mimeType": mimeType,
							"text":     content,
						},
					},
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
			case "generate_mcp_server":
				// Get the server specification
				specObj, ok := args["spec"].(map[string]interface{})
				if !ok {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: "Invalid params: spec must be an object",
					}
					break
				}

				// Convert the spec to a ServerSpec struct
				specJSON, err := json.Marshal(specObj)
				if err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to marshal spec: %v", err),
					}
					break
				}

				var spec ServerSpec
				if err := json.Unmarshal(specJSON, &spec); err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to unmarshal spec: %v", err),
					}
					break
				}

				// Get the output directory
				outputDir, ok := args["outputDir"].(string)
				if !ok {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: "Invalid params: outputDir must be a string",
					}
					break
				}

				// Create the output directory if it doesn't exist
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to create output directory: %v", err),
					}
					break
				}

				// Create a template function map
				funcMap := template.FuncMap{
					"marshal": func(v interface{}) string {
						a, _ := json.Marshal(v)
						return string(a)
					},
				}

				// Parse and execute the server template
				tmpl, err := template.New("server").Funcs(funcMap).Parse(ServerTemplate)
				if err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to parse server template: %v", err),
					}
					break
				}

				serverFile := filepath.Join(outputDir, spec.Name+".go")
				f, err := os.Create(serverFile)
				if err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to create server file: %v", err),
					}
					break
				}
				defer f.Close()

				if err := tmpl.Execute(f, spec); err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to execute server template: %v", err),
					}
					break
				}

				// Parse and execute the README template
				readmeTmpl, err := template.New("readme").Parse(ReadmeTemplate)
				if err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to parse README template: %v", err),
					}
					break
				}

				readmeFile := filepath.Join(outputDir, "README.md")
				rf, err := os.Create(readmeFile)
				if err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to create README file: %v", err),
					}
					break
				}
				defer rf.Close()

				if err := readmeTmpl.Execute(rf, spec); err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to execute README template: %v", err),
					}
					break
				}

				response.Result = map[string]interface{}{
					"content": []TextContent{
						{
							Type: "text",
							Text: fmt.Sprintf("Successfully generated MCP server in %s:\n- %s\n- %s", outputDir, serverFile, readmeFile),
						},
					},
				}

			case "validate_mcp_server":
				// Get the server specification
				specObj, ok := args["spec"].(map[string]interface{})
				if !ok {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: "Invalid params: spec must be an object",
					}
					break
				}

				// Convert the spec to a ServerSpec struct
				specJSON, err := json.Marshal(specObj)
				if err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to marshal spec: %v", err),
					}
					break
				}

				var spec ServerSpec
				if err := json.Unmarshal(specJSON, &spec); err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to unmarshal spec: %v", err),
					}
					break
				}

				// Validate the spec
				var validationErrors []string

				// Check for required fields
				if spec.Name == "" {
					validationErrors = append(validationErrors, "Server name is required")
				}
				if spec.Description == "" {
					validationErrors = append(validationErrors, "Server description is required")
				}
				if len(spec.Tools) == 0 && len(spec.Resources) == 0 {
					validationErrors = append(validationErrors, "Server must have at least one tool or resource")
				}

				// Check tools
				for i, tool := range spec.Tools {
					if tool.Name == "" {
						validationErrors = append(validationErrors, fmt.Sprintf("Tool %d: Name is required", i))
					}
					if tool.Description == "" {
						validationErrors = append(validationErrors, fmt.Sprintf("Tool %d: Description is required", i))
					}
					if tool.InputSchema == nil {
						validationErrors = append(validationErrors, fmt.Sprintf("Tool %d: InputSchema is required", i))
					}
					if tool.Logic == "" {
						validationErrors = append(validationErrors, fmt.Sprintf("Tool %d: Logic is required", i))
					}
				}

				// Check resources
				for i, resource := range spec.Resources {
					if resource.URI == "" {
						validationErrors = append(validationErrors, fmt.Sprintf("Resource %d: URI is required", i))
					}
					if resource.Name == "" {
						validationErrors = append(validationErrors, fmt.Sprintf("Resource %d: Name is required", i))
					}
					if resource.MimeType == "" {
						validationErrors = append(validationErrors, fmt.Sprintf("Resource %d: MimeType is required", i))
					}
					if resource.Description == "" {
						validationErrors = append(validationErrors, fmt.Sprintf("Resource %d: Description is required", i))
					}
				}

				if len(validationErrors) > 0 {
					response.Result = map[string]interface{}{
						"content": []TextContent{
							{
								Type: "text",
								Text: fmt.Sprintf("Validation failed:\n- %s", strings.Join(validationErrors, "\n- ")),
							},
						},
						"isError": true,
					}
				} else {
					response.Result = map[string]interface{}{
						"content": []TextContent{
							{
								Type: "text",
								Text: "Validation successful! The server specification is valid.",
							},
						},
					}
				}

			case "build_mcp_server":
				// Get the server path
				serverPath, ok := args["serverPath"].(string)
				if !ok {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: "Invalid params: serverPath must be a string",
					}
					break
				}

				// Check if the server path exists
				if _, err := os.Stat(serverPath); os.IsNotExist(err) {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: fmt.Sprintf("Server path does not exist: %s", serverPath),
					}
					break
				}

				// Get the server name from the directory name
				serverName := filepath.Base(serverPath)

				// Build the server
				cmd := exec.Command("go", "build", "-o", serverName, serverName+".go")
				cmd.Dir = serverPath
				output, err := cmd.CombinedOutput()
				if err != nil {
					response.Error = &ErrorObject{
						Code:    -32603,
						Message: fmt.Sprintf("Failed to build server: %v\n%s", err, string(output)),
					}
					break
				}

				response.Result = map[string]interface{}{
					"content": []TextContent{
						{
							Type: "text",
							Text: fmt.Sprintf("Successfully built MCP server: %s/%s", serverPath, serverName),
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
