// Package log cung cấp hệ thống logging linh hoạt, dễ mở rộng và thread-safe cho ứng dụng Go.
//
// # Tổng quan
//
// Package này triển khai hệ thống logging với nhiều cấp độ nghiêm trọng, các output
// handler khác nhau (console, file, v.v.), và interface quản lý tập trung. Nó được thiết kế
// để thread-safe và quản lý tài nguyên hiệu quả trong các ứng dụng concurrent.
//
// # Tính năng
//
//   - Lọc theo cấp độ log (Debug, Info, Warning, Error, Fatal)
//   - Nhiều output handler hoạt động đồng thời
//   - Hoạt động thread-safe
//   - Hỗ trợ chuỗi định dạng
//   - Khả năng mở rộng handler tùy chỉnh
//   - Output console có màu
//   - Tự động xoay vòng file log
//   - Hỗ trợ dependency injection
//
// # Sử dụng trong Go-Fork Framework
//
// Log package là Core Provider được tự động đăng ký khi khởi tạo ứng dụng.
// Fork HTTP Framework (package fork) cung cấp web context và routing,
// trong khi Fork Application (package app) quản lý dependency injection và lifecycle.
//
//	// 1. Khởi tạo ứng dụng Go-Fork (log được tự động đăng ký như Core Provider)
//	config := map[string]interface{}{
//	    "name": "myapp",
//	    "path": "./configs",
//	}
//	app := app.New(config)
//
//	// 2. Tạo file cấu hình configs/myapp.yaml
//	log:
//	  level: "info"
//	  console:
//	    enabled: true
//	    colored: true
//	  file:
//	    enabled: true
//	    path: "storage/logs/app.log"
//	    max_size: 10485760  # 10MB
//	  stack:
//	    enabled: false
//	    handlers:
//	      console: true
//	      file: true
//
//	// 3. Sử dụng log trong controller
//	func (c *UserController) Create(ctx *fork.Context) error {
//	    logger := ctx.App().Log()
//
//	    logger.Info("Creating new user: %s", userData.Email)
//
//	    user, err := c.userService.Create(userData)
//	    if err != nil {
//	        logger.Error("Failed to create user: %v", err)
//	        return ctx.JSON(500, map[string]string{"error": "Internal server error"})
//	    }
//
//	    logger.Info("User created successfully: ID=%d", user.ID)
//	    return ctx.JSON(201, user)
//	}
//
//	// 4. Sử dụng trong middleware
//	func LoggingMiddleware() fork.MiddlewareFunc {
//	    return func(c *fork.Context) error {
//	        logger := c.App().Log()
//
//	        start := time.Now()
//	        err := c.Next()
//	        duration := time.Since(start)
//
//	        logger.Info("HTTP %s %s - %d (%v)",
//	            c.Request().Method,
//	            c.Request().URL.Path,
//	            c.Response().StatusCode,
//	            duration)
//
//	        return err
//	    }
//	}
//
// # Sử dụng nâng cao trong Go-Fork
//
// Đối với các yêu cầu logging phức tạp, có thể truy cập log manager từ bất kỳ đâu trong ứng dụng:
//
//	// Trong service layer với dependency injection
//	type UserService struct {
//	    app app.Application
//	}
//
//	func NewUserService(app app.Application) *UserService {
//	    return &UserService{
//	        app: app,
//	    }
//	}
//
//	func (s *UserService) ProcessPayment(userID int, amount float64) error {
//	    logger := s.app.Log()
//
//	    logger.Info("Processing payment for user %d: $%.2f", userID, amount)
//
//	    // Xử lý logic payment
//	    if amount <= 0 {
//	        logger.Warning("Invalid payment amount for user %d: $%.2f", userID, amount)
//	        return errors.New("invalid amount")
//	    }
//
//	    logger.Debug("Payment validation passed for user %d", userID)
//	    return nil
//	}
//
//	// Cấu hình log level động từ ứng dụng
//	func (app *Application) SetLogLevel(level string) error {
//	    logger := app.Log()
//
//	    switch level {
//	    case "debug":
//	        logger.SetMinLevel(handler.DebugLevel)
//	    case "info":
//	        logger.SetMinLevel(handler.InfoLevel)
//	    case "warning":
//	        logger.SetMinLevel(handler.WarningLevel)
//	    case "error":
//	        logger.SetMinLevel(handler.ErrorLevel)
//	    default:
//	        return fmt.Errorf("invalid log level: %s", level)
//	    }
//
//	    logger.Info("Log level changed to: %s", level)
//	    return nil
//	}
//
// # Tích hợp với Go-Fork Components
//
// Log package tích hợp hoàn toàn với các component khác trong Go-Fork framework:
//
// ## Database Operations
//
//	func (r *UserRepository) Create(user *User) error {
//	    logger := r.app.Log()
//
//	    logger.Debug("Creating user in database: %+v", user)
//
//	    result := r.db.Create(user)
//	    if result.Error != nil {
//	        logger.Error("Database error creating user: %v", result.Error)
//	        return result.Error
//	    }
//
//	    logger.Info("User created successfully: ID=%d, Email=%s", user.ID, user.Email)
//	    return nil
//	}
//
// ## Queue Jobs
//
//	func (j *EmailJob) Handle(data []byte) error {
//	    logger := j.app.Log()
//
//	    logger.Info("Processing email job: %s", string(data))
//
//	    if err := j.sendEmail(data); err != nil {
//	        logger.Error("Failed to send email: %v", err)
//	        return err
//	    }
//
//	    logger.Info("Email sent successfully")
//	    return nil
//	}
//
// ## Scheduled Tasks
//
//	func (t *CleanupTask) Run() error {
//	    logger := t.app.Log()
//
//	    logger.Info("Starting cleanup task")
//
//	    deleted, err := t.cleanOldFiles()
//	    if err != nil {
//	        logger.Error("Cleanup task failed: %v", err)
//	        return err
//	    }
//
//	    logger.Info("Cleanup completed: %d files deleted", deleted)
//	    return nil
//	}
//
// # Cấu hình Log Handlers
//
// Log package hỗ trợ 3 loại handler chính được cấu hình qua YAML:
//
// - Console Handler: Ghi log ra console với hỗ trợ màu sắc
// - File Handler: Ghi log vào file với tự động rotation
// - Stack Handler: Kết hợp nhiều handler cùng lúc
//
// Tất cả handlers được quản lý tự động bởi ServiceProvider và không cần
// cấu hình thủ công trong code ứng dụng.
//
// # Xem thêm
//
// - Interface Manager và triển khai DefaultManager cho các thao tác logging chính
// - Package handler để hiểu về các loại output handler
// - go.fork.vn/config để cấu hình log qua YAML files
// - go.fork.vn/di để hiểu về dependency injection trong Go-Fork
//
// # Tương thích Go-Fork Framework
//
// Module này được thiết kế đặc biệt cho Go-Fork framework version v0.1.1 trở lên,
// triển khai đầy đủ interface ServiceProvider với các phương thức Register, Boot,
// Requires và Providers theo chuẩn Go-Fork dependency injection.
package log
