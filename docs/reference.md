# Log Package API Reference

## Interfaces

### Logger

The `Logger` interface provides all logging functionality.

```go
type Logger interface {
    // Basic logging methods
    Debug(msg string, keysAndValues ...interface{})
    Info(msg string, keysAndValues ...interface{})
    Warn(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
    Fatal(msg string, keysAndValues ...interface{})
    
    // Level checking methods
    IsDebugEnabled() bool
    IsInfoEnabled() bool
    IsWarnEnabled() bool
    IsErrorEnabled() bool
    
    // Configuration methods
    SetLevel(level Level)
    GetLevel() Level
    
    // Context methods
    WithContext(ctx context.Context) Logger
    WithFields(keysAndValues ...interface{}) Logger
    
    // Output methods
    SetOutput(output io.Writer)
    GetOutput() io.Writer
}
```

### Level

Log level enumeration for controlling log output.

```go
type Level int

const (
    DebugLevel Level = iota
    InfoLevel
    WarnLevel
    ErrorLevel
    FatalLevel
)
```

## Methods Documentation

### Basic Logging Methods

#### Debug(msg string, keysAndValues ...interface{})
Logs a debug message with optional key-value pairs.

**Parameters:**
- `msg` - Log message
- `keysAndValues` - Optional key-value pairs for structured logging

**Example:**
```go
logger.Debug("Processing request", "request_id", "req-123", "user_id", 456)
```

#### Info(msg string, keysAndValues ...interface{})
Logs an informational message with optional key-value pairs.

**Parameters:**
- `msg` - Log message
- `keysAndValues` - Optional key-value pairs for structured logging

**Example:**
```go
logger.Info("User created successfully", "user_id", 123, "email", "user@example.com")
```

#### Warn(msg string, keysAndValues ...interface{})
Logs a warning message with optional key-value pairs.

**Parameters:**
- `msg` - Log message
- `keysAndValues` - Optional key-value pairs for structured logging

**Example:**
```go
logger.Warn("High memory usage", "usage_percent", 85, "threshold", 80)
```

#### Error(msg string, keysAndValues ...interface{})
Logs an error message with optional key-value pairs.

**Parameters:**
- `msg` - Log message
- `keysAndValues` - Optional key-value pairs for structured logging

**Example:**
```go
logger.Error("Database connection failed", "error", err, "host", "db.example.com")
```

#### Fatal(msg string, keysAndValues ...interface{})
Logs a fatal error message and terminates the application.

**Parameters:**
- `msg` - Log message
- `keysAndValues` - Optional key-value pairs for structured logging

**Example:**
```go
logger.Fatal("Cannot bind to port", "error", err, "port", 8080)
```

### Level Checking Methods

#### IsDebugEnabled() bool
Checks if debug level logging is enabled.

**Returns:**
- `bool` - True if debug logging is enabled

**Example:**
```go
if logger.IsDebugEnabled() {
    logger.Debug("Expensive debug info", "data", expensiveOperation())
}
```

#### IsInfoEnabled() bool
Checks if info level logging is enabled.

**Returns:**
- `bool` - True if info logging is enabled

#### IsWarnEnabled() bool
Checks if warn level logging is enabled.

**Returns:**
- `bool` - True if warn logging is enabled

#### IsErrorEnabled() bool
Checks if error level logging is enabled.

**Returns:**
- `bool` - True if error logging is enabled

### Configuration Methods

#### SetLevel(level Level)
Sets the minimum log level for output.

**Parameters:**
- `level` - Minimum log level (DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel)

**Example:**
```go
logger.SetLevel(log.InfoLevel) // Only Info, Warn, Error, Fatal will be logged
```

#### GetLevel() Level
Gets the current minimum log level.

**Returns:**
- `Level` - Current log level

**Example:**
```go
currentLevel := logger.GetLevel()
if currentLevel == log.DebugLevel {
    // Debug logging is enabled
}
```

### Context Methods

#### WithContext(ctx context.Context) Logger
Creates a new logger with the provided context.

**Parameters:**
- `ctx` - Context to attach to the logger

**Returns:**
- `Logger` - New logger instance with context

**Example:**
```go
contextLogger := logger.WithContext(ctx)
contextLogger.Info("Request processed", "user_id", 123)
```

#### WithFields(keysAndValues ...interface{}) Logger
Creates a new logger with pre-populated fields.

**Parameters:**
- `keysAndValues` - Key-value pairs to attach to all log entries

**Returns:**
- `Logger` - New logger instance with attached fields

**Example:**
```go
userLogger := logger.WithFields("user_id", 123, "session_id", "sess-456")
userLogger.Info("User action") // Will include user_id and session_id
userLogger.Error("User error", "error", err) // Will include user_id, session_id, and error
```

### Output Methods

#### SetOutput(output io.Writer)
Sets the output destination for log entries.

**Parameters:**
- `output` - Writer to send log output to

**Example:**
```go
// Log to file
file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    log.Fatal(err)
}
logger.SetOutput(file)

// Log to stdout
logger.SetOutput(os.Stdout)
```

#### GetOutput() io.Writer
Gets the current output destination.

**Returns:**
- `io.Writer` - Current output writer

**Example:**
```go
currentOutput := logger.GetOutput()
if currentOutput == os.Stdout {
    // Logging to stdout
}
```

## Constants and Types

### Log Levels

```go
const (
    DebugLevel Level = iota // Most verbose
    InfoLevel              // General information
    WarnLevel              // Warning messages
    ErrorLevel             // Error messages
    FatalLevel             // Fatal errors (causes exit)
)
```

### Level Methods

#### String() string
Returns the string representation of the log level.

**Example:**
```go
level := log.InfoLevel
fmt.Println(level.String()) // Output: "INFO"
```

#### ParseLevel(s string) (Level, error)
Parses a string into a log level.

**Parameters:**
- `s` - String representation of level ("debug", "info", "warn", "error", "fatal")

**Returns:**
- `Level` - Parsed log level
- `error` - Error if string cannot be parsed

**Example:**
```go
level, err := log.ParseLevel("info")
if err != nil {
    log.Fatal("Invalid log level")
}
logger.SetLevel(level)
```

## Usage with Fork Framework

### Service Provider Integration

```go
// The logger service is automatically available
config := map[string]interface{}{
    "name": "myapp",
    "path": "./configs",
}
app := app.New(config)
logger := app.MustMake("logger").(log.Logger)
```

### Configuration

The logger can be configured through the config system:

```yaml
logging:
  level: "info"
  output: "stdout"  # or file path
  format: "json"    # or "text"
```

```go
// In service provider or bootstrap
configManager := app.MustMake("config").(config.Manager)
logger := app.MustMake("logger").(log.Logger)

if level, exists := configManager.GetString("logging.level"); exists {
    if logLevel, err := log.ParseLevel(level); err == nil {
        logger.SetLevel(logLevel)
    }
}
```

### Testing

```go
import "go.fork.vn/log/mocks"

func TestMyService(t *testing.T) {
    mockLogger := mocks.NewMockLogger(t)
    mockLogger.On("Info", "User created", "user_id", 123).Return()
    
    service := &MyService{logger: mockLogger}
    service.CreateUser(user)
    
    mockLogger.AssertExpectations(t)
}
```

## Error Handling

The logger interface methods do not return errors. Internal errors are handled as follows:

- **Output errors**: Logged to stderr if possible
- **Formatting errors**: Invalid key-value pairs are logged as best effort
- **Level errors**: Invalid levels default to InfoLevel

## Performance Notes

### Allocation Optimization
The logger is optimized to minimize allocations:

```go
// Efficient - no string concatenation
logger.Info("User action", "user_id", 123, "action", "login")

// Less efficient - creates temporary strings
logger.Info(fmt.Sprintf("User %d performed %s", 123, "login"))
```

### Level Checking
Use level checking for expensive operations:

```go
// Good - only executes if debug is enabled
if logger.IsDebugEnabled() {
    logger.Debug("Request body", "body", string(requestBody))
}

// Less efficient - always converts body to string
logger.Debug("Request body", "body", string(requestBody))
```

### Field Reuse
Reuse loggers with common fields:

```go
// Create once, use many times
userLogger := logger.WithFields("user_id", 123)

userLogger.Info("Login attempt")
userLogger.Info("Login successful")
userLogger.Error("Login failed", "error", err)
```

## Thread Safety

All logger methods are thread-safe and can be called concurrently from multiple goroutines without synchronization.

## Examples

### Basic Usage
```go
logger := app.MustMake("logger").(log.Logger)

logger.Debug("Starting request processing")
logger.Info("User authenticated", "user_id", 123)
logger.Warn("Rate limit approaching", "current", 95, "limit", 100)
logger.Error("Database error", "error", err, "query", "SELECT * FROM users")
```

### Structured Logging
```go
logger.Info("HTTP request completed",
    "method", "POST",
    "path", "/api/users",
    "status", 201,
    "duration_ms", 45,
    "user_id", 123,
    "ip", "192.168.1.1",
)
```

### Context Logging
```go
requestLogger := logger.WithFields("request_id", requestID, "user_id", userID)

requestLogger.Info("Processing payment")
requestLogger.Error("Payment failed", "error", err, "amount", 100.00)
```

### Conditional Logging
```go
if logger.IsDebugEnabled() {
    logger.Debug("Complex data structure", "data", complexObject.String())
}
```
