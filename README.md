# Go Fork Log Package

[![Go Version](https://img.shields.io/github/go-mod/go-version/go-fork/log)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/go-fork/log)](https://github.com/go-fork/log/releases)
[![License](https://img.shields.io/github/license/go-fork/log)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/go.fork.vn/log)](https://goreportcard.com/report/go.fork.vn/log)
[![Coverage](https://img.shields.io/badge/coverage-95%25-brightgreen)](coverage.html)

Package log cung cấp hệ thống logging có cấu trúc và hiệu suất cao cho Fork Framework, được thiết kế với kiến trúc **Shared Handlers** và **Contextual Loggers**.

## ✨ Tính Năng Chính

- 🎯 **Contextual Logging** - Logger được tạo theo context để phân biệt nguồn gốc
- 🔄 **Shared Handlers** - Handlers được chia sẻ giữa nhiều logger instances
- 🎚️ **Multi-Level Support** - Debug, Info, Warning, Error, Fatal levels
- ⚙️ **Flexible Configuration** - Cấu hình linh hoạt cho console, file, stack handlers
- 🔧 **Runtime Management** - Quản lý handlers và loggers trong runtime
- ⚡ **High Performance** - Tối ưu hóa cho ứng dụng hiệu suất cao
- 🚀 **Fork Framework Integration** - Tích hợp sâu với DI container và service providers

## 📦 Cài Đặt

```bash
go get go.fork.vn/log
```

## 🚀 Khởi Động Nhanh

### Sử Dụng Standalone

```go
package main

import (
    "go.fork.vn/log"
    "go.fork.vn/log/handler"
)

func main() {
    // Tạo cấu hình
    config := &log.Config{
        Level: handler.InfoLevel,
        Console: log.ConsoleConfig{
            Enabled: true,
            Colored: true,
        },
        File: log.FileConfig{
            Enabled: true,
            Path:    "logs/app.log",
            MaxSize: 10 * 1024 * 1024, // 10MB
        },
    }
    
    // Khởi tạo manager
    manager := log.NewManager(config)
    defer manager.Close()
    
    // Lấy logger theo context
    userLogger := manager.GetLogger("UserService")
    orderLogger := manager.GetLogger("OrderService")
    
    // Structured logging
    userLogger.Info("User logged in", "user_id", 12345, "username", "john_doe")
    orderLogger.Warning("Low stock", "product_id", "ABC123", "stock", 2)
    userLogger.Error("Database error", "error", "connection timeout")
}
```

### Sử Dụng với Fork Framework

```go
// main.go
func main() {
    app := fork.NewApplication()
    
    // Đăng ký log service provider
    logConfig := log.DefaultConfig()
    app.RegisterProvider(providers.NewLogProvider(logConfig))
    
    // Đăng ký service khác
    app.RegisterProvider(providers.NewUserServiceProvider())
    
    app.Run()
}

// services/user_service.go
type UserService struct {
    logger log.Logger
    db     *database.DB
}

func NewUserService(container *container.Container) *UserService {
    manager := container.Get("log").(log.Manager)
    
    return &UserService{
        logger: manager.GetLogger("UserService"),
        db:     container.Get("database").(*database.DB),
    }
}

func (s *UserService) CreateUser(user *User) error {
    s.logger.Info("Creating user", "username", user.Username)
    
    if err := s.db.Create(user); err != nil {
        s.logger.Error("Failed to create user", 
            "username", user.Username, 
            "error", err.Error(),
        )
        return err
    }
    
    s.logger.Info("User created successfully", 
        "user_id", user.ID,
        "username", user.Username,
    )
    
    return nil
}
```

## 🏗️ Kiến Trúc

Package log sử dụng kiến trúc **Shared Handlers** với **Contextual Loggers**:

- **Manager**: Quản lý tập trung các handlers và loggers
- **Shared Handlers**: Console, File, Stack handlers được chia sẻ
- **Contextual Loggers**: Mỗi logger có context riêng (UserService, OrderService, etc.)
- **Runtime Management**: Thêm/xóa handlers và loggers trong runtime

## 📚 Handlers

### Console Handler
- Xuất logs ra stdout/stderr  
- Hỗ trợ màu sắc theo level
- Tối ưu cho development

### File Handler
- Ghi logs vào file với rotation
- Tự động tạo directory
- Production-ready

### Stack Handler
- Kết hợp nhiều handlers
- Log đồng thời ra nhiều đích  
- Flexible configuration

### Custom Handlers
```go
type DatabaseHandler struct {
    db *sql.DB
}

func (h *DatabaseHandler) Log(level handler.Level, message string, args ...interface{}) error {
    // Custom implementation
    return nil
}

func (h *DatabaseHandler) Close() error {
    return nil
}
```

## ⚙️ Cấu Hình

### Default Configuration

```go
config := log.DefaultConfig()
// Level: InfoLevel
// Console: Enabled=true, Colored=true
// File: Enabled=false
// Stack: Enabled=false
```

### Environment-Specific Configs

```go
// Development
devConfig := &log.Config{
    Level: handler.DebugLevel,
    Console: log.ConsoleConfig{Enabled: true, Colored: true},
    File:    log.FileConfig{Enabled: true, Path: "logs/dev.log"},
    Stack:   log.StackConfig{Enabled: true, Handlers: log.StackHandlers{Console: true, File: true}},
}

// Production
prodConfig := &log.Config{
    Level: handler.InfoLevel,
    Console: log.ConsoleConfig{Enabled: false},
    File:    log.FileConfig{Enabled: true, Path: "/var/log/app/app.log", MaxSize: 100*1024*1024},
    Stack:   log.StackConfig{Enabled: false},
}
```

## 📊 Log Levels

| Level | Value | Usage |
|-------|-------|-------|
| Debug | 0 | Chi tiết implementation, chỉ development |
| Info | 1 | Thông tin chung về application flow |
| Warning | 2 | Cảnh báo, tình huống bất thường nhưng không critical |
| Error | 3 | Lỗi xảy ra nhưng application vẫn tiếp tục |
| Fatal | 4 | Lỗi nghiêm trọng, có thể dẫn đến application dừng |

## 🛠️ Advanced Usage

### Middleware Logging

```go
func LoggingMiddleware(manager log.Manager) func(http.Handler) http.Handler {
    logger := manager.GetLogger("HTTPMiddleware")
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            logger.Info("Request started",
                "method", r.Method,
                "path", r.URL.Path,
                "remote_addr", r.RemoteAddr,
            )
            
            next.ServeHTTP(w, r)
            
            logger.Info("Request completed",
                "method", r.Method,
                "path", r.URL.Path,
                "duration", time.Since(start).String(),
            )
        })
    }
}
```

### Performance Monitoring

```go
type PerformanceLogger struct {
    logger log.Logger
}

func (p *PerformanceLogger) TimeOperation(operation string, fn func() error, context ...interface{}) error {
    start := time.Now()
    
    p.logger.Info("Operation started", append(context, "operation", operation)...)
    
    err := fn()
    duration := time.Since(start)
    
    logArgs := append(context, "operation", operation, "duration_ms", duration.Milliseconds())
    
    if err != nil {
        p.logger.Error("Operation failed", append(logArgs, "error", err.Error())...)
    } else {
        p.logger.Info("Operation completed", logArgs...)
    }
    
    return err
}
```

### Runtime Handler Management

```go
// Thêm custom handler
dbHandler := NewDatabaseHandler(db)
manager.AddHandler("database", dbHandler)

// Set handler cho logger cụ thể
manager.SetHandler("AuditService", "database")

// Remove handler
manager.RemoveHandler("database")
```

## 🧪 Testing

```go
func TestUserService_CreateUser(t *testing.T) {
    // Setup mock logger
    mockLogger := &MockLogger{}
    service := &UserService{logger: mockLogger}
    
    // Execute
    user := &User{Username: "testuser"}
    err := service.CreateUser(user)
    
    // Assert
    assert.NoError(t, err)
    assert.Contains(t, mockLogger.Logs, "Creating user")
    assert.Contains(t, mockLogger.Logs, "User created successfully")
}
```

## 📖 Documentation

Tài liệu chi tiết bằng tiếng Việt:

- **[Tổng Quan](docs/overview.md)** - Kiến trúc và khái niệm logging
- **[Cấu Hình](docs/configuration.md)** - Hướng dẫn cấu hình chi tiết
- **[Handlers](docs/handler.md)** - Tài liệu về các handler types
- **[Logger](docs/logger.md)** - Interface và sử dụng logger
- **[Workflows](docs/workflows.md)** - Quy trình làm việc và patterns

## 📈 Performance

- **Shared Handlers**: Giảm memory footprint bằng cách chia sẻ handlers
- **Level Filtering**: Logs được filter sớm để tránh xử lý không cần thiết
- **Concurrent Safe**: Thread-safe cho các ứng dụng concurrent
- **Zero Allocation**: Tối ưu allocation cho hot paths

## 🔄 Migration

Đang upgrade từ version cũ? Xem [Migration Guide](releases/next/MIGRATION.md).

## 📋 Requirements

- Go 1.19+
- Fork Framework (optional, for full integration)

## 🤝 Contributing

1. Fork repository
2. Tạo feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push branch: `git push origin feature/amazing-feature`  
5. Submit Pull Request

## 📄 License

Dự án này được phân phối dưới [MIT License](LICENSE).

## 🙏 Acknowledgments

- [Fork Framework](https://github.com/go-fork) ecosystem
- Go community để inspiration và best practices
- Contributors và maintainers

---

**Lưu ý**: Package này được thiết kế đặc biệt cho Fork Framework. Để tận dụng đầy đủ tính năng, hãy sử dụng cùng với [Fork Framework](https://github.com/go-fork) và hệ sinh thái của nó.
