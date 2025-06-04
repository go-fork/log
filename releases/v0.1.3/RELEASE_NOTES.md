# Release Notes - v0.1.3

## Overview
This release focuses on test function naming convention refactoring, dependency upgrades, and bug fixes to align with Fork Framework standards and improve code quality.

## What's New
### 🔧 Improvements
- **Test Function Naming Convention Refactoring**: Migrated all test functions to Fork Framework naming convention
  - Updated pattern from legacy naming to `Test{TypeName}_{MethodName}[_{Scenario}]`
  - Updated benchmark functions to `Benchmark{TypeName}_{MethodName}[_{Scenario}]`
  - Enhanced test readability and consistency across the codebase

- **Dependencies Updated**:
  - Upgraded `go.fork.vn/config` from v0.1.2 to v0.1.3
  - Upgraded `go.fork.vn/di` from v0.1.2 to v0.1.3
  - Updated indirect dependency `github.com/spf13/cast` to v1.9.2

### 🐛 Bug Fixes
- **ServiceProvider Boot Method**: Fixed panic handling for nil container validation
  - Added proper container nil check in `Boot()` method for consistency with `Register()`
  - Resolved test case `TestServiceProvider_Boot/application_with_nil_container` failure
  - Ensured proper error handling throughout service provider lifecycle

### 📚 Documentation
- Updated CHANGELOG.md with comprehensive change documentation
- Enhanced migration guide with specific naming convention examples

## Breaking Changes
### ⚠️ Important Notes
- **Test Function Naming**: If you're extending the log package with custom tests, update function names to follow `Test{TypeName}_{MethodName}[_{Scenario}]` pattern
- **Benchmark Function Naming**: Update benchmark functions to follow `Benchmark{TypeName}_{MethodName}[_{Scenario}]` pattern

## Migration Guide
See [MIGRATION.md](./MIGRATION.md) for detailed migration instructions.

## Dependencies
### Updated
- go.fork.vn/config: v0.1.2 → v0.1.3
- go.fork.vn/di: v0.1.2 → v0.1.3
- github.com/spf13/cast: v1.9.1 → v1.9.2 (indirect)

## Testing
- 100% test functions pass with new naming convention
- Zero issues from `go vet` and `golangci-lint`
- 13 files modified with 147 insertions, 113 deletions
- All existing functionality maintained with improved test organization

## Quality Metrics
- Test Coverage: Maintained at 100% pass rate
- Static Analysis: Zero warnings or errors
- Code Quality: Enhanced readability and consistency
- Performance: No regression in benchmark tests

## Examples
### Test Function Naming Examples
**Before (v0.1.2 and earlier)**:
```go
func TestValidateConfig(t *testing.T) { ... }
func TestNewManager(t *testing.T) { ... }
func BenchmarkValidateConfig(b *testing.B) { ... }
```

**After (v0.1.3+)**:
```go
func TestConfig_Validate(t *testing.T) { ... }
func TestManager_New(t *testing.T) { ... }
func BenchmarkConfig_Validate(b *testing.B) { ... }
```

## Download
- Source code: [go.fork.vn/log@v0.1.3]
- Module: `go get go.fork.vn/log@v0.1.3`
- Documentation: [pkg.go.dev/go.fork.vn/log@v0.1.4]

---
Release Date: 2025-06-04
