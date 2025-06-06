// Package log cung cấp hệ thống logging có cấu trúc và hiệu suất cao cho Fork Framework,
// được thiết kế với kiến trúc Shared Handlers và Contextual Loggers.
//
// # Tổng quan
//
// Package log triển khai một hệ thống logging toàn diện với kiến trúc Shared Handlers,
// nơi các handler được chia sẻ giữa nhiều logger instances để tối ưu hiệu suất và tài nguyên.
// Mỗi logger có context riêng (VD: UserService, OrderService) để dễ dàng theo dõi nguồn gốc log.
//
// # Tính năng chính
//
//   - Contextual Logging: Logger được tạo theo context để phân biệt nguồn gốc
//   - Shared Handlers: Handlers được chia sẻ giữa nhiều logger instances
//   - Multi-Level Support: Debug, Info, Warning, Error, Fatal levels
//   - Flexible Configuration: Cấu hình linh hoạt cho console, file, stack handlers
//   - Runtime Management: Quản lý handlers và loggers trong runtime
//   - High Performance: Tối ưu hóa cho ứng dụng hiệu suất cao
//   - Fork Framework Integration: Tích hợp sâu với DI container và service providers
//   - Thread-Safe: An toàn cho các ứng dụng concurrent
//   - Zero Allocation: Tối ưu allocation cho hot paths
//
// # Kiến trúc Shared Handlers
//
// Package sử dụng kiến trúc Shared Handlers với các thành phần chính:
//
//   - Manager: Quản lý tập trung các handlers và loggers
//   - Shared Handlers: Console, File, Stack handlers được chia sẻ
//   - Contextual Loggers: Mỗi logger có context riêng
//   - Runtime Management: Thêm/xóa handlers và loggers trong runtime
//
// # Sử dụng cơ bản
//
//	package main
//
//	import (
//	    "go.fork.vn/log"
//	    "go.fork.vn/log/handler"
//	)
//
//	func main() {
//	    // Tạo cấu hình
//	    config := &log.Config{
//	        Level: handler.InfoLevel,
//	        Console: log.ConsoleConfig{
//	            Enabled: true,
//	            Colored: true,
//	        },
//	        File: log.FileConfig{
//	            Enabled: true,
//	            Path:    "logs/app.log",
//	            MaxSize: 10 * 1024 * 1024, // 10MB
//	        },
//	    }
//
//	    // Khởi tạo manager
//	    manager := log.NewManager(config)
//	    defer manager.Close()
//
//	    // Lấy logger theo context
//	    userLogger := manager.GetLogger("UserService")
//	    orderLogger := manager.GetLogger("OrderService")
//
//	    // Structured logging
//	    userLogger.Info("User logged in", "user_id", 12345, "username", "john_doe")
//	    orderLogger.Warning("Low stock", "product_id", "ABC123", "stock", 2)
//	    userLogger.Error("Database error", "error", "connection timeout")
//	}
//
// # Tích hợp với Fork Framework
//
// Package log tích hợp mượt mà với Fork Framework thông qua Service Provider pattern:
//
//	// main.go
//	func main() {
//	    app := fork.NewApplication()
//
//	    // Đăng ký log service provider
//	    logConfig := log.DefaultConfig()
//	    app.RegisterProvider(providers.NewLogProvider(logConfig))
//
//	    // Đăng ký service khác
//	    app.RegisterProvider(providers.NewUserServiceProvider())
//
//	    app.Run()
//	}
//
// # Service Integration Pattern
//
// Trong Fork Framework, services có thể inject log manager thông qua DI container:
//
//	// services/user_service.go
//	type UserService struct {
//	    logger log.Logger
//	    db     *database.DB
//	}
//
//	func NewUserService(container *container.Container) *UserService {
//	    manager := container.Get("log").(log.Manager)
//
//	    return &UserService{
//	        logger: manager.GetLogger("UserService"),
//	        db:     container.Get("database").(*database.DB),
//	    }
//	}
//
//	func (s *UserService) CreateUser(user *User) error {
//	    s.logger.Info("Creating user", "username", user.Username)
//
//	    if err := s.db.Create(user); err != nil {
//	        s.logger.Error("Failed to create user",
//	            "username", user.Username,
//	            "error", err.Error(),
//	        )
//	        return err
//	    }
//
//	    s.logger.Info("User created successfully",
//	        "user_id", user.ID,
//	        "username", user.Username,
//	    )
//
//	    return nil
//	}
//
// # Handlers
//
// Package hỗ trợ 3 loại handler chính:
//
// ## Console Handler
//
// Xuất logs ra stdout/stderr với hỗ trợ màu sắc:
//
//	consoleConfig := log.ConsoleConfig{
//	    Enabled: true,
//	    Colored: true, // Màu sắc theo level
//	}
//
// ## File Handler
//
// Ghi logs vào file với rotation tự động:
//
//	fileConfig := log.FileConfig{
//	    Enabled: true,
//	    Path:    "/var/log/app/app.log",
//	    MaxSize: 100 * 1024 * 1024, // 100MB
//	}
//
// ## Stack Handler
//
// Kết hợp nhiều handlers để log đồng thời ra nhiều đích:
//
//	stackConfig := log.StackConfig{
//	    Enabled: true,
//	    Handlers: log.StackHandlers{
//	        Console: true,
//	        File:    true,
//	    },
//	}
//
// # Custom Handlers
//
// Có thể tạo custom handlers bằng cách implement interface Handler:
//
//	type DatabaseHandler struct {
//	    db *sql.DB
//	}
//
//	func (h *DatabaseHandler) Log(level handler.Level, message string, args ...interface{}) error {
//	    query := `INSERT INTO logs (level, message, data, created_at) VALUES (?, ?, ?, ?)`
//	    data, _ := json.Marshal(h.argsToMap(args))
//	    _, err := h.db.Exec(query, level.String(), message, string(data), time.Now())
//	    return err
//	}
//
//	func (h *DatabaseHandler) Close() error {
//	    return nil
//	}
//
// # Log Levels
//
// Package hỗ trợ 5 levels logging chuẩn:
//
//   - Debug (0): Chi tiết implementation, chỉ development
//   - Info (1): Thông tin chung về application flow
//   - Warning (2): Cảnh báo, tình huống bất thường nhưng không critical
//   - Error (3): Lỗi xảy ra nhưng application vẫn tiếp tục
//   - Fatal (4): Lỗi nghiêm trọng, có thể dẫn đến application dừng
//
// # Cấu hình Environment-Specific
//
// Cấu hình khác nhau cho từng môi trường:
//
//	// Development
//	devConfig := &log.Config{
//	    Level: handler.DebugLevel,
//	    Console: log.ConsoleConfig{Enabled: true, Colored: true},
//	    File:    log.FileConfig{Enabled: true, Path: "logs/dev.log"},
//	    Stack:   log.StackConfig{Enabled: true, Handlers: log.StackHandlers{Console: true, File: true}},
//	}
//
//	// Production
//	prodConfig := &log.Config{
//	    Level: handler.InfoLevel,
//	    Console: log.ConsoleConfig{Enabled: false},
//	    File:    log.FileConfig{Enabled: true, Path: "/var/log/app/app.log", MaxSize: 100*1024*1024},
//	    Stack:   log.StackConfig{Enabled: false},
//	}
//
// # Middleware Logging
//
// Tích hợp với HTTP middleware trong Fork Framework:
//
//	func LoggingMiddleware(manager log.Manager) func(http.Handler) http.Handler {
//	    logger := manager.GetLogger("HTTPMiddleware")
//
//	    return func(next http.Handler) http.Handler {
//	        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	            start := time.Now()
//
//	            logger.Info("Request started",
//	                "method", r.Method,
//	                "path", r.URL.Path,
//	                "remote_addr", r.RemoteAddr,
//	            )
//
//	            next.ServeHTTP(w, r)
//
//	            logger.Info("Request completed",
//	                "method", r.Method,
//	                "path", r.URL.Path,
//	                "duration", time.Since(start).String(),
//	            )
//	        })
//	    }
//	}
//
// # Performance Monitoring
//
// Pattern để monitor performance với logging:
//
//	type PerformanceLogger struct {
//	    logger log.Logger
//	}
//
//	func (p *PerformanceLogger) TimeOperation(operation string, fn func() error, context ...interface{}) error {
//	    start := time.Now()
//
//	    p.logger.Info("Operation started", append(context, "operation", operation)...)
//
//	    err := fn()
//	    duration := time.Since(start)
//
//	    logArgs := append(context, "operation", operation, "duration_ms", duration.Milliseconds())
//
//	    if err != nil {
//	        p.logger.Error("Operation failed", append(logArgs, "error", err.Error())...)
//	    } else {
//	        p.logger.Info("Operation completed", logArgs...)
//	    }
//
//	    return err
//	}
//
// # Runtime Handler Management
//
// Quản lý handlers trong runtime:
//
//	// Thêm custom handler
//	dbHandler := NewDatabaseHandler(db)
//	manager.AddHandler("database", dbHandler)
//
//	// Set handler cho logger cụ thể
//	manager.SetHandler("AuditService", "database")
//
//	// Remove handler
//	manager.RemoveHandler("database")
//
// # Testing với Mock Loggers
//
// Package cung cấp mock interfaces cho testing:
//
//	func TestUserService_CreateUser(t *testing.T) {
//	    // Setup mock logger
//	    mockLogger := &MockLogger{}
//	    service := &UserService{logger: mockLogger}
//
//	    // Execute
//	    user := &User{Username: "testuser"}
//	    err := service.CreateUser(user)
//
//	    // Assert
//	    assert.NoError(t, err)
//	    assert.Contains(t, mockLogger.Logs, "Creating user")
//	    assert.Contains(t, mockLogger.Logs, "User created successfully")
//	}
//
// # Default Configuration
//
// Package cung cấp cấu hình mặc định phù hợp cho hầu hết use cases:
//
//	config := log.DefaultConfig()
//	// Level: InfoLevel
//	// Console: Enabled=true, Colored=true
//	// File: Enabled=false
//	// Stack: Enabled=false
//
// # Structured Logging
//
// Sử dụng structured logging với key-value pairs:
//
//	logger.Info("User action",
//	    "user_id", 12345,
//	    "action", "login",
//	    "ip_address", "192.168.1.100",
//	    "user_agent", "Mozilla/5.0...",
//	    "success", true,
//	)
//
// # Performance Characteristics
//
//   - Shared Handlers: Giảm memory footprint bằng cách chia sẻ handlers
//   - Level Filtering: Logs được filter sớm để tránh xử lý không cần thiết
//   - Concurrent Safe: Thread-safe cho các ứng dụng concurrent
//   - Zero Allocation: Tối ưu allocation cho hot paths
//   - Resource Management: Tự động cleanup và proper handler closure
//
// # Documentation
//
// Tài liệu chi tiết bằng tiếng Việt có sẵn tại:
//   - docs/index.md: Trang chủ với quick start guide
//   - docs/overview.md: Kiến trúc với Mermaid diagrams
//   - docs/configuration.md: Hướng dẫn cấu hình chi tiết
//   - docs/handler.md: Tài liệu về các handler types
//   - docs/logger.md: Interface và patterns sử dụng
//   - docs/workflows.md: Quy trình và integration workflows
//
// # Requirements
//
//   - Go 1.19+
//   - Fork Framework (optional, for full integration)
//
// # Installation
//
//	go get go.fork.vn/log@latest
//
// # Compatibility
//
// Package này được thiết kế đặc biệt cho Fork Framework và tương thích với:
//   - go.fork.vn/di (Dependency Injection)
//   - go.fork.vn/config (Configuration Management)
//   - Fork Framework service provider pattern
//   - Standard Go logging interfaces
package log
