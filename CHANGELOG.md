# Changelog

## [Unreleased]

## v0.1.3 - 2025-06-04

### Added
- **Project Structure Enhancement**
  - Added `.github/` directory with comprehensive CI/CD workflows
  - Added `releases/` directory with structured release management
  - Added `scripts/` directory with automation tools for release management

- **GitHub Integration**
  - CI workflow (`ci.yml`) for automated testing and quality checks
  - Release workflow (`release.yml`) for automated releases with proper tagging
  - Issue templates for bug reports and feature requests
  - Pull request template with comprehensive checklist
  - CODEOWNERS file for automatic review assignments
  - Dependabot configuration focused on critical dependencies

- **Release Management System**
  - Structured release documentation in `releases/` directory
  - Release notes, summaries, and migration guides for v0.1.0, v0.1.1, v0.1.2, v0.1.3
  - Template system for future releases in `releases/next/`
  - Automated release archiving and template generation scripts

- **Automation Scripts**
  - `create_release_templates.sh` - Generate release documentation templates
  - `archive_release.sh` - Archive completed releases and prepare for next version
  - Version management with semantic versioning support

### Improved
- **Documentation Structure**
  - Updated README.md with project structure information
  - Enhanced release management documentation
  - Improved migration guides for all versions
  - Better organization of release artifacts

- **Development Workflow**
  - Streamlined CI/CD pipeline with proper Go module handling
  - Focused dependency management (go.fork.vn/config, go.fork.vn/di)
  - Automated quality checks and testing
  - Consistent release process with documentation

- **Repository Organization**
  - Clear separation of concerns with dedicated directories
  - Standardized file structure following Go best practices
  - Improved maintainability with automation scripts

### Technical Details
- **CI/CD Pipeline**: Automated testing, linting, and release processes
- **Go Module**: Updated module proxy URLs for go.fork.vn/log
- **Dependencies**: Focused management of go.fork.vn dependencies
- **Release Process**: Automated with proper semantic versioning and documentation generation



### Improved
- **Test Function Naming Convention Refactoring**
  - Migrated all test functions to Fork Framework naming convention
  - Updated pattern from legacy naming to `Test{TypeName}_{MethodName}[_{Scenario}]`
  - Updated benchmark functions to `Benchmark{TypeName}_{MethodName}[_{Scenario}]`
  - Enhanced test readability and consistency across the codebase

- **Dependencies Updated**
  - Upgraded `go.fork.vn/config` from v0.1.2 to v0.1.3
  - Upgraded `go.fork.vn/di` from v0.1.2 to v0.1.3
  - Updated indirect dependency `github.com/spf13/cast` to v1.9.2

### Fixed
- **ServiceProvider Boot Method**: Fixed panic handling for nil container validation
  - Added proper container nil check in `Boot()` method for consistency with `Register()`
  - Resolved test case `TestServiceProvider_Boot/application_with_nil_container` failure
  - Ensured proper error handling throughout service provider lifecycle

### Technical Details
- **Dependencies**: Upgraded to go.fork.vn/config v0.1.3 and go.fork.vn/di v0.1.3
- **Test Coverage**: 100% test functions pass with new naming convention
- **Quality Checks**: Zero issues from `go vet` and `golangci-lint`
- **Code Quality**: 13 files modified with 147 insertions, 113 deletions

### Test Function Naming Examples
- **Before**: `TestValidateConfig()`, `TestNewManager()`, `TestLogMethod()`
- **After**: `TestConfig_Validate()`, `TestManager_New()`, `TestManager_LogMethods()`
- **Benchmark Before**: `BenchmarkValidateConfig()`, `BenchmarkNewManager()`
- **Benchmark After**: `BenchmarkConfig_Validate()`, `BenchmarkManager_New()`

## v0.1.2 - 2025-06-02

### Added

- **Enhanced Performance Monitoring**
  - ServiceProvider registration and bootstrapping benchmarks
  - Manager handler operations benchmarks
  - Configuration validation speed testing
  - Concurrent logging performance measurements
  - Memory allocation tracking across all operations

- **Advanced Benchmark Categories**
  - Creation benchmarks for object instantiation
  - Validation benchmarks for different scenarios
  - Concurrency benchmarks for parallel operations
  - Edge case benchmarks for error conditions
  - Memory footprint benchmarks with allocation tracking

### Improved
- **Test Infrastructure Reliability**
  - Fixed provider test expectations to match actual implementation
  - Enhanced mock management to prevent race conditions
  - Improved error handling in test scenarios
  - Better resource cleanup in benchmark tests

- **Code Quality and Compliance**
  - Resolved all static analysis warnings
  - Fixed unused variable warnings in benchmarks
  - Improved type safety in mock setups
  - Enhanced error message consistency

- **Performance Testing Coverage**
  - 100% benchmark coverage for public APIs
  - Comprehensive error path performance testing
  - Multi-level configuration validation benchmarks
  - Handler management performance testing

### Fixed
- **Provider Tests**: Corrected service registration expectations
- **Benchmark Stability**: Fixed mock conflicts and memory leaks
- **Log Level Validation**: Updated valid log level names in tests
- **Type Assertions**: Improved error handling in test scenarios
- **Resource Management**: Better cleanup in concurrent tests

### Technical Details
- **Benchmark Commands Added**:
  ```bash
  go test -bench=. -benchmem                    # All benchmarks with memory stats
  go test -bench=BenchmarkServiceProvider       # Provider-specific benchmarks
  go test -bench=BenchmarkManager               # Manager-specific benchmarks
  go test -bench=BenchmarkConfig                # Config-specific benchmarks
  ```

- **Performance Profiling Support**:
  ```bash
  go test -bench=. -cpuprofile=cpu.prof         # CPU profiling
  go test -bench=. -memprofile=mem.prof         # Memory profiling
  ```

- **Quality Metrics Maintained**:
  - Test coverage: 80%+
  - Zero static analysis warnings
  - Comprehensive benchmark coverage
  - Robust error handling

## v0.1.1 - 2025-06-02

### Changed
- Updated dependencies to latest versions
- Enhanced stability and performance improvements

### Added
- Comprehensive test suite for Config validation and error handling

## v0.1.0 - 2025-05-31

### Added
- **Logging Manager**: Comprehensive logging management system for Go applications
- **Multiple Log Levels**: Support for Debug, Info, Warning, Error, and Fatal levels
- **Multiple Handlers**: Console handler with color support, File handler with rotation, Stack handler for multiple destinations
- **Thread-Safe**: Safe concurrent logging from multiple goroutines
- **Printf-Style Formatting**: Support for placeholder formatting in log messages
- **Configurable Levels**: Minimum log levels configuration for each handler
- **DI Integration**: Seamless integration with Dependency Injection container
- **Extensible API**: Custom handler support through Handler interface
- **Resource Management**: Automatic cleanup and proper handler closure
- **Performance Optimized**: Efficient concurrent logging with minimal lock contention
- **Error Resilience**: Individual handler failures don't crash the logging system
- **Flexible Configuration**: Runtime handler management and dynamic reconfiguration
- **Color Support**: Advanced color formatting for console output
- **File Rotation**: Automatic file rotation based on size and time triggers
- **Context Logging**: Structured metadata support with context
- **OpenTelemetry Integration**: Trace context support for distributed tracing
- **Structured Logging**: Integration with structured logging standards
- **Memory Optimization**: Optimized for high-throughput logging scenarios

### Technical Details
- Initial release as standalone module `go.fork.vn/log`
- Repository located at `github.com/Fork/log`
- Built with Go 1.23.9
- Full test coverage (100% pass rate - 47 tests) and documentation included
- Thread-safe logging manager with minimal lock contention
- Memory leak prevention with proper resource management
- Easy mock regeneration with testing utilities

### Dependencies
- `go.fork.vn/di`: Dependency injection integration

[Unreleased]: github.com/go-fork/log/compare/v0.1.3...HEAD
[v0.1.3]: github.com/go-fork/log/compare/v0.1.2...v0.1.3
[v0.1.2]: github.com/go-fork/log/compare/v0.1.1...v0.1.2
[v0.1.1]: github.com/go-fork/log/compare/v0.1.0...v0.1.1
[v0.1.0]: github.com/go-fork/log/releases/tag/v0.1.0
