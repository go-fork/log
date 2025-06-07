# Release Notes - v0.1.5

## Overview
This release focuses on enhanced file handler validation, security improvements, and test infrastructure consistency. The main changes include stricter directory validation requirements and comprehensive test naming standardization.

## What's New
### 🔒 Security Enhancements
- Enhanced `NewFileHandler` directory validation - now requires existing directories
- Removed automatic directory creation for better security and explicit control
- Added comprehensive write permission validation before file operations

### 🐛 Bug Fixes
- Fixed directory validation logic to provide clear error messages
- Enhanced error handling for file system operations during validation
- Improved DefaultConfig to set File.Enabled: false by default for safer defaults

### 🔧 Improvements
- Updated Config.Validate() to always check File.Path regardless of File.Enabled status
- Enhanced error messages with specific validation failures
- Improved file handler initialization with detailed permission checks

### 📚 Documentation
- Updated configuration documentation to reflect new validation behavior
- Added comprehensive validation examples and error scenarios
- Enhanced CHANGELOG with detailed breaking changes documentation

## Breaking Changes
### ⚠️ Important Notes
- **NewFileHandler Behavior Change**: Now requires the directory to exist beforehand
- **No Automatic Directory Creation**: Directories must be created manually before initializing file handlers
- **Enhanced Validation**: File.Path is always validated regardless of File.Enabled status

## Migration Guide
See [MIGRATION.md](./MIGRATION.md) for detailed migration instructions.

## Dependencies
No dependency changes in this release.

## Performance
- Enhanced validation performance with early directory checks
- Reduced runtime errors through comprehensive upfront validation

## Security
- **Directory Security**: Removed automatic directory creation prevents unintended file system modifications
- **Permission Validation**: Added comprehensive write permission checks before file operations
- **Clear Error Messages**: Enhanced error reporting for better debugging and security auditing

## Testing
- Added comprehensive test cases for directory validation scenarios
- Converted all test case names to snake_case format for consistency
- Enhanced test coverage for file handler edge cases and error conditions
- Improved benchmark test naming conventions
- Added validation scenarios for non-existent directories and write permissions

## Contributors
Thanks to all contributors who made this release possible:
- @ntnghia0921
- @zinzinday

## Download
- Source code: [go.fork.vn/log@v0.1.5]
- Documentation: [pkg.go.dev/go.fork.vn/log@v0.1.5]

---
Release Date: 2025-06-07
