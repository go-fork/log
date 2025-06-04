# Log Package Documentation

The log package provides structured logging capabilities for Fork Framework applications, integrating seamlessly with the dependency injection container and service provider pattern.

## Quick Navigation

- **[Overview](overview.md)** - Architecture and logging concepts
- **[Usage Guide](usage.md)** - Practical examples and patterns
- **[API Reference](reference.md)** - Complete interface documentation

## Quick Start

```go
package main

import (
    "go.fork.vn/core"
)

func main() {
    // Create Fork application
    config := map[string]interface{}{
	    "name": "myapp",
	    "path": "./configs",
	}
	app := app.New(config)
    
    // Access logger through DI container
    logger := app.MustMake("logger").(log.Logger)
    
    // Basic logging
    logger.Info("Application started")
    logger.Error("Something went wrong", "error", err)
    logger.Debug("Debug information", "user_id", 123)
    
    app.Run()
}
```

## Key Features

- **Framework Integration**: Seamless integration with Fork Framework's DI container
- **Structured Logging**: Support for structured log entries with key-value pairs
- **Multiple Levels**: Support for Debug, Info, Warn, Error, and Fatal log levels
- **Contextual Logging**: Add context to log entries for better debugging
- **Configurable Output**: Support for console, file, and custom output destinations
- **Performance**: Optimized for high-performance applications

## Installation

```bash
go get go.fork.vn/log@latest
```

The log service is automatically registered when using Fork Framework.

## Package Structure

```
log/
├── docs/
│   ├── index.md      # This file
│   ├── overview.md   # Architecture and concepts
│   ├── usage.md      # Usage examples
│   └── reference.md  # API documentation
├── mocks/            # Generated mocks
├── logger.go         # Core logger interface
├── manager.go        # Logger implementation
├── service_provider.go # Framework integration
└── doc.go           # Package documentation
```

## Log Levels

The package supports standard log levels:

- **Debug**: Detailed information for diagnosing problems
- **Info**: General information about application flow
- **Warn**: Warning messages for potentially harmful situations
- **Error**: Error events that might still allow the application to continue
- **Fatal**: Very severe error events that will presumably lead the application to abort

## Getting Help

- Read the [Usage Guide](usage.md) for practical examples
- Check the [API Reference](reference.md) for detailed interface documentation
- Review the [Overview](overview.md) for architectural concepts
