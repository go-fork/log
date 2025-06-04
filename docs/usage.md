# Log Package Usage Guide

This guide provides practical examples of using the log package within Fork Framework applications.

## Basic Framework Integration

### Application Setup

```go
package main

import (
    "go.fork.vn/core"
)

func main() {
    // Create Fork application
    app := core.NewApplication()
    
    // Logger service is automatically registered
    logger := app.MustMake("logger").(log.Logger)
    
    // Basic logging
    logger.Info("Application starting")
    logger.Debug("Debug mode enabled")
    
    app.Run()
}
```

## Service Integration Patterns

### In Service Providers

```go
package providers

import (
    "go.fork.vn/di"
    "go.fork.vn/log"
)

type DatabaseServiceProvider struct{}

func (p *DatabaseServiceProvider) Register(container di.Container) error {
    logger := container.MustMake("logger").(log.Logger)
    logger.Debug("Registering database service provider")
    
    return container.Singleton("database", func(container di.Container) (interface{}, error) {
        logger.Info("Creating database connection")
        
        configManager := container.MustMake("config").(config.Manager)
        
        var dbConfig DatabaseConfig
        if err := configManager.UnmarshalKey("database", &dbConfig); err != nil {
            logger.Error("Failed to load database config", "error", err)
            return nil, err
        }
        
        db, err := newDatabaseConnection(dbConfig)
        if err != nil {
            logger.Error("Failed to create database connection", 
                "error", err, 
                "host", dbConfig.Host,
                "database", dbConfig.Name,
            )
            return nil, err
        }
        
        logger.Info("Database connection established", 
            "host", dbConfig.Host,
            "database", dbConfig.Name,
        )
        return db, nil
    })
}

func (p *DatabaseServiceProvider) Boot(container di.Container) error {
    logger := container.MustMake("logger").(log.Logger)
    logger.Debug("Booting database service provider")
    return nil
}
```

### In Controllers

```go
package controllers

import (
    "net/http"
    "time"
    
    "go.fork.vn/core"
    "go.fork.vn/log"
)

type UserController struct {
    app    app.Application
    logger log.Logger
}

func NewUserController(app app.Application) *UserController {
    return &UserController{
        app:    app,
        logger: app.MustMake("logger").(log.Logger),
    }
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    // Create request-scoped logger
    requestLogger := c.logger.WithFields(
        "request_id", r.Header.Get("X-Request-ID"),
        "method", r.Method,
        "path", r.URL.Path,
        "user_agent", r.UserAgent(),
        "ip", r.RemoteAddr,
    )
    
    requestLogger.Info("Processing user creation request")
    
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        requestLogger.Error("Failed to decode request body", "error", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Validate user (don't log sensitive data)
    if err := user.Validate(); err != nil {
        requestLogger.Warn("User validation failed", 
            "validation_errors", err.Error(),
            "email", user.Email, // OK to log email
        )
        http.Error(w, "Validation failed", http.StatusBadRequest)
        return
    }
    
    userService := c.app.MustMake("user_service").(UserService)
    createdUser, err := userService.CreateUser(user)
    if err != nil {
        requestLogger.Error("Failed to create user", 
            "error", err,
            "email", user.Email,
        )
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    duration := time.Since(start)
    requestLogger.Info("User creation completed",
        "user_id", createdUser.ID,
        "email", createdUser.Email,
        "duration_ms", duration.Milliseconds(),
    )
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdUser)
}
```

### In Business Services

```go
package services

import (
    "errors"
    "time"
    
    "go.fork.vn/core"
    "go.fork.vn/log"
)

type UserService struct {
    app    app.Application
    logger log.Logger
    db     Database
}

func NewUserService(app app.Application) *UserService {
    return &UserService{
        app:    app,
        logger: app.MustMake("logger").(log.Logger),
        db:     app.MustMake("database").(Database),
    }
}

func (s *UserService) CreateUser(user User) (*User, error) {
    // Create service-scoped logger
    serviceLogger := s.logger.WithFields("service", "user_service", "operation", "create_user")
    
    serviceLogger.Debug("Starting user creation", "email", user.Email)
    
    // Check if user exists
    if s.userExists(user.Email) {
        serviceLogger.Warn("Attempted to create duplicate user", "email", user.Email)
        return nil, errors.New("user already exists")
    }
    
    // Hash password (don't log the password!)
    hashedPassword, err := s.hashPassword(user.Password)
    if err != nil {
        serviceLogger.Error("Failed to hash password", "error", err, "email", user.Email)
        return nil, err
    }
    
    user.Password = hashedPassword
    user.CreatedAt = time.Now()
    
    // Save to database
    if err := s.db.Create(&user); err != nil {
        serviceLogger.Error("Database error during user creation", 
            "error", err,
            "email", user.Email,
            "table", "users",
        )
        return nil, err
    }
    
    serviceLogger.Info("User created successfully", 
        "user_id", user.ID,
        "email", user.Email,
        "created_at", user.CreatedAt,
    )
    
    // Send welcome email
    go s.sendWelcomeEmail(user, serviceLogger)
    
    return &user, nil
}

func (s *UserService) sendWelcomeEmail(user User, parentLogger log.Logger) {
    emailLogger := parentLogger.WithFields("operation", "send_welcome_email", "user_id", user.ID)
    
    emailLogger.Debug("Sending welcome email")
    
    emailService := s.app.MustMake("email_service").(EmailService)
    if err := emailService.SendWelcomeEmail(user.Email, user.Name); err != nil {
        emailLogger.Error("Failed to send welcome email", 
            "error", err,
            "email", user.Email,
        )
        return
    }
    
    emailLogger.Info("Welcome email sent successfully", "email", user.Email)
}
```

## Configuration Integration

### Logger Configuration

```yaml
# configs/app.yaml
logging:
  level: "info"
  format: "json"
  output: "stdout"
  
app:
  debug: false
  
# Override in development
# configs/app.development.yaml
logging:
  level: "debug"
  format: "text"
  
app:
  debug: true
```

### Bootstrap Configuration

```go
// bootstrap/app.go
package bootstrap

import (
    "os"
    "go.fork.vn/core"
    "go.fork.vn/log"
)

func setupLogging(app app.Application) {
    configManager := app.MustMake("config").(config.Manager)
    logger := app.MustMake("logger").(log.Logger)
    
    // Set log level from config
    if levelStr, exists := configManager.GetString("logging.level"); exists {
        if level, err := log.ParseLevel(levelStr); err == nil {
            logger.SetLevel(level)
        }
    }
    
    // Set output destination
    if output, exists := configManager.GetString("logging.output"); exists {
        switch output {
        case "stdout":
            logger.SetOutput(os.Stdout)
        case "stderr":
            logger.SetOutput(os.Stderr)
        default:
            // File output
            if file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
                logger.SetOutput(file)
            }
        }
    }
    
    // Enable debug logging if app is in debug mode
    if debug, exists := configManager.GetBool("app.debug"); exists && debug {
        logger.SetLevel(log.DebugLevel)
        logger.Debug("Debug logging enabled")
    }
}
```

## Structured Logging Patterns

### HTTP Request Logging

```go
func (c *BaseController) logRequest(r *http.Request, statusCode int, duration time.Duration) {
    c.logger.Info("HTTP request completed",
        "method", r.Method,
        "path", r.URL.Path,
        "query", r.URL.RawQuery,
        "status", statusCode,
        "duration_ms", duration.Milliseconds(),
        "user_agent", r.UserAgent(),
        "ip", c.getClientIP(r),
        "content_length", r.ContentLength,
    )
}
```

### Database Operation Logging

```go
func (db *Database) logQuery(query string, args []interface{}, duration time.Duration, err error) {
    fields := []interface{}{
        "query", query,
        "duration_ms", duration.Milliseconds(),
        "args_count", len(args),
    }
    
    if err != nil {
        fields = append(fields, "error", err)
        db.logger.Error("Database query failed", fields...)
    } else {
        db.logger.Debug("Database query executed", fields...)
    }
}
```

### Business Logic Logging

```go
func (s *OrderService) ProcessOrder(order Order) error {
    orderLogger := s.logger.WithFields(
        "service", "order_service",
        "operation", "process_order",
        "order_id", order.ID,
        "customer_id", order.CustomerID,
    )
    
    orderLogger.Info("Processing order", 
        "total_amount", order.TotalAmount,
        "item_count", len(order.Items),
    )
    
    // Validate inventory
    for _, item := range order.Items {
        if !s.inventory.HasStock(item.ProductID, item.Quantity) {
            orderLogger.Warn("Insufficient inventory", 
                "product_id", item.ProductID,
                "requested", item.Quantity,
                "available", s.inventory.GetStock(item.ProductID),
            )
            return errors.New("insufficient inventory")
        }
    }
    
    // Process payment
    if err := s.processPayment(order, orderLogger); err != nil {
        return err
    }
    
    // Update inventory
    if err := s.updateInventory(order, orderLogger); err != nil {
        return err
    }
    
    orderLogger.Info("Order processed successfully")
    return nil
}
```

## Error Logging Patterns

### Error Context

```go
func (s *PaymentService) ProcessPayment(amount float64, customerID int) error {
    paymentLogger := s.logger.WithFields(
        "service", "payment_service",
        "customer_id", customerID,
        "amount", amount,
    )
    
    paymentLogger.Info("Processing payment")
    
    // Validate customer
    customer, err := s.getCustomer(customerID)
    if err != nil {
        paymentLogger.Error("Failed to retrieve customer", 
            "error", err,
            "customer_id", customerID,
        )
        return fmt.Errorf("customer validation failed: %w", err)
    }
    
    // Process with payment gateway
    result, err := s.gateway.Charge(customer.PaymentToken, amount)
    if err != nil {
        paymentLogger.Error("Payment gateway error", 
            "error", err,
            "gateway", s.gateway.Name(),
            "transaction_id", result.TransactionID,
            "gateway_code", result.ErrorCode,
        )
        return fmt.Errorf("payment failed: %w", err)
    }
    
    paymentLogger.Info("Payment processed successfully", 
        "transaction_id", result.TransactionID,
        "gateway_fee", result.Fee,
    )
    
    return nil
}
```

### Error Recovery

```go
func (s *EmailService) SendEmail(to, subject, body string) error {
    emailLogger := s.logger.WithFields(
        "service", "email_service",
        "to", to,
        "subject", subject,
    )
    
    emailLogger.Debug("Sending email")
    
    // Try primary SMTP server
    err := s.primarySMTP.Send(to, subject, body)
    if err != nil {
        emailLogger.Warn("Primary SMTP failed, trying backup", 
            "error", err,
            "smtp_host", s.primarySMTP.Host,
        )
        
        // Try backup SMTP server
        err = s.backupSMTP.Send(to, subject, body)
        if err != nil {
            emailLogger.Error("All SMTP servers failed", 
                "primary_error", err,
                "backup_error", err,
                "primary_host", s.primarySMTP.Host,
                "backup_host", s.backupSMTP.Host,
            )
            return fmt.Errorf("email delivery failed: %w", err)
        }
        
        emailLogger.Info("Email sent via backup SMTP", "smtp_host", s.backupSMTP.Host)
    } else {
        emailLogger.Info("Email sent successfully", "smtp_host", s.primarySMTP.Host)
    }
    
    return nil
}
```

## Performance Optimization

### Level-Based Logging

```go
func (s *DataProcessor) ProcessLargeDataset(data []DataPoint) {
    // Only log debug info if debug is enabled
    if s.logger.IsDebugEnabled() {
        s.logger.Debug("Processing dataset", 
            "size", len(data),
            "sample_data", data[:min(5, len(data))], // Only first 5 items
        )
    }
    
    s.logger.Info("Starting data processing", "size", len(data))
    
    for i, point := range data {
        // Expensive debug logging
        if s.logger.IsDebugEnabled() && i%1000 == 0 {
            s.logger.Debug("Processing progress", 
                "processed", i,
                "total", len(data),
                "percent", float64(i)/float64(len(data))*100,
            )
        }
        
        if err := s.processDataPoint(point); err != nil {
            s.logger.Error("Failed to process data point", 
                "error", err,
                "index", i,
                "point_id", point.ID,
            )
        }
    }
    
    s.logger.Info("Data processing completed", "processed", len(data))
}
```

### Logger Reuse

```go
type RequestHandler struct {
    baseLogger log.Logger
}

func (h *RequestHandler) HandleRequest(r *http.Request) {
    // Create request-scoped logger once
    requestLogger := h.baseLogger.WithFields(
        "request_id", r.Header.Get("X-Request-ID"),
        "method", r.Method,
        "path", r.URL.Path,
    )
    
    // Pass logger down to other functions
    h.authenticateUser(r, requestLogger)
    h.processRequest(r, requestLogger)
    h.sendResponse(r, requestLogger)
}

func (h *RequestHandler) authenticateUser(r *http.Request, logger log.Logger) {
    // Use the same logger with additional context
    authLogger := logger.WithFields("operation", "authenticate")
    
    authLogger.Debug("Authenticating user")
    // ... authentication logic
}
```

## Testing with Logging

### Mock-Based Testing

```go
package services_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "go.fork.vn/log/mocks"
)

func TestUserService_CreateUser_Success(t *testing.T) {
    mockLogger := mocks.NewMockLogger(t)
    
    // Expect specific log calls
    mockLogger.On("Debug", "Starting user creation", "email", "test@example.com").Return()
    mockLogger.On("Info", "User created successfully", 
        "user_id", 123,
        "email", "test@example.com",
        mock.AnythingOfType("time.Time"),
    ).Return()
    
    service := &UserService{logger: mockLogger}
    user := User{Email: "test@example.com"}
    
    result, err := service.CreateUser(user)
    
    assert.NoError(t, err)
    assert.Equal(t, 123, result.ID)
    mockLogger.AssertExpectations(t)
}

func TestUserService_CreateUser_Error(t *testing.T) {
    mockLogger := mocks.NewMockLogger(t)
    
    // Expect error log
    mockLogger.On("Error", "Database error during user creation",
        "error", mock.AnythingOfType("*errors.errorString"),
        "email", "test@example.com",
        "table", "users",
    ).Return()
    
    service := &UserService{logger: mockLogger}
    user := User{Email: "test@example.com"}
    
    _, err := service.CreateUser(user)
    
    assert.Error(t, err)
    mockLogger.AssertExpectations(t)
}
```

### Integration Testing

```go
func TestUserController_Integration(t *testing.T) {
    config := map[string]interface{}{
	    "name": "myapp",
	    "path": "./configs",
	}
	app := appMocks.New(config)
    
    // Configure test logging
    logger := app.MustMake("logger").(log.Logger)
    logger.SetLevel(log.DebugLevel)
    
    controller := NewUserController(app)
    
    // Create test request
    user := User{Email: "test@example.com", Name: "Test User"}
    body, _ := json.Marshal(user)
    req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Request-ID", "test-req-123")
    
    w := httptest.NewRecorder()
    controller.CreateUser(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    
    // Verify response
    var result User
    json.NewDecoder(w.Body).Decode(&result)
    assert.Equal(t, "test@example.com", result.Email)
}
```

## Common Patterns

### Middleware Logging

```go
func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Generate request ID if not present
            requestID := r.Header.Get("X-Request-ID")
            if requestID == "" {
                requestID = generateRequestID()
                r.Header.Set("X-Request-ID", requestID)
            }
            
            // Create request logger
            // Wrap response writer to capture status code
            wrapped := &responseWriter{ResponseWriter: w}
            
                "status", wrapped.statusCode,
                "duration_ms", duration.Milliseconds(),
    jobLogger := w.logger.WithFields(
        "worker_id", w.ID,
        "job_id", job.ID,
        "job_type", job.Type,
    )
    
    jobLogger.Info("Starting job processing")
    
    defer func() {
        if r := recover(); r != nil {
            jobLogger.Error("Job panicked", 
                "panic", r,
                "stack", debug.Stack(),
            )
        }
    }()
    
    start := time.Now()
    
    if err := w.executeJob(job, jobLogger); err != nil {
        jobLogger.Error("Job failed", 
            "error", err,
            "duration_ms", time.Since(start).Milliseconds(),
        )
        w.requeueJob(job, jobLogger)
        return
    }
    
    jobLogger.Info("Job completed successfully", 
        "duration_ms", time.Since(start).Milliseconds(),
    )
}
```

### Configuration-Based Logger Factory

```go
type LoggerFactory struct {
    config config.Manager
}

func NewLoggerFactory(app app.Application) *LoggerFactory {
    return &LoggerFactory{
        config: app.MustMake("config").(config.Manager),
    }
}

func (f *LoggerFactory) CreateServiceLogger(serviceName string) log.Logger {
    baseLogger := f.app.MustMake("logger").(log.Logger)
    
    // Add service-specific context
    serviceLogger := baseLogger.WithFields("service", serviceName)
    
    // Service-specific log level
    if level, exists := f.config.GetString(fmt.Sprintf("logging.services.%s.level", serviceName)); exists {
        if logLevel, err := log.ParseLevel(level); err == nil {
            serviceLogger.SetLevel(logLevel)
        }
    }
    
    return serviceLogger
}
```
