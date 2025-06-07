# Migration Guide - v0.1.5

## Overview
This guide helps you migrate from v0.1.4 to v0.1.5, focusing on the breaking changes in file handler validation and directory requirements.

## Prerequisites
- Go 1.23 or later
- v0.1.4 or earlier installed

## Quick Migration Checklist
- [ ] Create directories before initializing file handlers
- [ ] Update any code that relied on automatic directory creation
- [ ] Verify write permissions for log directories
- [ ] Update test cases to use snake_case naming (if developing)
- [ ] Run tests to ensure compatibility

## Breaking Changes

### File Handler Directory Requirements
The most significant change in v0.1.5 is that `NewFileHandler` no longer automatically creates directories.

#### Before (v0.1.4 and earlier)
```go
// This worked even if /var/log/myapp didn't exist
handler, err := handler.NewFileHandler("/var/log/myapp/app.log", 10*1024*1024)
if err != nil {
    log.Fatal("Failed to create file handler:", err)
}
```

#### After (v0.1.5)
```go
// You must ensure the directory exists first
logDir := "/var/log/myapp"
if err := os.MkdirAll(logDir, 0755); err != nil {
    log.Fatal("Failed to create log directory:", err)
}

// Now create the file handler
handler, err := handler.NewFileHandler("/var/log/myapp/app.log", 10*1024*1024)
if err != nil {
    log.Fatal("Failed to create file handler:", err)
}
```

### Configuration Changes
#### Default Configuration
The default configuration now sets `File.Enabled: false` by default.

#### Before (v0.1.4 and earlier)
```go
config := log.DefaultConfig()
// config.File.Enabled was true by default
```

#### After (v0.1.5)
```go
config := log.DefaultConfig()
// config.File.Enabled is now false by default
// You must explicitly enable file logging:
config.File.Enabled = true
config.File.Path = "/path/to/log/file.log"
```
    Field1 string
    Field2 int64 // Changed from int
    Field3 bool  // New field
}
```

### Configuration Changes
If you're using configuration files:

```yaml
# Old configuration format
old_setting: value
deprecated_option: true

# New configuration format
new_setting: value
# deprecated_option removed
new_option: false
```

## Step-by-Step Migration

### Step 1: Update Dependencies
```bash
go get go.fork.vn/log@v0.1.5
go mod tidy
```

### Step 2: Update Import Statements
```go
// If import paths changed
import (
    "go.fork.vn/log" // Updated import
)
```

### Step 3: Update Code
Replace deprecated function calls:

```go
// Before
result := log.OldFunction(param)

// After
result := log.NewFunction(param, defaultValue)
```

### Step 4: Update Configuration
Update your configuration files according to the new schema.

### Step 5: Run Tests
```bash
go test ./...
```

## Common Issues and Solutions

### Issue 1: Function Not Found
**Problem**: `undefined: log.OldFunction`  
**Solution**: Replace with `log.NewFunction`

### Issue 2: Type Mismatch
**Problem**: `cannot use int as int64`  
**Solution**: Cast the value or update variable type

## Getting Help
- Check the [documentation](https://pkg.go.dev/go.fork.vn/log@v0.1.5)
- Search [existing issues](https://github.com/go-fork/log/issues)
- Create a [new issue](https://github.com/go-fork/log/issues/new) if needed

## Rollback Instructions
If you need to rollback:

```bash
go get go.fork.vn/log@previous-version
go mod tidy
```

Replace `previous-version` with your previous version tag.

---
**Need Help?** Feel free to open an issue or discussion on GitHub.
