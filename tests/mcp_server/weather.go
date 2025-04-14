// Package main provides a simple weather MCP server for testing
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
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

// WeatherData represents weather data for a city
type WeatherData struct {
	Temperature float64 `json:"temperature"`
	Conditions  string  `json:"conditions"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"windSpeed"`
}

// Cities is a map of city names to weather data
var Cities = map[string]WeatherData{
	"new york": {
		Temperature: 22.5,
		Conditions:  "Partly Cloudy",
		Humidity:    65,
		WindSpeed:   10.2,
	},
	"london": {
		Temperature: 18.0,
		Conditions:  "Rainy",
		Humidity:    80,
		WindSpeed:   15.5,
	},
	"tokyo": {
		Temperature: 28.3,
		Conditions:  "Sunny",
		Humidity:    45,
		WindSpeed:   8.7,
	},
	"sydney": {
		Temperature: 26.8,
		Conditions:  "Clear",
		Humidity:    50,
		WindSpeed:   12.3,
	},
	"paris": {
		Temperature: 20.1,
		Conditions:  "Cloudy",
		Humidity:    70,
		WindSpeed:   9.8,
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Log to stderr for debugging
	logf := func(format string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}

	logf("Weather MCP server starting")

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
						"name":        "get_weather",
						"description": "Get weather information for a city",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"city": map[string]interface{}{
									"type":        "string",
									"description": "City name",
								},
							},
							"required": []string{"city"},
						},
					},
					{
						"name":        "get_forecast",
						"description": "Get weather forecast for a city",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"city": map[string]interface{}{
									"type":        "string",
									"description": "City name",
								},
								"days": map[string]interface{}{
									"type":        "number",
									"description": "Number of days (1-5)",
								},
							},
							"required": []string{"city"},
						},
					},
				},
			}

		case "listResources":
			response.Result = map[string]interface{}{
				"resources": []map[string]interface{}{
					{
						"uri":         "weather://cities/list",
						"name":        "List of available cities",
						"mimeType":    "application/json",
						"description": "List of cities with available weather data",
					},
				},
			}

		case "listResourceTemplates":
			response.Result = map[string]interface{}{
				"resourceTemplates": []map[string]interface{}{
					{
						"uriTemplate": "weather://cities/{city}",
						"name":        "Weather for a specific city",
						"mimeType":    "application/json",
						"description": "Weather information for a specific city",
					},
				},
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

			if uri == "weather://cities/list" {
				// Return list of cities
				cityList := make([]string, 0, len(Cities))
				for city := range Cities {
					cityList = append(cityList, city)
				}

				cityListJSON, _ := json.Marshal(cityList)

				response.Result = map[string]interface{}{
					"contents": []map[string]interface{}{
						{
							"uri":      uri,
							"mimeType": "application/json",
							"text":     string(cityListJSON),
						},
					},
				}
			} else {
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
			case "get_weather":
				city, ok := args["city"].(string)
				if !ok {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: "Invalid params: city is required",
					}
					break
				}

				// Convert city to lowercase for case-insensitive lookup
				city = strings.ToLower(city)

				// Get weather data for city
				weatherData, ok := Cities[city]
				if !ok {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: fmt.Sprintf("City not found: %s", city),
					}
					break
				}

				// Format weather data as text
				weatherText := fmt.Sprintf("Weather for %s:\nTemperature: %.1f°C\nConditions: %s\nHumidity: %d%%\nWind Speed: %.1f km/h",
					city, weatherData.Temperature, weatherData.Conditions, weatherData.Humidity, weatherData.WindSpeed)

				response.Result = map[string]interface{}{
					"content": []TextContent{
						{
							Type: "text",
							Text: weatherText,
						},
					},
				}

			case "get_forecast":
				city, ok := args["city"].(string)
				if !ok {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: "Invalid params: city is required",
					}
					break
				}

				// Convert city to lowercase for case-insensitive lookup
				city = strings.ToLower(city)

				// Get weather data for city
				weatherData, ok := Cities[city]
				if !ok {
					response.Error = &ErrorObject{
						Code:    -32602,
						Message: fmt.Sprintf("City not found: %s", city),
					}
					break
				}

				// Get number of days (default to 3)
				days := 3
				if daysArg, ok := args["days"].(float64); ok {
					days = int(daysArg)
					if days < 1 {
						days = 1
					} else if days > 5 {
						days = 5
					}
				}

				// Generate forecast
				now := time.Now()
				var forecastText strings.Builder
				forecastText.WriteString(fmt.Sprintf("Weather forecast for %s:\n\n", city))

				for i := 0; i < days; i++ {
					date := now.AddDate(0, 0, i)
					// Slightly vary the weather data for each day
					tempVariation := (float64(i) * 0.5) - 1.0
					humidityVariation := i * 2
					windVariation := (float64(i) * 0.3) - 0.6

					forecastText.WriteString(fmt.Sprintf("Date: %s\n", date.Format("2006-01-02")))
					forecastText.WriteString(fmt.Sprintf("Temperature: %.1f°C\n", weatherData.Temperature+tempVariation))
					forecastText.WriteString(fmt.Sprintf("Conditions: %s\n", weatherData.Conditions))
					forecastText.WriteString(fmt.Sprintf("Humidity: %d%%\n", weatherData.Humidity+humidityVariation))
					forecastText.WriteString(fmt.Sprintf("Wind Speed: %.1f km/h\n", weatherData.WindSpeed+windVariation))
					forecastText.WriteString("\n")
				}

				response.Result = map[string]interface{}{
					"content": []TextContent{
						{
							Type: "text",
							Text: forecastText.String(),
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
