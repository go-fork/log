# Log Package - Fork Framework

[![Go Version](https://img.shields.io/badge/go-1.23.9+-blue.svg)](https://golang.org)
[![Fork Framework](https://img.shields.io/badge/fork-v0.1.2-green.svg)](https://fork.vn)
[![License](https://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)

Package log cung cấp hệ thống logging linh hoạt, có thể mở rộng và thread-safe cho các ứng dụng Go trong hệ sinh thái Fork Framework.

## Tổng quan

Log package là Core Provider được tự động đăng ký và cấu hình khi khởi tạo ứng dụng Fork Framework. Package này triển khai hệ thống logging toàn diện với nhiều mức độ nghiêm trọng, các handler output khác nhau, và giao diện quản lý tập trung.

## Tính năng

- ✅ **Tích hợp hoàn toàn với Fork Framework** - Tự động đăng ký như Core Provider
- ✅ **Lọc log đa cấp** - Debug, Info, Warning, Error, Fatal  
- ✅ **Nhiều handler output** - Console, File, Stack handlers
- ✅ **Thread-safe operations** - Thao tác an toàn trong môi trường đồng thời
- ✅ **Cấu hình YAML** - Quản lý cấu hình qua file YAML
- ✅ **Dependency Injection** - Tích hợp DI container của Fork
- ✅ **Auto file rotation** - Tự động xoay file với kích thước và thời gian trigger
- ✅ **Colored console output** - Hỗ trợ màu sắc cho development
- ✅ **Printf-style formatting** - Định dạng kiểu Printf
- ✅ **Memory leak prevention** - Quản lý tài nguyên đúng cách

## Cài đặt

```bash
go get go.fork.vn/log@v0.1.2
```

## Sử dụng trong Fork Framework

### 1. Cấu hình trong YAML

Tạo file `configs/app.yaml`:

```yaml
log:
  level: "info"
  console:
    enabled: true
    colored: true
  file:
    enabled: true
    path: "storage/logs/app.log"
    max_size: 10485760  # 10MB
  stack:
    enabled: false
    handlers:
      console: true
      file: true
```

### 2. Khởi tạo Fork Application

```go
package main

import (
    "go.fork.vn/core"
    "go.fork.vn/fork"
)

func main() {
    // Log được auto-register như Core Provider
    config := map[string]interface{}{
        "name": "myapp",
        "path": "./configs",
    }
    app := app.New(config)
    
    // Tạo Fork HTTP server
    server := fork.New(app)
    
    // Log đã sẵn sàng sử dụng trong controllers
    server.Start(":8080")
}
```

### 3. Sử dụng trong Controller

```go
func (c *UserController) Create(ctx *fork.Context) error {
    logger := ctx.App().Log()
    
    logger.Info("Creating new user: %s", userData.Email)
    
    user, err := c.userService.Create(userData)
    if err != nil {
        logger.Error("Failed to create user: %v", err)
        return ctx.JSON(500, map[string]string{"error": "Internal server error"})
    }
    
    logger.Info("User created successfully: ID=%d", user.ID)
    return ctx.JSON(201, user)
}
```

### 4. Sử dụng trong Middleware

```go
func LoggingMiddleware() fork.MiddlewareFunc {
    return func(c *fork.Context) error {
        logger := c.App().Log()
        
        start := time.Now()
        err := c.Next()
        duration := time.Since(start)
        
        logger.Info("HTTP %s %s - %d (%v)",
            c.Request().Method,
            c.Request().URL.Path,
            c.Response().StatusCode,
            duration)
        
        return err
    }
}
```

### 5. Sử dụng trong Service Layer

```go
type UserService struct {
    app app.Application
}

func NewUserService(app app.Application) *UserService {
    return &UserService{app: app}
}

func (s *UserService) ProcessPayment(userID int, amount float64) error {
    logger := s.app.Log()
    
    logger.Info("Processing payment for user %d: $%.2f", userID, amount)
    
    if amount <= 0 {
        logger.Warning("Invalid payment amount for user %d: $%.2f", userID, amount)
        return errors.New("invalid amount")
    }
    
    logger.Debug("Payment validation passed for user %d", userID)
    return nil
}
```

## Tích hợp với Fork Components

### Database Operations

```go
func (r *UserRepository) Create(user *User) error {
    logger := r.app.Log()
    
    logger.Debug("Creating user in database: %+v", user)
    
    result := r.db.Create(user)
    if result.Error != nil {
        logger.Error("Database error creating user: %v", result.Error)
        return result.Error
    }
    
    logger.Info("User created successfully: ID=%d, Email=%s", user.ID, user.Email)
    return nil
}
```

### Queue Jobs

```go
func (j *EmailJob) Handle(data []byte) error {
    logger := j.app.Log()
    
    logger.Info("Processing email job: %s", string(data))
    
    if err := j.sendEmail(data); err != nil {
        logger.Error("Failed to send email: %v", err)
        return err
    }
    
    logger.Info("Email sent successfully")
    return nil
}
```

### Scheduled Tasks

```go
func (t *CleanupTask) Run() error {
    logger := t.app.Log()
    
    logger.Info("Starting cleanup task")
    
    deleted, err := t.cleanOldFiles()
    if err != nil {
        logger.Error("Cleanup task failed: %v", err)
        return err
    }
    
    logger.Info("Cleanup completed: %d files deleted", deleted)
    return nil
}
```

## Log Handlers

Log package hỗ trợ 3 loại handler chính:

### Console Handler
Ghi log ra console với hỗ trợ màu sắc cho development

### File Handler  
Ghi log vào file với tự động rotation khi đạt kích thước tối đa

### Stack Handler
Kết hợp nhiều handler cùng lúc để ghi log ra nhiều đích

## Cấu hình Log Level

Có thể thay đổi log level động từ ứng dụng:

```go
func (app *Application) SetLogLevel(level string) error {
    logger := app.Log()
    
    switch level {
    case "debug":
        logger.SetMinLevel(handler.DEBUG)
    case "info":
        logger.SetMinLevel(handler.INFO)
    case "warning":
        logger.SetMinLevel(handler.WARNING)
    case "error":
        logger.SetMinLevel(handler.ERROR)
    default:
        return fmt.Errorf("invalid log level: %s", level)
    }
    
    logger.Info("Log level changed to: %s", level)
    return nil
}
```

## Tài liệu

- [Hướng dẫn sử dụng chi tiết](docs/usage.md)
- [Tổng quan kiến trúc](docs/overview.md)
- [API Documentation](https://pkg.go.dev/go.fork.vn/log)
- [Release Notes](releases/) - Chi tiết về từng phiên bản

## Release Management

Dự án sử dụng structured release management với:
- **Automated CI/CD**: GitHub Actions cho testing, building, và releasing
- **Release Notes**: Chi tiết cho từng version trong thư mục `releases/`
- **Migration Guides**: Hướng dẫn nâng cấp giữa các versions
- **Automated Dependencies**: Dependabot tự động cập nhật dependencies

## Changelog

Xem [CHANGELOG.md](CHANGELOG.md) để biết thông tin về các thay đổi trong từng phiên bản.

## License

MIT License. Xem [LICENSE](LICENSE) để biết thêm chi tiết.

## Fork Framework

Package này là một phần của [Fork Framework](https://fork.vn) - Modern web framework cho Go.

---

**Lưu ý**: Package này được thiết kế đặc biệt để sử dụng trong Fork Framework và không khuyến khích sử dụng độc lập.
