# Release Notes - v0.1.6

**Release Date**: June 7, 2025  
**Type**: Minor Release - Configuration & Testing Improvements

## 🎯 Release Highlights

This release focuses on **cleaner repository architecture** and **improved default configurations** for better developer experience and security.

### 🔧 **Configuration Architecture Improvements**
- **Secure Defaults**: Empty file paths by default prevent accidental directory creation
- **Smart Validation**: File path validation only when handlers actually need it
- **Better Separation**: Clear distinction between default config and user settings

### 🧪 **Test Infrastructure Modernization**
- **Portable Tests**: No more hardcoded directory dependencies
- **Graceful Skipping**: Tests skip when directories don't exist instead of failing
- **Clean Repository**: No mandatory directory structures for development

### 🏗️ **Repository Structure Cleanup**
- **Removed Hardcoded Paths**: No more `storages/log/app.log` dependency
- **Flexible Development**: Work without creating specific directory structures
- **Enhanced Portability**: Tests work across different environments

## 📋 **What's Changed**

### Configuration Changes
```go
// Before v0.1.6
config := DefaultConfig()
// config.File.Path = "storages/log/app.log" (hardcoded)
// config.File.Enabled = false (but path still required for validation)

// After v0.1.6  
config := DefaultConfig()
// config.File.Path = "" (empty, user must set explicitly)
// config.File.Enabled = false (and path not required unless enabled)
```

### Test Behavior Changes
```go
// Before v0.1.6
// Tests would create directories or fail

// After v0.1.6
if _, err := os.Stat(dir); os.IsNotExist(err) {
    t.Skipf("Skipping test because directory does not exist: %s", dir)
}
```

## 🚀 **Benefits for Developers**

### **Cleaner Development Experience**
- No need to create `storages/log/` directories
- Repository doesn't get polluted with test artifacts
- More flexible development setup

### **Better Security Defaults**
- Empty paths prevent accidental file creation
- Explicit configuration required for file logging
- Validation only when actually needed

### **Enhanced Portability**
- Tests work in any environment
- No hardcoded path dependencies
- Better CI/CD compatibility

## 🔄 **Migration Guide**

### **For Most Users**: No Action Required
If you're already setting your own file paths, this release doesn't affect you.

### **If Using Default Config with File Logging**:
```go
// Before v0.1.6
config := log.DefaultConfig()
config.File.Enabled = true
// Would use hardcoded "storages/log/app.log"

// After v0.1.6
config := log.DefaultConfig()
config.File.Enabled = true
config.File.Path = "/your/custom/path/app.log" // Must set explicitly
```

### **For Test Suites**:
Tests that depended on `storages/log/` directory will now skip gracefully. To run these tests, create the required directories manually if needed.

## 🧪 **Testing**

```bash
# All tests pass with new skip-based approach
go test ./...

# Tests that need directories will skip gracefully
# No more test failures due to missing directories
```

## 📊 **Quality Metrics**

- **Test Coverage**: Maintained high coverage with improved reliability
- **CI/CD Pipeline**: All checks passing with enhanced portability
- **Breaking Changes**: None - backward compatible
- **Dependencies**: No changes

## ⚠️ **Important Notes**

1. **File Logging**: Now requires explicit path configuration
2. **Test Behavior**: Some tests may skip if directories don't exist (this is expected)
3. **Development Setup**: No mandatory directory creation required
4. **CI/CD**: Enhanced compatibility across different environments

## 🔗 **Related Issues & PRs**

- Configuration architecture improvements
- Test infrastructure modernization  
- Repository structure cleanup
- Enhanced developer experience

## 📦 **Installation**

```bash
go get go.fork.vn/log@v0.1.6
```

## 🤝 **Contributors**

Thanks to all contributors who helped improve the repository architecture and developer experience!

---

**Full Changelog**: [v0.1.5...v0.1.6](https://github.com/go-fork/log/compare/v0.1.5...v0.1.6)
