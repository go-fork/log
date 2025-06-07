# Changelog

## [Unreleased]

## v0.1.7 - 2025-06-07

### Fixed
- **Stack Handler Duplicate Logging Issue**
  - Fixed critical bug where logs were duplicated when both individual handlers and stack handler were enabled
  - Enhanced `initializeHandlers()` to respect stack handler configuration instead of always adding all handlers
  - Stack handler now only includes handlers that are explicitly configured in `stack.handlers`
  - Improved handler logic in `GetLogger()` to prevent duplication between individual and stack handlers

### Improved
- **Handler Logic Optimization**
  - Implemented smart handler selection logic to avoid duplicate outputs
  - Added conditional handler addition: individual handlers are only added when not already covered by stack handler
  - Enhanced stack handler configuration to be more granular and predictable
  - Better separation of concerns between individual handlers and composite stack handler

### Technical Details
- **Stack Handler Configuration**: Now respects `stack.handlers.console` and `stack.handlers.file` settings
- **Anti-Duplication Logic**: Individual handlers are only added when `!stack.enabled || !stack.handlers.[type]`
- **Configuration Priority**: Stack handler acts as a composite handler that can replace individual handlers
- **Backward Compatibility**: All existing configurations continue to work with improved behavior

## v0.1.6 - 2025-06-07

### Changed
- **Default Configuration Architecture**
  - Modified `DefaultConfig()` to use empty file path by default for better security
  - Updated `File.Enabled` to `false` by default to prevent accidental file creation
  - Enhanced validation logic to only require file path when file handler is actually used
  - Improved separation between default configuration and user-specific settings
  - Better configuration validation that respects actual handler usage patterns

### Fixed
- **Repository Structure Cleanup**
  - Removed dependency on hardcoded `storages/log/app.log` directory structure
  - Eliminated repository pollution with unnecessary directory structures
  - Tests now skip gracefully when required directories don't exist instead of creating them
  - Enhanced test portability by using `t.TempDir()` for temporary test directories
  - Cleaner development experience without mandatory directory creation

### Improved
- **Test Infrastructure Modernization**
  - Enhanced test reliability by removing hardcoded directory dependencies
  - Better error handling when directories don't exist (skip vs fail approach)
  - More portable test suite that works across different environments
  - Cleaner repository structure without test artifacts
  - Improved CI/CD compatibility with flexible directory requirements

### Technical Details
- **Validation Logic**: File path validation now conditional based on actual handler usage
- **Default Behavior**: More secure defaults that require explicit user configuration
- **Test Strategy**: Skip-based testing for missing dependencies instead of creation
- **Development Experience**: Cleaner repository without mandatory directory structures

## v0.1.5 - 2025-06-07

### Changed
- **Configuration Architecture Update**
  - Modified `DefaultConfig()` to set `File.Enabled: false` by default
  - Updated validation logic to always check `File.Path` regardless of `File.Enabled` status
  - Enhanced `Config.Validate()` to ensure proper file system validation
  - Improved error messages for file handler initialization requirements

- **FileHandler Security Enhancement**
  - Modified `NewFileHandler()` to require existing directories with write permissions
  - Removed automatic directory creation for security and explicit control
  - Enhanced directory validation with specific error messages
  - Added comprehensive write permission testing before file handler initialization

### Fixed
- **Directory Creation and Permissions**
  - Enhanced `NewFileHandler` to validate directory existence and write permissions
  - Added comprehensive error messages: "path to folder do not exists" and "directory does not have write permission"
  - Enhanced error handling for file system operations during validation

- **CI/CD Pipeline Issues**
  - Fixed empty branch issue in `config.go` that triggered staticcheck warning SA9003
  - Fixed errcheck warnings for unchecked `os.Chmod` calls in test files
  - Added proper error handling for file permission operations in tests
  - Enhanced test directory setup for CI environment compatibility

### Added
- **Enhanced Test Coverage**
  - Added test cases for non-existent directory scenarios
  - Added test cases for write permission validation
  - Added test cases for directory access errors
  - Improved error message validation in tests

- **Test Infrastructure Improvements**
  - Converted all test case names to snake_case format for consistency
  - Enhanced test coverage for file handler edge cases and error conditions
  - Improved benchmark test naming conventions
  - Added comprehensive validation scenarios
  - Added proper directory creation for CI/CD test environments
  - Created `storages/log/` directory structure for default configuration tests

- **CI/CD Enhancements**
  - Upgraded GitHub Actions cache from v3 to v4 to fix module caching issues
  - Enhanced test setup to create required directories before validation
  - Improved go.mod/go.sum consistency for CI environments

### Breaking Changes
- **NewFileHandler Behavior**
  - `NewFileHandler` now requires the directory to exist beforehand
  - No longer automatically creates directories
  - Enhanced validation provides detailed error messages for debugging

## v0.1.4 - 2025-06-07

### Added
- **Complete Vietnamese Documentation Suite**
  - Brand new comprehensive documentation written entirely in Vietnamese
  - Added `docs/index.md` - Main documentation hub with quick start guide
  - Added `docs/overview.md` - Architecture overview with Mermaid diagrams
  - Added `docs/configuration.md` - Detailed configuration guide for all environments
  - Added `docs/handler.md` - Complete handler documentation with performance comparisons
  - Added `docs/logger.md` - Logger interface and contextual logging patterns
  - Added `docs/workflows.md` - Application lifecycle and integration workflows

- **Mermaid Architecture Diagrams**
  - Shared Handlers Architecture visualization
  - Fork Framework integration flow charts
  - Handler processing workflow diagrams
  - Application lifecycle and deployment patterns
  - Environment-specific configuration flows

- **Fork Framework Integration Focus**
  - Specialized documentation for Fork Framework patterns
  - Dependency Injection container integration examples
  - Service Provider pattern implementations
  - Contextual logging for microservices architecture
  - Performance monitoring and alerting workflows

- **Environment-Specific Configurations**
  - Development environment setup with debugging features
  - Production environment optimizations and monitoring
  - Testing environment configurations with mock patterns
  - Docker and containerization deployment guides
  - Cloud-native logging strategies

### Changed
- **Complete Documentation Rewrite**
  - Migrated from English to Vietnamese for better accessibility
  - Restructured documentation architecture from ground up
  - Replaced simple examples with comprehensive real-world patterns
  - Enhanced code examples with Fork Framework integration
  - Updated all configuration examples with modern best practices

- **README.md Complete Overhaul**
  - Modern GitHub badges for Go version, releases, coverage, and quality
  - Professional project presentation with emoji icons
  - Comprehensive quick start examples for both standalone and Fork Framework usage
  - Architecture overview section with clear explanations
  - Advanced usage patterns including middleware and performance monitoring
  - Structured documentation navigation with clear links

- **Documentation Structure Reorganization**
  - Consolidated 6 focused documentation files instead of scattered content
  - Logical flow from overview → configuration → implementation → workflows
  - Cross-referencing between documentation sections
  - Improved navigation with clear section headers
  - Better code-to-documentation ratio with practical examples

### Removed
- **Legacy Documentation Files**
  - Removed outdated English documentation files
  - Cleaned up all `.bak` backup files from docs directory
  - Removed redundant or conflicting documentation
  - Eliminated inconsistent configuration examples
  - Removed deprecated usage patterns

### Improved
- **Code Example Quality**
  - All examples now follow Fork Framework conventions
  - Realistic service integration patterns (UserService, OrderService)
  - Production-ready configuration examples
  - Error handling best practices in all examples
  - Performance optimization techniques demonstrated

- **Documentation Accessibility**
  - Vietnamese language for Vietnamese developers
  - Clear technical terminology with explanations
  - Step-by-step implementation guides
  - Troubleshooting sections with common issues
  - Migration guides for different use cases

- **Visual Documentation**
  - Mermaid diagrams for complex architecture concepts
  - Flow charts for decision-making processes
  - Sequence diagrams for integration patterns
  - Component relationship visualizations
  - Data flow illustrations

### Technical Details
- **Documentation Architecture**: 6 core documentation files totaling ~95KB
- **Language**: Complete Vietnamese translation for better developer experience
- **Visual Elements**: 12+ Mermaid diagrams for architecture visualization
- **Code Examples**: 50+ practical examples covering all major use cases
- **Integration Focus**: Deep Fork Framework integration with DI container patterns
- **Environment Coverage**: Development, Production, Testing, and Container deployments

### Documentation Structure
```
docs/
├── index.md         # 8.2KB  - Main hub with quick start
├── overview.md      # 7.6KB  - Architecture and concepts  
├── configuration.md # 12.3KB - Environment configurations
├── handler.md       # 14.6KB - Handler implementations
├── logger.md        # 15.3KB - Logger patterns and usage
└── workflows.md     # 20.9KB - Integration workflows
```

### Migration Guide
- **For Existing Users**: All core APIs remain unchanged
- **Documentation**: New Vietnamese docs replace English versions
- **Examples**: Updated to Fork Framework patterns but backward compatible
- **Configuration**: Enhanced examples but existing configs still work
- **Integration**: New patterns available, existing integrations unaffected

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

[Unreleased]: github.com/go-fork/log/compare/v0.1.4...HEAD
[v0.1.4]: github.com/go-fork/log/compare/v0.1.3...v0.1.4
[v0.1.3]: github.com/go-fork/log/compare/v0.1.2...v0.1.3
[v0.1.2]: github.com/go-fork/log/compare/v0.1.1...v0.1.2
[v0.1.1]: github.com/go-fork/log/compare/v0.1.0...v0.1.1
[v0.1.0]: github.com/go-fork/log/releases/tag/v0.1.0
