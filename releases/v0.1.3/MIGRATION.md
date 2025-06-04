# Migration Guide - v0.1.3

## Overview
This guide helps you migrate from v0.1.2 to v0.1.3.

## Prerequisites
- Go 1.23 or later
- Previous version (v0.1.2) installed

## Quick Migration Checklist
- [ ] Update dependencies to v0.1.3
- [ ] Update test function names (if you extend log package tests)
- [ ] Run tests to ensure compatibility
- [ ] Update documentation references

## Breaking Changes

### Test Function Naming Convention
If you're extending the log package or writing custom tests that follow our patterns:

#### Test Functions
```go
// Old naming pattern (v0.1.2 and earlier)
func TestValidateConfig(t *testing.T) { ... }
func TestNewManager(t *testing.T) { ... }
func TestLogMethod(t *testing.T) { ... }

// New naming pattern (v0.1.3+)
func TestConfig_Validate(t *testing.T) { ... }
func TestManager_New(t *testing.T) { ... }
func TestManager_LogMethods(t *testing.T) { ... }
```

#### Benchmark Functions
```go
// Old naming pattern (v0.1.2 and earlier)
func BenchmarkValidateConfig(b *testing.B) { ... }
func BenchmarkNewManager(b *testing.B) { ... }

// New naming pattern (v0.1.3+)
func BenchmarkConfig_Validate(b *testing.B) { ... }
func BenchmarkManager_New(b *testing.B) { ... }
```

### ServiceProvider Changes
**Internal API Enhancement** - No breaking changes for users:
- Enhanced container nil validation in `Boot()` method
- Improved error handling consistency between `Register()` and `Boot()` methods
## Step-by-Step Migration

### Step 1: Update Dependencies
```bash
go get go.fork.vn/log@v0.1.3
go mod tidy
```

### Step 2: Update Test Function Names (If Applicable)
If you have custom tests following our patterns, update them:

```go
// Before
func TestValidateMyConfig(t *testing.T) { ... }

// After  
func TestMyConfig_Validate(t *testing.T) { ... }
```

### Step 3: Verify ServiceProvider Usage
No changes needed - the fix is internal and maintains compatibility.

### Step 4: Run Tests
```bash
go test ./...
```

## Common Issues and Solutions

### Issue 1: Test Function Names Out of Sync
**Problem**: Custom tests not following new convention  
**Solution**: Update to `Test{TypeName}_{MethodName}[_{Scenario}]` pattern

### Issue 2: Dependency Version Conflicts
**Problem**: `go.fork.vn/config` or `go.fork.vn/di` version mismatch  
**Solution**: Update to v0.1.3: `go get go.fork.vn/config@v0.1.3 go.fork.vn/di@v0.1.3`

## Getting Help
- Check the [documentation](https://pkg.go.dev/go.fork.vn/log@v0.1.3)
- Search [existing issues](https://github.com/go-fork/log/issues)
- Create a [new issue](https://github.com/go-fork/log/issues/new) if needed

## Rollback Instructions
If you need to rollback:

```bash
go get go.fork.vn/log@v0.1.2
go mod tidy
```

---
**Need Help?** Feel free to open an issue or discussion on GitHub.
