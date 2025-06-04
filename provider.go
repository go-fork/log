package log

import (
	"os"
	"path/filepath"

	"go.fork.vn/config"
	"go.fork.vn/di"
	"go.fork.vn/log/handler"
)

// ServiceProvider triển khai interface di.ServiceProvider cho các dịch vụ logging.
//
// Provider này tự động hóa việc đăng ký các dịch vụ logging trong một container
// dependency injection, thiết lập các handlers cho console và file với các giá trị mặc định hợp lý.
type ServiceProvider struct{}

// NewServiceProvider tạo một provider dịch vụ log mới.
//
// Sử dụng hàm này để tạo một provider có thể được đăng ký với
// một instance di.Container.
//
// Trả về:
//   - di.ServiceProvider: một service provider cho logging
//
// Ví dụ:
//
//	app := myapp.New()
//	app.Register(log.NewServiceProvider())
func NewServiceProvider() di.ServiceProvider {
	return &ServiceProvider{}
}

// Register đăng ký các dịch vụ logging với container của ứng dụng.
//
// Phương thức này:
//   - Lấy config manager từ container bằng MustMake
//   - Unmarshal log configuration từ key "log"
//   - Tạo log manager với các handlers dựa trên configuration
//   - Đăng ký manager trong container DI
//
// Nếu không có config hoặc config không hợp lệ, sử dụng default configuration.
// Handlers được tạo dựa trên cấu hình: console, file, và stack handlers.
//
// Tham số:
//   - app: di.Application - instance của ứng dụng cung cấp Container()
func (p *ServiceProvider) Register(app di.Application) {
	if app == nil {
		panic("application cannot be nil")
	}

	c := app.Container()
	if c == nil {
		panic("container cannot be nil")
	}

	// Lấy config manager từ container bằng MustMake
	configManager, ok := c.MustMake("config").(config.Manager)
	if !ok {
		panic("config manager not found or invalid type")
	}

	// Khởi tạo với default config
	logConfig := DefaultConfig()

	// Unmarshal log configuration, nếu lỗi thì panic
	if err := configManager.UnmarshalKey("log", logConfig); err != nil {
		panic("failed to unmarshal log config: " + err.Error())
	}

	// Validate configuration, nếu lỗi thì panic
	if err := logConfig.Validate(); err != nil {
		panic("invalid log config: " + err.Error())
	}

	// Tạo log manager mới
	manager := NewManager()

	// Thêm handlers dựa trên configuration
	if logConfig.Console.Enabled {
		consoleHandler := handler.NewConsoleHandler(logConfig.Console.Colored)
		manager.AddHandler("console", consoleHandler)
	}

	if logConfig.File.Enabled && logConfig.File.Path != "" {
		// Đảm bảo thư mục log tồn tại
		logDir := filepath.Dir(logConfig.File.Path)

		// Sử dụng absolute path hoặc current working directory
		if !filepath.IsAbs(logConfig.File.Path) {
			// Nếu path là relative, sử dụng working directory
			workDir, _ := os.Getwd()
			logConfig.File.Path = filepath.Join(workDir, logConfig.File.Path)
			logDir = filepath.Dir(logConfig.File.Path)
		}

		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			if err := os.MkdirAll(logDir, 0755); err != nil {
				panic("failed to create log directory: " + err.Error())
			}
		}

		fileHandler, err := handler.NewFileHandler(logConfig.File.Path, logConfig.File.MaxSize)
		if err != nil {
			panic("failed to create file handler: " + err.Error())
		}
		manager.AddHandler("file", fileHandler)
	}

	if logConfig.Stack.Enabled {
		stackHandler := handler.NewStackHandler()

		// Thêm sub-handlers cho stack
		if logConfig.Stack.Handlers.Console {
			consoleHandler := handler.NewConsoleHandler(logConfig.Console.Colored)
			stackHandler.AddHandler(consoleHandler)
		}

		if logConfig.Stack.Handlers.File && logConfig.File.Path != "" {
			fileHandler, err := handler.NewFileHandler(logConfig.File.Path, logConfig.File.MaxSize)
			if err != nil {
				panic("failed to create stack file handler: " + err.Error())
			}
			stackHandler.AddHandler(fileHandler)
		}

		manager.AddHandler("stack", stackHandler)
	}

	// Đăng ký log manager trong container
	c.Instance("log", manager) // Dịch vụ logging chung
}

// Boot thực hiện thiết lập sau đăng ký cho dịch vụ logging.
//
// Đối với provider logging, hiện tại đây là no-op vì tất cả thiết lập
// được thực hiện trong quá trình đăng ký.
//
// Tham số:
//   - app: di.Application - instance của ứng dụng
func (p *ServiceProvider) Boot(app di.Application) {
	// Không yêu cầu thiết lập bổ sung sau khi đăng ký
	if app == nil {
		panic("application cannot be nil")
	}

	c := app.Container()
	if c == nil {
		panic("container cannot be nil")
	}
}

// Requires trả về danh sách các provider mà log provider phụ thuộc vào.
//
// Log provider không có dependency bắt buộc với provider khác, nên phương thức này
// trả về một slice rỗng.
//
// Returns:
//   - []string: Một slice rỗng vì không có dependencies bắt buộc
func (p *ServiceProvider) Requires() []string {
	return []string{} // Không có dependencies bắt buộc
}

// Providers trả về danh sách các service mà log provider đăng ký.
//
// Log provider đăng ký log manager vào container với key:
// - "log": Service logging chính (Manager)
//
// Returns:
//   - []string: Mảng chứa tên của các services được đăng ký
func (p *ServiceProvider) Providers() []string {
	return []string{"log"}
}
