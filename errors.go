// This file is deprecated. Error types have been moved to internal/errors/errors.go
// This file is kept for backward compatibility and will be removed in a future version.
package main

import (
	"github.com/bloxsome/local_ai_agents/internal/errors"
)

// OllamaAgentsError is an alias for errors.OllamaAgentsError
type OllamaAgentsError = errors.OllamaAgentsError

// NewOllamaAgentsError creates a new OllamaAgentsError
func NewOllamaAgentsError(message string) OllamaAgentsError {
	return errors.NewOllamaAgentsError(message)
}

// ConfigurationError is an alias for errors.ConfigurationError
type ConfigurationError = errors.ConfigurationError

// NewConfigurationError creates a new ConfigurationError
func NewConfigurationError(message string) ConfigurationError {
	return errors.NewConfigurationError(message)
}

// APIConnectionError is an alias for errors.APIConnectionError
type APIConnectionError = errors.APIConnectionError

// NewAPIConnectionError creates a new APIConnectionError
func NewAPIConnectionError(message string) APIConnectionError {
	return errors.NewAPIConnectionError(message)
}

// InputError is an alias for errors.InputError
type InputError = errors.InputError

// NewInputError creates a new InputError
func NewInputError(message string) InputError {
	return errors.NewInputError(message)
}

// MemoryError is an alias for errors.MemoryError
type MemoryError = errors.MemoryError

// NewMemoryError creates a new MemoryError
func NewMemoryError(message string) MemoryError {
	return errors.NewMemoryError(message)
}

// FileOperationError is an alias for errors.FileOperationError
type FileOperationError = errors.FileOperationError

// NewFileOperationError creates a new FileOperationError
func NewFileOperationError(message string) FileOperationError {
	return errors.NewFileOperationError(message)
}

// CommandExecutionError is an alias for errors.CommandExecutionError
type CommandExecutionError = errors.CommandExecutionError

// NewCommandExecutionError creates a new CommandExecutionError
func NewCommandExecutionError(message string) CommandExecutionError {
	return errors.NewCommandExecutionError(message)
}

// LogicProcessingError is an alias for errors.LogicProcessingError
type LogicProcessingError = errors.LogicProcessingError

// NewLogicProcessingError creates a new LogicProcessingError
func NewLogicProcessingError(message string) LogicProcessingError {
	return errors.NewLogicProcessingError(message)
}

// ModelInferenceError is an alias for errors.ModelInferenceError
type ModelInferenceError = errors.ModelInferenceError

// NewModelInferenceError creates a new ModelInferenceError
func NewModelInferenceError(message string) ModelInferenceError {
	return errors.NewModelInferenceError(message)
}

// DataProcessingError is an alias for errors.DataProcessingError
type DataProcessingError = errors.DataProcessingError

// NewDataProcessingError creates a new DataProcessingError
func NewDataProcessingError(message string) DataProcessingError {
	return errors.NewDataProcessingError(message)
}
