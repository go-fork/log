package log

import (
	"go.fork.vn/config"
	"go.fork.vn/di"
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

	// Tạo log manager mới với config
	manager := NewManager(logConfig)

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
