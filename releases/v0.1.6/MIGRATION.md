# Migration Guide - v0.1.6

## Overview
This guide helps you migrate from v0.1.5 to v0.1.6, focusing on the **configuration architecture improvements** and **test infrastructure changes**.

## Prerequisites
- Go 1.23 or later
- v0.1.5 or earlier installed

## Quick Migration Checklist
- [ ] Review file path configuration if using default config
- [ ] Update any hardcoded directory dependencies in tests
- [ ] Verify file logging setup if using default configuration
- [ ] No breaking changes - existing configurations still work

## ⚠️ **Important**: No Breaking Changes
**v0.1.6 is fully backward compatible**. Existing code will continue to work without modifications.

## Changes Overview

### 🔧 **Configuration Architecture Changes**

#### Default Configuration Updates
The main change affects users who rely on `DefaultConfig()` without setting custom file paths.

#### Before (v0.1.5 and earlier)
```go
config := log.DefaultConfig()
// config.File.Path was "storages/log/app.log" (hardcoded)
// config.File.Enabled was false but path validation still occurred
```

#### After (v0.1.6)
```go
config := log.DefaultConfig()
// config.File.Path is now "" (empty)
// config.File.Enabled is false and no path validation when disabled
```

### 📝 **Migration Scenarios**

## Scenario 1: Using Custom File Paths (Most Common)
### ✅ **No Action Required**

If you're already setting custom file paths, this release doesn't affect you:

```go
// This code continues to work exactly the same
config := log.DefaultConfig()
config.File.Enabled = true
config.File.Path = "/var/log/myapp/app.log"  // Custom path
config.File.MaxSize = 50 * 1024 * 1024      // 50MB

manager, err := log.NewManager(config)
// No changes needed
```

## Scenario 2: Using Default Config with Console Only
### ✅ **No Action Required**

If you only use console logging, nothing changes:

```go
// This continues to work the same
config := log.DefaultConfig()
// config.File.Enabled is false by default
// Only console logging is active

manager, err := log.NewManager(config)
// No changes needed
```

## Scenario 3: Using Default Config with File Logging
### ⚠️ **Minor Update Required**

If you enabled file logging without setting a custom path (relying on the old hardcoded path):

#### Before (v0.1.5)
```go
config := log.DefaultConfig()
config.File.Enabled = true
// Used hardcoded "storages/log/app.log"

manager, err := log.NewManager(config)
```

#### After (v0.1.6)
```go
config := log.DefaultConfig()
config.File.Enabled = true
config.File.Path = "/your/preferred/path/app.log" // Must set explicitly

// Ensure directory exists (same as before)
logDir := filepath.Dir(config.File.Path)
if err := os.MkdirAll(logDir, 0755); err != nil {
    log.Fatal("Failed to create log directory:", err)
}

manager, err := log.NewManager(config)
```

## Scenario 4: Test Suites
### ✅ **Enhanced Reliability**

Tests that previously depended on hardcoded directories now skip gracefully:

#### Before (v0.1.5)
```go
// Tests might fail if storages/log/ didn't exist
func TestSomething(t *testing.T) {
    config := log.DefaultConfig()
    config.File.Enabled = true
    // Would try to validate "storages/log/app.log"
}
```

#### After (v0.1.6)
```go
// Tests skip gracefully when directories don't exist
func TestSomething(t *testing.T) {
    config := log.DefaultConfig()
    config.File.Enabled = true
    config.File.Path = filepath.Join(t.TempDir(), "app.log") // Use temp dir
    // More reliable testing
}
```

## 🛠️ **Recommended Migration Steps**

### Step 1: Update Dependencies
```bash
go get go.fork.vn/log@v0.1.6
go mod tidy
```

### Step 2: Review Configuration
Check if you're using `DefaultConfig()` with file logging:

```bash
# Search for potential usage patterns
grep -r "DefaultConfig" . --include="*.go"
grep -r "File.Enabled.*true" . --include="*.go"
```

### Step 3: Update File Paths (If Needed)
Only if you were relying on the hardcoded `storages/log/app.log` path:

```go
// Add explicit path configuration
config := log.DefaultConfig()
config.File.Enabled = true
config.File.Path = "/your/preferred/path/app.log" // Add this line
```

### Step 4: Test Your Application
```bash
go test ./...
go run your-application
```

## 🧪 **Testing Changes**

### Enhanced Test Reliability
Tests now handle missing directories more gracefully:

```go
// Tests may now skip instead of fail
// This is expected and indicates better portability
=== RUN   TestConfig_Validate/storages_directory_test
    config_test.go:XXX: Skipping test because directory does not exist: storages/log
--- SKIP: TestConfig_Validate/storages_directory_test (0.00s)
```

### For Test Development
Use temporary directories for more reliable tests:

```go
func TestWithFileLogging(t *testing.T) {
    // Use temporary directory
    tempDir := t.TempDir()
    
    config := log.DefaultConfig()
    config.File.Enabled = true
    config.File.Path = filepath.Join(tempDir, "test.log")
    
    // Test with reliable path
    manager, err := log.NewManager(config)
    assert.NoError(t, err)
}
```

## 📊 **Benefits After Migration**

### **Enhanced Security**
- No accidental file creation in unexpected locations
- Explicit configuration prevents surprise behaviors
- Better control over file system interactions

### **Improved Portability**
- Tests work across different environments
- No mandatory directory structures
- Better CI/CD compatibility

### **Cleaner Development**
- No need to create specific directories during development
- Repository stays clean without test artifacts
- More flexible development setup

## 🔍 **Troubleshooting**

### Issue: "path is required for file handler initialization"
**Cause**: Using file logging with empty path  
**Solution**: Set explicit file path
```go
config.File.Path = "/your/path/app.log"
```

### Issue: Tests failing with directory errors
**Cause**: Hardcoded directory dependencies  
**Solution**: Use temporary directories or skip tests
```go
// Option 1: Use temp directory
config.File.Path = filepath.Join(t.TempDir(), "app.log")

// Option 2: Skip when directory missing
if _, err := os.Stat(dir); os.IsNotExist(err) {
    t.Skip("Directory does not exist:", dir)
}
```

### Issue: Application can't find log files
**Cause**: Changed default path  
**Solution**: Set explicit path in configuration
```go
config.File.Path = "/var/log/myapp/app.log" // Your preferred location
```

## 🚀 **Rollback Instructions**

If you need to rollback to v0.1.5:

```bash
go get go.fork.vn/log@v0.1.5
go mod tidy
```

The old hardcoded path behavior will be restored.

## 🎯 **Validation**

After migration, verify your setup:

```bash
# 1. Dependencies updated
go list -m go.fork.vn/log

# 2. Tests pass
go test ./...

# 3. Application runs correctly
go run your-main.go

# 4. Log files created in expected locations
ls -la /your/log/directory/
```

## 📞 **Getting Help**

- Check the [documentation](https://pkg.go.dev/go.fork.vn/log@v0.1.6)
- Search [existing issues](https://github.com/go-fork/log/issues)
- Create a [new issue](https://github.com/go-fork/log/issues/new) if needed

---

**Migration Difficulty**: 🟢 **Easy** (Mostly backward compatible)  
**Time Required**: 🕐 **5-15 minutes** (Only if using default file paths)  
**Risk Level**: 🟢 **Low** (No breaking changes)
