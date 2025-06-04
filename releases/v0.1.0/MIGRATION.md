# Migration Guide - v0.1.0

## Overview
This is the initial release of go-fork/log package. No migration is needed as this is the first version.

## Installation

### Step 1: Install Package
```bash
go get go.fork.vn/log@v0.1.0
```

### Step 2: Basic Usage
```go
package main

import (
    "go.fork.vn/log"
    "go.fork.vn/di"
)

func main() {
    // Create DI application
    app := di.New()
    
    // Register log service provider
    app.Register(log.NewServiceProvider())
    
    // Get logging manager
    container := app.Container()
    logger := container.MustMake("log").(log.Manager)
    
    // Use logging
    logger.Info("Application started")
    logger.Error("Something went wrong", map[string]interface{}{
        "error": "sample error",
        "user":  "john_doe",
    })
}
```

## Getting Started
- Check the [documentation](https://pkg.go.dev/go.fork.vn/log@v0.1.0)
- Review examples in the repository
- Read the README.md for detailed usage instructions

---
**Welcome to go-fork/log!** This is the beginning of a powerful logging management solution.
