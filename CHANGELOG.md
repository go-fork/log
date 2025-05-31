# Changelog

## [Unreleased]

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
- Repository located at `github.com/go-fork/log`
- Built with Go 1.23.9
- Full test coverage (100% pass rate - 47 tests) and documentation included
- Thread-safe logging manager with minimal lock contention
- Memory leak prevention with proper resource management
- Easy mock regeneration with testing utilities

### Dependencies
- `go.fork.vn/di`: Dependency injection integration

[Unreleased]: https://github.com/go-fork/log/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/go-fork/log/releases/tag/v0.1.0
