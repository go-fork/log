# Log Package Overview

## Architecture

The log package is designed around Fork Framework's core principles of dependency injection, service providers, and interface-based design. It provides a unified logging system that integrates seamlessly with the framework's ecosystem while maintaining high performance and flexibility.

### Core Components

#### Logger Interface
The central `log.Logger` interface provides all logging operations:
- Structured logging with key-value pairs
- Multiple log levels (Debug, Info, Warn, Error, Fatal)
- Contextual logging with additional metadata
- Performance-optimized logging methods

#### Service Provider
The `ServiceProvider` registers the logger as a singleton in the DI container, making it available throughout the application via `app.MustMake("logger")`.

#### Manager Implementation
The manager provides the concrete implementation of the Logger interface, handling:
- Log formatting and output
- Level filtering
- Context management
- Performance optimization

## Design Principles

### 1. Framework Integration First
The log package is designed specifically for Fork Framework applications. All examples and patterns assume usage within the framework's DI container system.

### 2. Structured Logging
All logging methods support structured logging with key-value pairs for better searchability and analysis:

```go
logger.Info("User logged in", "user_id", 123, "ip", "192.168.1.1")
```

### 3. Performance Oriented
The logger is optimized for high-performance applications:
- Minimal allocations during logging
- Efficient string formatting
- Level-based filtering to avoid unnecessary work

### 4. Contextual Information
Support for adding context to log entries:
- Request IDs for tracing
- User information
- Application state

## Logging Levels

### Debug Level
Use for detailed diagnostic information, typically only enabled during development:

```go
logger.Debug("Processing user request", "request_id", reqID, "user_id", userID)
```

### Info Level
Use for general informational messages about application flow:

```go
logger.Info("User registration completed", "user_id", userID, "email", email)
```

### Warn Level
Use for potentially harmful situations that don't prevent operation:

```go
logger.Warn("High memory usage detected", "usage_percent", 85)
```

### Error Level
Use for error events that might still allow the application to continue:

```go
logger.Error("Database connection failed", "error", err, "retry_count", retries)
```

### Fatal Level
Use for very severe errors that will cause the application to terminate:

```go
logger.Fatal("Cannot start server", "error", err, "port", port)
```

## Framework Integration Patterns

### Service Provider Registration
```go
// Automatic registration in Fork Framework
config := map[string]interface{}{
	    "name": "myapp",
	    "path": "./configs",
}
app := app.New(config)
// Logger provider is auto-registered

// Access via DI container
logger := app.MustMake("logger").(log.Logger)
```

### Dependency Injection in Services
```go
type UserService struct {
    logger log.Logger
}

func NewUserService(app app.Application) *UserService {
    return &UserService{
        logger: app.MustMake("logger").(log.Logger),
    }
}

func (s *UserService) CreateUser(user User) error {
    s.logger.Info("Creating new user", "email", user.Email)
    
    if err := s.validateUser(user); err != nil {
        s.logger.Error("User validation failed", "error", err, "email", user.Email)
        return err
    }
    
    s.logger.Info("User created successfully", "user_id", user.ID)
    return nil
}
```

### Service Provider Logging
```go
func (p *DatabaseServiceProvider) Register(container di.Container) error {
    logger := container.MustMake("logger").(log.Logger)
    logger.Debug("Registering database service provider")
    
    return container.Singleton("database", func(container di.Container) (interface{}, error) {
        logger.Info("Creating database connection")
        // ...create database connection
        return db, nil
    })
}
```

## Structured Logging

### Key-Value Pairs
The logger supports structured logging with key-value pairs:

```go
logger.Info("HTTP request completed",
    "method", "POST",
    "path", "/api/users",
    "status", 201,
    "duration", "45ms",
    "user_id", 123,
)
```

### Consistent Field Names
Use consistent field names across your application:

```go
// Good - consistent naming
logger.Info("User action", "user_id", 123, "action", "login")
logger.Error("User error", "user_id", 123, "error", err)

// Avoid - inconsistent naming
logger.Info("User action", "userID", 123, "action", "login")
logger.Error("User error", "user", 123, "err", err)
```

### Nested Context
For complex data structures, use nested context:

```go
logger.Info("Order processed",
    "order_id", order.ID,
    "customer_id", order.CustomerID,
    "items_count", len(order.Items),
    "total_amount", order.Total,
    "shipping_address", order.ShippingAddress.String(),
)
```

## Performance Considerations

### Level-Based Filtering
The logger performs level-based filtering to avoid unnecessary work:

```go
// Only evaluated if Debug level is enabled
logger.Debug("Expensive operation", "result", expensiveComputation())

// Better - check level first for expensive operations
if logger.IsDebugEnabled() {
    logger.Debug("Expensive operation", "result", expensiveComputation())
}
```

### String Formatting
Avoid expensive string operations in log arguments:

```go
// Good - let the logger handle formatting
logger.Info("Processing items", "count", len(items), "type", itemType)

// Avoid - expensive string operations
logger.Info(fmt.Sprintf("Processing %d items of type %s", len(items), itemType))
```

### Batch Logging
For high-frequency logging, consider batching:

```go
type BatchLogger struct {
    logger log.Logger
    buffer []LogEntry
    mu     sync.Mutex
}

func (b *BatchLogger) flush() {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    for _, entry := range b.buffer {
        b.logger.Log(entry.Level, entry.Message, entry.Fields...)
    }
    b.buffer = b.buffer[:0]
}
```

## Error Handling

### Logging Errors
Always include error information in log entries:

```go
if err := userService.CreateUser(user); err != nil {
    logger.Error("Failed to create user",
        "error", err,
        "user_email", user.Email,
        "validation_errors", user.ValidationErrors(),
    )
    return err
}
```

### Error Context
Provide context about where and why the error occurred:

```go
func (s *UserService) processPayment(userID int, amount float64) error {
    logger.Info("Processing payment", "user_id", userID, "amount", amount)
    
    if err := s.paymentGateway.Charge(userID, amount); err != nil {
        logger.Error("Payment processing failed",
            "error", err,
            "user_id", userID,
            "amount", amount,
            "gateway", s.paymentGateway.Name(),
            "operation", "charge",
        )
        return fmt.Errorf("payment failed: %w", err)
    }
    
    logger.Info("Payment processed successfully", "user_id", userID, "amount", amount)
    return nil
}
```

## Testing Strategy

### Unit Testing
Mock the log.Logger interface for unit tests:

```go
mockLogger := mocks.NewMockLogger(t)
mockLogger.On("Info", "User created", "user_id", 123).Return()
```

### Integration Testing
Use real logger with test-specific configuration:

```go
func TestUserService_Integration(t *testing.T) {
    app := core.NewApplication()
    logger := app.MustMake("logger").(log.Logger)
    
    // Set test log level
    logger.SetLevel(log.DebugLevel)
    
    service := NewUserService(app)
    // ...test implementation
}
```

### Log Verification
Verify that expected log entries are created:

```go
func TestUserService_CreateUser_LogsSuccess(t *testing.T) {
    mockLogger := mocks.NewMockLogger(t)
    
    // Expect success log
    mockLogger.On("Info", "User created successfully", "user_id", 123).Return()
    
    service := &UserService{logger: mockLogger}
    service.CreateUser(user)
    
    mockLogger.AssertExpectations(t)
}
```

## Best Practices

### 1. Use Structured Logging
Always use key-value pairs for structured logging:

```go
// Good
logger.Info("User login", "user_id", 123, "ip", clientIP)

// Avoid
logger.Info(fmt.Sprintf("User %d logged in from %s", 123, clientIP))
```

### 2. Include Relevant Context
Include context that helps with debugging:

```go
logger.Error("Database query failed",
    "error", err,
    "query", query,
    "params", params,
    "duration", time.Since(start),
    "connection_id", connID,
)
```

### 3. Use Appropriate Log Levels
Choose the right log level for each message:

```go
logger.Debug("Entering function", "user_id", userID)  // Development only
logger.Info("User action completed", "action", action) // Normal operations
logger.Warn("Rate limit approaching", "current", current, "limit", limit) // Warnings
logger.Error("Operation failed", "error", err) // Errors
logger.Fatal("Cannot start application", "error", err) // Fatal errors
```

### 4. Avoid Logging Sensitive Information
Never log sensitive data like passwords, tokens, or personal information:

```go
// Good
logger.Info("User authenticated", "user_id", user.ID)

// Never do this
logger.Info("User authenticated", "password", user.Password, "token", token)
```

### 5. Performance Optimization
Check log level for expensive operations:

```go
if logger.IsDebugEnabled() {
    logger.Debug("Request details", "body", string(requestBody))
}
```
