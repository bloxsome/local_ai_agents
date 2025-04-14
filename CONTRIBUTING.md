# Contributing to Local_AI_Agents

Thank you for your interest in contributing to Local_AI_Agents! This document provides guidelines and instructions for contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
  - [Development Environment](#development-environment)
  - [Project Structure](#project-structure)
- [How to Contribute](#how-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Pull Requests](#pull-requests)
- [Development Workflow](#development-workflow)
  - [Branching Strategy](#branching-strategy)
  - [Commit Messages](#commit-messages)
  - [Testing](#testing)
  - [Documentation](#documentation)
- [Style Guidelines](#style-guidelines)
  - [Go Code Style](#go-code-style)
  - [Documentation Style](#documentation-style)
- [Community](#community)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

### Development Environment

1. **Prerequisites**:
   - Go 1.18 or higher
   - Git
   - Ollama (for running local AI models)

2. **Setup**:
   ```bash
   # Clone the repository
   git clone https://github.com/yourusername/Local_AI_Agents.git
   cd Local_AI_Agents

   # Install dependencies
   go mod download

   # Build the application
   go build
   ```

3. **Running Tests**:
   ```bash
   # Run all tests
   cd tests
   ./run_tests.sh  # On Unix-like systems
   run_tests.bat   # On Windows
   ```

### Project Structure

- `cmd/` - Command-line applications
- `docs/` - Documentation files
- `examples/` - Example code and usage
- `internal/` - Internal packages
  - `modules/` - Core modules (agents, context, commands, etc.)
  - `config/` - Configuration handling
- `tests/` - Test files and utilities

## How to Contribute

### Reporting Bugs

1. **Check Existing Issues**: Before creating a new issue, please check if it already exists.

2. **Create a Bug Report**: If your issue is not already reported, create a new issue using the bug report template. Include:
   - A clear title and description
   - Steps to reproduce the issue
   - Expected behavior
   - Actual behavior
   - Environment details (OS, Go version, etc.)
   - Any relevant logs or screenshots

### Suggesting Enhancements

1. **Check Existing Issues**: Before suggesting an enhancement, please check if it has already been suggested.

2. **Create an Enhancement Request**: If your enhancement is not already suggested, create a new issue using the feature request template. Include:
   - A clear title and description
   - The problem your enhancement solves
   - How your enhancement would work
   - Any alternatives you've considered
   - Any relevant examples or mockups

### Pull Requests

1. **Fork the Repository**: Fork the repository to your GitHub account.

2. **Create a Branch**: Create a branch for your changes.
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Changes**: Make your changes following the style guidelines.

4. **Write Tests**: Add tests for your changes to ensure they work as expected.

5. **Run Tests**: Ensure all tests pass.
   ```bash
   cd tests
   ./run_tests.sh
   ```

6. **Update Documentation**: Update any relevant documentation.

7. **Commit Changes**: Commit your changes with a clear commit message.
   ```bash
   git commit -m "Add feature: your feature description"
   ```

8. **Push Changes**: Push your changes to your fork.
   ```bash
   git push origin feature/your-feature-name
   ```

9. **Create a Pull Request**: Create a pull request from your fork to the main repository. Include:
   - A clear title and description
   - Reference to any related issues
   - Changes you've made
   - Any additional information that might be helpful

10. **Code Review**: Participate in the code review process, making changes as requested.

## Development Workflow

### Branching Strategy

- `main`: The main branch contains the stable version of the code.
- `feature/*`: Feature branches for new features or enhancements.
- `bugfix/*`: Bugfix branches for bug fixes.
- `docs/*`: Documentation branches for documentation changes.

### Commit Messages

Write clear, concise commit messages that explain what changes were made and why. Follow this format:

```
<type>: <subject>

<body>
```

Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or modifying tests
- `chore`: Maintenance tasks

Example:
```
feat: Add support for multiple agent personalities

Implement a system for managing different agent personalities with customizable traits, prompt templates, and response formatting.
```

### Testing

- Write tests for all new features and bug fixes.
- Ensure all tests pass before submitting a pull request.
- Aim for high test coverage.

### Documentation

- Update documentation for all new features and changes.
- Document public APIs, modules, and important functions.
- Keep the README.md up to date.

## Style Guidelines

### Go Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) guidelines.
- Use `gofmt` to format your code.
- Follow the standard Go naming conventions:
  - Use camelCase for variable and function names.
  - Use PascalCase for exported names.
  - Use all lowercase for package names.
- Write clear, concise comments for functions and complex code blocks.

### Documentation Style

- Use Markdown for documentation.
- Keep documentation clear, concise, and up to date.
- Include examples where appropriate.
- Use proper headings, lists, and code blocks.

## Community

- Be respectful and considerate of others.
- Help others who have questions or issues.
- Share your knowledge and expertise.
- Provide constructive feedback.

Thank you for contributing to Local_AI_Agents!
