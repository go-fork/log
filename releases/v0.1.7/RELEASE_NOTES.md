# Release Notes - v0.1.7

## Overview
This release focuses on fixing a critical duplicate logging issue and improving the Stack Handler architecture. The main achievement is resolving the problem where logs appeared twice when both individual handlers and stack handlers were enabled simultaneously.

## What's New
### 🐛 Bug Fixes
- **Fixed Stack Handler Duplicate Logging**: Resolved critical issue where log messages were output twice when both individual handlers (console/file) and stack handler were enabled
- **Enhanced Stack Handler Configuration**: Stack handler now properly respects the `stack.handlers.console` and `stack.handlers.file` configuration settings

### 🔧 Improvements
- **Smart Handler Selection Logic**: Implemented intelligent handler selection to prevent duplicate outputs
- **Conditional Handler Addition**: Individual handlers are now only added when not already covered by the stack handler
- **Better Configuration Granularity**: Stack handler configuration is now more predictable and follows user intent
- **Improved Handler Separation**: Better separation of concerns between individual handlers and composite stack handler

### 📚 Documentation
- Enhanced inline documentation for handler logic
- Improved comments explaining the anti-duplication strategy

## Breaking Changes
### ⚠️ Important Notes
**No Breaking Changes** - This release maintains full backward compatibility. All existing configurations continue to work with improved behavior.

## Migration Guide
No migration required for this release. Existing code and configurations will work without changes and benefit from the improved logging behavior.

## Technical Implementation
### Handler Logic Improvements
- **Stack Handler Initialization**: Now respects individual handler configurations in `initializeHandlers()`
- **Logger Creation Logic**: Enhanced `GetLogger()` method with anti-duplication logic
- **Configuration Priority**: Stack handler can now act as a true composite handler replacement

### Anti-Duplication Strategy
```go
// Individual handlers only added when not covered by stack
if !m.config.Stack.Enabled || (m.config.Console.Enabled && !m.config.Stack.Handlers.Console) {
    // Add console handler
}
```

## Performance
- **Reduced Output Overhead**: Eliminated duplicate log writes, improving performance
- **Memory Efficiency**: Reduced redundant handler operations

## Testing
- All existing tests continue to pass
- Enhanced test coverage for stack handler configurations
- Verified fix with multiple configuration scenarios

## Contributors
Thanks to all contributors who made this release possible:
- @contributor1
- @contributor2

## Download
- Source code: [go.fork.vn/log@v0.1.7]
- Documentation: [pkg.go.dev/go.fork.vn/log@v0.1.7]

---
Release Date: 2025-06-07
