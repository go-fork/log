# Migration Guide - v0.1.7

## Overview
This is a **bug fix release** with **zero breaking changes**. No migration is required.

## Prerequisites
- Go 1.23 or later
- Previous version installed

## Quick Migration Checklist
- [x] ✅ **No changes required** - This release maintains full backward compatibility
- [x] ✅ **No API changes** - All function signatures remain the same
- [x] ✅ **No configuration changes** - All existing configs work better than before
- [x] ✅ **No code updates needed** - Your existing code will benefit automatically

## Breaking Changes
**None** - This release has zero breaking changes.

## What's Improved (No Action Required)
### Stack Handler Behavior
Your existing stack handler configurations will now work more predictably:

```yaml
# This configuration now works correctly without duplication
log:
  console:
    enabled: true
  stack:
    enabled: true
    handlers:
      console: true  # No longer causes duplicate output
      file: true
```

### Before v0.1.7 (Problem)
```
2025/06/07 15:04:34 [INFO] [test] Message  # From console handler
2025/06/07 15:04:34 [INFO] [test] Message  # From stack->console (duplicate!)
```

### After v0.1.7 (Fixed)
```
2025/06/07 15:04:34 [INFO] [test] Message  # Clean, single output
```

## Step-by-Step Migration

### Step 1: Update Dependencies
```bash
go get go.fork.vn/log@v0.1.7
go mod tidy
```

### Step 2: Enjoy the Fix
That's it! Your code will automatically benefit from:
- ✅ No more duplicate log entries
- ✅ Better stack handler behavior
- ✅ Improved performance (fewer redundant writes)

### Step 3: Optional - Simplify Configuration
If you were working around the duplication issue, you can now simplify:

```yaml
# You can now safely use both without duplication
console:
  enabled: true
stack:
  enabled: true
  handlers:
    console: true
    file: true
```

## Common Improvements You'll Notice

### Improvement 1: Clean Log Output
**Before**: Duplicate log entries with mixed configurations  
**After**: Clean, single log entries as expected

### Improvement 2: Predictable Stack Handler
**Before**: Stack handler always included all handlers regardless of config  
**After**: Stack handler respects your `handlers.console` and `handlers.file` settings

## Getting Help
- Check the [documentation](https://pkg.go.dev/go.fork.vn/log@v0.1.7)
- Search [existing issues](https://github.com/go-fork/log/issues)
- Create a [new issue](https://github.com/go-fork/log/issues/new) if needed

## Rollback Instructions
If you need to rollback (not recommended):

```bash
go get go.fork.vn/log@v0.1.6
go mod tidy
```

**Note**: Rolling back will restore the duplicate logging issue.

---
**Need Help?** This release should "just work" - but feel free to open an issue if you notice anything unexpected.

## Step-by-Step Migration

### Step 1: Update Dependencies
```bash
go get go.fork.vn/log@v0.1.7
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
- Check the [documentation](https://pkg.go.dev/go.fork.vn/log@v0.1.7)
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
