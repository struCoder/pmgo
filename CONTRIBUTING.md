# Contributing to PMGO

Thank you for your interest in contributing to PMGO! This document provides guidelines for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Style Guidelines](#style-guidelines)

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your changes
4. Make your changes
5. Test your changes
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make

### Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/pmgo.git
cd pmgo

# Install dependencies
make deps

# Build the project
make build

# Run tests
make test
```

### Development Tools

We recommend using these tools for development:

- [golangci-lint](https://golangci-lint.run/) for linting
- [air](https://github.com/cosmtrek/air) for live reloading
- [swag](https://github.com/swaggo/swag) for API documentation

Install them with:

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install air
go install github.com/cosmtrek/air@latest

# Install swag
go install github.com/swaggo/swag/cmd/swag@latest
```

## Making Changes

### Branch Naming

Use descriptive branch names:

- `feature/add-web-interface`
- `bugfix/fix-memory-leak`
- `docs/update-readme`

### Commit Messages

Follow conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(api): add process restart endpoint
fix(daemon): resolve memory leak in watcher
docs(readme): update installation instructions
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/master

# Run tests with verbose output
go test -v ./...
```

### Writing Tests

- Write unit tests for all new functionality
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for >80% test coverage

Example test structure:

```go
func TestProcessManager_Start(t *testing.T) {
    tests := []struct {
        name    string
        input   ProcessConfig
        want    error
        wantErr bool
    }{
        {
            name: "valid process",
            input: ProcessConfig{
                Name:    "test",
                Command: "echo",
                Args:    []string{"hello"},
            },
            want:    nil,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            pm := NewProcessManager()
            err := pm.Start(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if err != nil && err.Error() != tt.want.Error() {
                t.Errorf("Start() error = %v, want %v", err, tt.want)
            }
        })
    }
}
```

## Submitting Changes

### Pull Request Process

1. Update documentation if needed
2. Add tests for new functionality
3. Ensure all tests pass
4. Update CHANGELOG.md
5. Submit pull request with clear description

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
```

## Style Guidelines

### Go Code Style

Follow standard Go conventions:

- Use `gofmt` for formatting
- Follow naming conventions
- Use meaningful variable names
- Add comments for exported functions
- Keep functions small and focused

### API Design

- Use RESTful principles
- Return consistent error formats
- Use appropriate HTTP status codes
- Include API versioning
- Provide comprehensive documentation

### Documentation

- Use clear, concise language
- Include examples where helpful
- Keep documentation up to date
- Follow markdown best practices

## Project Structure

```
pmgo/
├── cmd/                    # Application entrypoints
│   └── pmgo/              # Main application
├── internal/              # Private application code
│   ├── api/              # HTTP API handlers
│   ├── config/           # Configuration management
│   ├── daemon/           # Daemon implementation
│   ├── master/           # Process master
│   ├── process/          # Process management
│   └── web/              # Web interface
├── pkg/                   # Public library code
├── configs/              # Configuration files
├── docs/                 # Documentation
├── scripts/              # Build and deployment scripts
└── tests/                # Integration tests
```

## Release Process

1. Update version in code
2. Update CHANGELOG.md
3. Create release branch
4. Run full test suite
5. Build and test binaries
6. Create GitHub release
7. Update documentation

## Getting Help

- Check existing issues and documentation
- Ask questions in GitHub Discussions
- Join our Discord server
- Contact maintainers directly

Thank you for contributing to PMGO!