package log

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	configMocks "go.fork.vn/config/mocks"
	"go.fork.vn/di"
	diMocks "go.fork.vn/di/mocks"
	"go.fork.vn/log/handler"
)

// setupMockApplication thiết lập một mock Application với Container đã cấu hình
func setupMockApplication() (*diMocks.Application, *di.Container) {
	container := di.New()

	mockApp := new(diMocks.Application)
	mockApp.On("Container").Return(container)
	mockApp.On("BasePath", mock.Anything).Return(filepath.Join(os.TempDir(), "logs"))

	return mockApp, container
}

// setupConfigManager thiết lập một mockConfigManager đã cấu hình
func setupConfigManager() *configMocks.MockManager {
	mockConfigManager := new(configMocks.MockManager)

	// Cấu hình mockConfigManager để xử lý UnmarshalKey với đúng config
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
		// Gán cấu hình cho log config
		config := args.Get(1).(*Config)
		config.Level = "info"
		config.Console.Enabled = true
		config.Console.Colored = true
		config.File.Enabled = true
		config.File.Path = filepath.Join(os.TempDir(), "logs", "app.log")
		config.File.MaxSize = 10485760
	}).Return(nil)

	return mockConfigManager
}

func TestNewServiceProvider(t *testing.T) {
	provider := NewServiceProvider()
	assert.NotNil(t, provider, "NewServiceProvider() không được trả về nil")
}

func TestServiceProviderRegister(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication()

	// Tạo mock config manager
	mockConfigManager := setupConfigManager()

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider với application
	provider.Register(mockApp)

	// Xác minh mockApp.Container đã được gọi
	mockApp.AssertCalled(t, "Container")
	mockConfigManager.AssertCalled(t, "UnmarshalKey", "log", mock.AnythingOfType("*log.Config"))

	// Kiểm tra binding "log"
	managerInstance, err := container.Make("log")
	assert.NoError(t, err, "ServiceProvider phải đăng ký binding 'log'")

	manager, ok := managerInstance.(Manager)
	assert.True(t, ok, "Binding 'log' phải là kiểu Manager, nhưng nhận được %T", managerInstance)

	// Kiểm tra binding "log.manager"
	managerInstance2, err := container.Make("log.manager")
	assert.NoError(t, err, "ServiceProvider phải đăng ký binding 'log.manager'")
	assert.Equal(t, managerInstance, managerInstance2, "'log' và 'log.manager' phải trỏ đến cùng một instance")

	// Kiểm tra handlers được thiết lập đúng
	// Kiểm tra console handler
	consoleHandler := manager.GetHandler("console")
	assert.NotNil(t, consoleHandler, "Manager phải có console handler")
	_, ok = consoleHandler.(*handler.ConsoleHandler)
	assert.True(t, ok, "Console handler phải có kiểu đúng, nhưng nhận được %T", consoleHandler)

	// Kiểm tra file handler
	fileHandler := manager.GetHandler("file")
	assert.NotNil(t, fileHandler, "Manager phải có file handler")
	_, ok = fileHandler.(*handler.FileHandler)
	assert.True(t, ok, "File handler phải có kiểu đúng, nhưng nhận được %T", fileHandler)

	// Dọn dẹp
	if err := manager.Close(); err != nil {
		t.Logf("Không thể đóng manager: %v", err)
	}
}

func TestServiceProviderBoot(t *testing.T) {
	// 1. Test với mock application hợp lệ
	mockApp, container := setupMockApplication()

	// Thêm binding "log" để kiểm tra xác minh trong Boot
	container.Instance("log", NewManager())

	// Tạo service provider
	provider := NewServiceProvider()

	// Boot không làm gì nhưng nên chạy không lỗi
	provider.Boot(mockApp)

	// Không có hành động cụ thể trong Boot, nên chỉ kiểm tra không có panic
	// Không cần mock.AssertExpectations vì không có lời gọi cụ thể nào được mong đợi

	// 2. Test với nil application
	var nilApp di.Application = nil
	provider.Boot(nilApp) // Không nên panic

	// 3. Test với application có container nil
	mockAppNilContainer := new(diMocks.Application)
	mockAppNilContainer.On("Container").Return(nil)
	provider.Boot(mockAppNilContainer) // Không nên panic
}

func TestServiceProviderWithConfigError(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication()

	// Tạo mock config manager với lỗi
	mockConfigManager := new(configMocks.MockManager)
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Return(
		errors.New("config error"))

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Register nên panic khi config manager trả về lỗi
	assert.Panics(t, func() {
		provider.Register(mockApp)
	}, "ServiceProvider.Register nên panic khi config manager trả về lỗi")
}

func TestServiceProviderWithInvalidConfig(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication()

	// Tạo mock config manager với cấu hình không hợp lệ
	mockConfigManager := new(configMocks.MockManager)
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
		// Gán cấu hình không hợp lệ cho log config
		config := args.Get(1).(*Config)
		config.Level = "invalid_level" // Level không hợp lệ
	}).Return(nil)

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Register nên panic khi validation config trả về lỗi
	assert.Panics(t, func() {
		provider.Register(mockApp)
	}, "ServiceProvider.Register nên panic khi cấu hình không hợp lệ")
}

func TestServiceProviderWithStackHandler(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication()

	// Tạo mock config manager với cấu hình stack handler
	mockConfigManager := new(configMocks.MockManager)
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
		// Cấu hình với stack handler
		config := args.Get(1).(*Config)
		config.Level = "info"
		config.Console.Enabled = true
		config.Console.Colored = true
		config.File.Enabled = true
		config.File.Path = filepath.Join(os.TempDir(), "logs", "app.log")
		config.File.MaxSize = 10485760
		config.Stack.Enabled = true
		config.Stack.Handlers.Console = true
		config.Stack.Handlers.File = true
	}).Return(nil)

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider với application
	provider.Register(mockApp)

	// Xác minh mockApp.Container đã được gọi
	mockApp.AssertCalled(t, "Container")
	mockConfigManager.AssertCalled(t, "UnmarshalKey", "log", mock.AnythingOfType("*log.Config"))

	// Kiểm tra binding "log"
	managerInstance, err := container.Make("log")
	assert.NoError(t, err, "ServiceProvider phải đăng ký binding 'log'")

	manager, ok := managerInstance.(Manager)
	assert.True(t, ok, "Binding 'log' phải là kiểu Manager")

	// Kiểm tra stack handler
	stackHandler := manager.GetHandler("stack")
	assert.NotNil(t, stackHandler, "Manager phải có stack handler")
	_, ok = stackHandler.(*handler.StackHandler)
	assert.True(t, ok, "Stack handler phải có kiểu đúng, nhưng nhận được %T", stackHandler)

	// Dọn dẹp
	if err := manager.Close(); err != nil {
		t.Logf("Không thể đóng manager: %v", err)
	}
}

func TestContainerBindingResolution(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication()

	// Tạo mock config manager
	mockConfigManager := setupConfigManager()

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider
	provider.Register(mockApp)

	// Thêm một binding phụ thuộc vào log manager
	container.Bind("custom.logger", func(c *di.Container) interface{} {
		// Lấy log manager từ container
		manager, err := c.Make("log")
		if err != nil {
			t.Fatal("Không thể resolve dependency 'log':", err)
		}

		// Trả về một struct sử dụng log manager
		return struct {
			LogManager Manager
			Name       string
		}{
			LogManager: manager.(Manager),
			Name:       "CustomLogger",
		}
	})

	// Giải quyết binding
	customLogger, err := container.Make("custom.logger")
	assert.NoError(t, err, "Phải resolve binding 'custom.logger' thành công")

	// Kiểm tra cấu trúc được trả về
	loggerStruct, ok := customLogger.(struct {
		LogManager Manager
		Name       string
	})

	assert.True(t, ok, "Binding 'custom.logger' phải trả về kiểu đúng, nhưng nhận được: %T", customLogger)
	assert.Equal(t, "CustomLogger", loggerStruct.Name, "Tên phải đúng")
	assert.NotNil(t, loggerStruct.LogManager, "LogManager không được là nil")
}

// TestServiceProviderRequires kiểm tra method Requires() trả về giá trị đúng
func TestServiceProviderRequires(t *testing.T) {
	// Tạo service provider
	provider := NewServiceProvider()

	// Lấy danh sách dependencies
	requires := provider.Requires()

	// Log provider không phụ thuộc vào provider nào khác
	assert.Empty(t, requires, "Log provider không nên phụ thuộc vào bất kỳ provider nào")
}

// TestServiceProviderProviders kiểm tra method Providers() trả về giá trị đúng
func TestServiceProviderProviders(t *testing.T) {
	// Tạo service provider
	provider := NewServiceProvider()

	// Lấy danh sách services được đăng ký
	providers := provider.Providers()

	// Kiểm tra số lượng và giá trị cụ thể
	expectedServices := []string{"log", "log.manager"}
	assert.ElementsMatch(t, expectedServices, providers,
		"Provider phải đăng ký đúng các services: %v, nhưng nhận được: %v",
		expectedServices, providers)
}

// TestRegisterWithInvalidInputs kiểm tra các trường hợp đầu vào không hợp lệ cho Register
func TestRegisterWithInvalidInputs(t *testing.T) {
	// Tạo service provider
	provider := NewServiceProvider()

	// 1. Đăng ký provider với nil application
	assert.Panics(t, func() {
		var nilApp di.Application = nil
		provider.Register(nilApp)
	}, "ServiceProvider.Register nên panic khi app là nil")

	// 2. Đăng ký provider với application có nil container
	mockAppNilContainer := new(diMocks.Application)
	mockAppNilContainer.On("Container").Return(nil)
	assert.Panics(t, func() {
		provider.Register(mockAppNilContainer)
	}, "ServiceProvider.Register nên panic khi container là nil")

	// 3. Đăng ký provider với container không có config manager
	mockApp, container := setupMockApplication()
	assert.Panics(t, func() {
		// Không cài đặt config manager
		provider.Register(mockApp)
	}, "ServiceProvider.Register nên panic khi config manager không tồn tại")

	// 4. Đăng ký provider với config manager có kiểu không đúng
	container.Instance("config", "not-a-config-manager")
	assert.Panics(t, func() {
		provider.Register(mockApp)
	}, "ServiceProvider.Register nên panic khi config manager có kiểu không đúng")
}
