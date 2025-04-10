// Package errors provides custom error types for the Ollama Agents application
package errors

// OllamaAgentsError is the base error type for Ollama Agents
type OllamaAgentsError struct {
	Message string
}

// Error returns the error message
func (e OllamaAgentsError) Error() string {
	return e.Message
}

// NewOllamaAgentsError creates a new OllamaAgentsError
func NewOllamaAgentsError(message string) OllamaAgentsError {
	return OllamaAgentsError{Message: message}
}

// ConfigurationError is raised when there's a configuration issue
type ConfigurationError struct {
	OllamaAgentsError
}

// NewConfigurationError creates a new ConfigurationError
func NewConfigurationError(message string) ConfigurationError {
	return ConfigurationError{OllamaAgentsError{Message: message}}
}

// APIConnectionError is raised when there's an issue connecting to an API
type APIConnectionError struct {
	OllamaAgentsError
}

// NewAPIConnectionError creates a new APIConnectionError
func NewAPIConnectionError(message string) APIConnectionError {
	return APIConnectionError{OllamaAgentsError{Message: message}}
}

// InputError is raised when there's an issue with user input
type InputError struct {
	OllamaAgentsError
}

// NewInputError creates a new InputError
func NewInputError(message string) InputError {
	return InputError{OllamaAgentsError{Message: message}}
}

// MemoryError is raised when there's an issue with memory operations
type MemoryError struct {
	OllamaAgentsError
}

// NewMemoryError creates a new MemoryError
func NewMemoryError(message string) MemoryError {
	return MemoryError{OllamaAgentsError{Message: message}}
}

// FileOperationError is raised when there's an issue with file operations
type FileOperationError struct {
	OllamaAgentsError
}

// NewFileOperationError creates a new FileOperationError
func NewFileOperationError(message string) FileOperationError {
	return FileOperationError{OllamaAgentsError{Message: message}}
}

// CommandExecutionError is raised when there's an error executing a command
type CommandExecutionError struct {
	OllamaAgentsError
}

// NewCommandExecutionError creates a new CommandExecutionError
func NewCommandExecutionError(message string) CommandExecutionError {
	return CommandExecutionError{OllamaAgentsError{Message: message}}
}

// LogicProcessingError is raised when there's an error in the logic processing steps
type LogicProcessingError struct {
	OllamaAgentsError
}

// NewLogicProcessingError creates a new LogicProcessingError
func NewLogicProcessingError(message string) LogicProcessingError {
	return LogicProcessingError{OllamaAgentsError{Message: message}}
}

// ModelInferenceError is raised when there's an error during model inference
type ModelInferenceError struct {
	OllamaAgentsError
}

// NewModelInferenceError creates a new ModelInferenceError
func NewModelInferenceError(message string) ModelInferenceError {
	return ModelInferenceError{OllamaAgentsError{Message: message}}
}

// DataProcessingError is raised when there's an error processing data
type DataProcessingError struct {
	OllamaAgentsError
}

// NewDataProcessingError creates a new DataProcessingError
func NewDataProcessingError(message string) DataProcessingError {
	return DataProcessingError{OllamaAgentsError{Message: message}}
}
