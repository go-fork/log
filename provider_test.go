package log

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.fork.vn/config/mocks"
	"go.fork.vn/di"
	diMocks "go.fork.vn/di/mocks"
	"go.fork.vn/log/handler"
)

// setupMockApplication thiết lập một mock Application với Container đã cấu hình
func setupMockApplication(t *testing.T) (*diMocks.MockApplication, di.Container) {
	container := di.New()

	mockApp := diMocks.NewMockApplication(t)
	mockApp.On("Container").Return(container).Maybe()

	return mockApp, container
}

func TestServiceProvider_New(t *testing.T) {
	provider := NewServiceProvider()
	assert.NotNil(t, provider, "NewServiceProvider() không được trả về nil")
}

func TestServiceProvider_Register(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo thư mục log trước khi chạy test
	logDir := filepath.Join(os.TempDir(), "logs")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}
	defer os.RemoveAll(logDir)

	// Tạo mock config manager
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
		config := args.Get(1).(*Config)
		config.Level = handler.InfoLevel
		config.Console.Enabled = true
		config.Console.Colored = true
		config.File.Enabled = true
		config.File.Path = filepath.Join(os.TempDir(), "logs", "app.log")
		config.File.MaxSize = 10485760
	}).Return(nil).Once()

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider với application
	provider.Register(mockApp)

	// Kiểm tra binding "log"
	managerInstance, err := container.Make("log")
	assert.NoError(t, err, "ServiceProvider phải đăng ký binding 'log'")

	manager, ok := managerInstance.(Manager)
	assert.True(t, ok, "Binding 'log' phải là kiểu Manager, nhưng nhận được %T", managerInstance)

	// Kiểm tra handlers được thiết lập đúng
	// Kiểm tra console handler
	consoleHandler := manager.GetHandler(HandlerTypeConsole)
	assert.NotNil(t, consoleHandler, "Manager phải có console handler")
	_, ok = consoleHandler.(*handler.ConsoleHandler)
	assert.True(t, ok, "Console handler phải có kiểu đúng, nhưng nhận được %T", consoleHandler)

	// Kiểm tra file handler
	fileHandler := manager.GetHandler(HandlerTypeFile)
	assert.NotNil(t, fileHandler, "Manager phải có file handler")
	_, ok = fileHandler.(*handler.FileHandler)
	assert.True(t, ok, "File handler phải có kiểu đúng, nhưng nhận được %T", fileHandler)

	// Dọn dẹp
	if err := manager.Close(); err != nil {
		t.Logf("Không thể đóng manager: %v", err)
	}
}

func TestServiceProvider_Boot(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func() di.Application
		expectPanic bool
	}{
		{
			name: "valid_application_with_log_binding",
			setupMocks: func() di.Application {
				mockApp, container := setupMockApplication(t)
				config := createTestConfigForProvider()
				container.Instance("log", NewManager(config))
				return mockApp
			},
			expectPanic: false,
		},
		{
			name: "nil_application",
			setupMocks: func() di.Application {
				return nil
			},
			expectPanic: true,
		},
		{
			name: "application_with_nil_container",
			setupMocks: func() di.Application {
				mockApp := diMocks.NewMockApplication(t)
				mockApp.On("Container").Return(nil).Maybe()
				return mockApp
			},
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewServiceProvider()
			app := tt.setupMocks()

			if tt.expectPanic {
				assert.Panics(t, func() {
					provider.Boot(app)
				})
			} else {
				assert.NotPanics(t, func() {
					provider.Boot(app)
				})
			}
		})
	}
}

func TestServiceProvider_WithConfigError(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo mock config manager với lỗi
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Return(
		errors.New("config error")).Once()

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Register nên panic khi config manager trả về lỗi
	assert.Panics(t, func() {
		provider.Register(mockApp)
	}, "ServiceProvider.Register nên panic khi config manager trả về lỗi")
}

func TestServiceProvider_WithInvalidConfig(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo mock config manager với cấu hình không hợp lệ
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
		// Gán cấu hình không hợp lệ cho log config
		config := args.Get(1).(*Config)
		config.Level = handler.Level(99) // Level không hợp lệ
	}).Return(nil).Once()

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Register nên panic khi validation config trả về lỗi
	assert.Panics(t, func() {
		provider.Register(mockApp)
	}, "ServiceProvider.Register nên panic khi cấu hình không hợp lệ")
}

func TestServiceProvider_WithStackHandler(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo thư mục log trước khi chạy test
	logDir := filepath.Join(os.TempDir(), "logs")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}
	defer os.RemoveAll(logDir)

	// Tạo mock config manager với cấu hình stack handler
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
		// Cấu hình với stack handler
		config := args.Get(1).(*Config)
		config.Level = handler.InfoLevel
		config.Console.Enabled = true
		config.Console.Colored = true
		config.File.Enabled = true
		config.File.Path = filepath.Join(os.TempDir(), "logs", "app.log")
		config.File.MaxSize = 10485760
		config.Stack.Enabled = true
		config.Stack.Handlers.Console = true
		config.Stack.Handlers.File = true
	}).Return(nil).Once()

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider với application
	provider.Register(mockApp)

	// Kiểm tra binding "log"
	managerInstance, err := container.Make("log")
	assert.NoError(t, err, "ServiceProvider phải đăng ký binding 'log'")

	manager, ok := managerInstance.(Manager)
	assert.True(t, ok, "Binding 'log' phải là kiểu Manager")

	// Kiểm tra stack handler
	stackHandler := manager.GetHandler(HandlerTypeStack)
	assert.NotNil(t, stackHandler, "Manager phải có stack handler")
	_, ok = stackHandler.(*handler.StackHandler)
	assert.True(t, ok, "Stack handler phải có kiểu đúng, nhưng nhận được %T", stackHandler)

	// Dọn dẹp
	if err := manager.Close(); err != nil {
		t.Logf("Không thể đóng manager: %v", err)
	}
}

func TestServiceProvider_ContainerBindingResolution(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo thư mục log trước khi chạy test
	logDir := filepath.Join(os.TempDir(), "logs")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}
	defer os.RemoveAll(logDir)

	// Tạo mock config manager
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
		config := args.Get(1).(*Config)
		config.Level = handler.InfoLevel
		config.Console.Enabled = true
		config.Console.Colored = true
		config.File.Enabled = true
		config.File.Path = filepath.Join(os.TempDir(), "logs", "app.log")
		config.File.MaxSize = 10485760
	}).Return(nil).Once()

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider
	provider.Register(mockApp)

	// Thêm một binding phụ thuộc vào log manager
	container.Bind("custom.logger", func(c di.Container) interface{} {
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
func TestServiceProvider_Requires(t *testing.T) {
	// Tạo service provider
	provider := NewServiceProvider()

	// Lấy danh sách dependencies
	requires := provider.Requires()

	// Log provider không phụ thuộc vào provider nào khác
	assert.Empty(t, requires, "Log provider không nên phụ thuộc vào bất kỳ provider nào")
}

// TestServiceProviderProviders kiểm tra method Providers() trả về giá trị đúng
func TestServiceProvider_Providers(t *testing.T) {
	// Tạo service provider
	provider := NewServiceProvider()

	// Lấy danh sách services được đăng ký
	providers := provider.Providers()

	// Kiểm tra số lượng và giá trị cụ thể - chỉ "log" được đăng ký
	expectedServices := []string{"log"}
	assert.ElementsMatch(t, expectedServices, providers,
		"Provider phải đăng ký đúng các services: %v, nhưng nhận được: %v",
		expectedServices, providers)
}

// TestRegisterWithInvalidInputs kiểm tra các trường hợp đầu vào không hợp lệ cho Register
func TestServiceProvider_RegisterWithInvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func() (di.Application, di.Container)
		expectPanic bool
		description string
	}{
		{
			name: "nil_application",
			setupMocks: func() (di.Application, di.Container) {
				return nil, nil
			},
			expectPanic: true,
			description: "ServiceProvider.Register nên panic khi app là nil",
		},
		{
			name: "application_with_nil_container",
			setupMocks: func() (di.Application, di.Container) {
				mockApp := diMocks.NewMockApplication(t)
				mockApp.On("Container").Return(nil).Once()
				return mockApp, nil
			},
			expectPanic: true,
			description: "ServiceProvider.Register nên panic khi container là nil",
		},
		{
			name: "container_without_config_manager",
			setupMocks: func() (di.Application, di.Container) {
				mockApp, container := setupMockApplication(t)
				return mockApp, container
			},
			expectPanic: true,
			description: "ServiceProvider.Register nên panic khi config manager không tồn tại",
		},
		{
			name: "container_with_invalid_config_manager_type",
			setupMocks: func() (di.Application, di.Container) {
				mockApp, container := setupMockApplication(t)
				container.Instance("config", "not-a-config-manager")
				return mockApp, container
			},
			expectPanic: true,
			description: "ServiceProvider.Register nên panic khi config manager có kiểu không đúng",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewServiceProvider()
			app, _ := tt.setupMocks()

			if tt.expectPanic {
				assert.Panics(t, func() {
					provider.Register(app)
				}, tt.description)
			} else {
				assert.NotPanics(t, func() {
					provider.Register(app)
				}, tt.description)
			}
		})
	}
}

// Helper function để tạo test config cho provider tests
func createTestConfigForProvider() *Config {
	return &Config{
		Level: handler.InfoLevel,
		Console: ConsoleConfig{
			Enabled: true,
			Colored: false,
		},
		File: FileConfig{
			Enabled: true,
			Path:    "/tmp/provider_test.log",
			MaxSize: 1024 * 1024, // 1MB
		},
		Stack: StackConfig{
			Enabled: true,
			Handlers: StackHandlers{
				Console: true,
				File:    true,
			},
		},
	}
}

// BenchmarkNewServiceProvider đo hiệu suất tạo ServiceProvider mới
func BenchmarkServiceProvider_New(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider := NewServiceProvider()
		_ = provider
	}
}

// BenchmarkServiceProviderRegister đo hiệu suất đăng ký ServiceProvider
func BenchmarkServiceProvider_Register(b *testing.B) {
	// Tạo thư mục log trước khi chạy benchmark
	logDir := filepath.Join(os.TempDir(), "logs")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		b.Fatalf("Failed to create log directory: %v", err)
	}
	defer os.RemoveAll(logDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Setup cho mỗi iteration để tránh mock conflicts
		container := di.New()

		// Tạo mock application với clean state
		mockApp := &diMocks.MockApplication{}
		mockApp.On("Container").Return(container).Once()

		// Tạo mock config manager với clean state
		mockConfigManager := &mocks.MockManager{}
		mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
			config := args.Get(1).(*Config)
			config.Level = handler.InfoLevel
			config.Console.Enabled = true
			config.Console.Colored = true
			config.File.Enabled = true
			config.File.Path = filepath.Join(os.TempDir(), "logs", "bench.log")
			config.File.MaxSize = 10485760
		}).Return(nil).Once()

		container.Instance("config", mockConfigManager)

		provider := NewServiceProvider()

		b.StartTimer()
		provider.Register(mockApp)
		b.StopTimer()

		// Cleanup để tránh memory leak
		if manager, err := container.Make("log"); err == nil {
			if m, ok := manager.(Manager); ok {
				_ = m.Close()
			}
		}
	}
}

// BenchmarkServiceProviderRegisterWithStackHandler đo hiệu suất với stack handler
func BenchmarkServiceProvider_Register_WithStackHandler(b *testing.B) {
	// Tạo thư mục log trước khi chạy benchmark
	logDir := filepath.Join(os.TempDir(), "logs")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		b.Fatalf("Failed to create log directory: %v", err)
	}
	defer os.RemoveAll(logDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		container := di.New()

		mockApp := &diMocks.MockApplication{}
		mockApp.On("Container").Return(container).Once()

		mockConfigManager := &mocks.MockManager{}
		mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
			config := args.Get(1).(*Config)
			config.Level = handler.InfoLevel
			config.Console.Enabled = true
			config.Console.Colored = true
			config.File.Enabled = true
			config.File.Path = filepath.Join(os.TempDir(), "logs", "bench_stack.log")
			config.File.MaxSize = 10485760
			config.Stack.Enabled = true
			config.Stack.Handlers.Console = true
			config.Stack.Handlers.File = true
		}).Return(nil).Once()

		container.Instance("config", mockConfigManager)
		provider := NewServiceProvider()

		b.StartTimer()
		provider.Register(mockApp)
		b.StopTimer()

		// Cleanup
		if manager, err := container.Make("log"); err == nil {
			if m, ok := manager.(Manager); ok {
				_ = m.Close()
			}
		}
	}
}

// BenchmarkServiceProviderBoot đo hiệu suất Boot method
func BenchmarkServiceProvider_Boot(b *testing.B) {
	// Setup một lần
	container := di.New()
	mockApp := diMocks.NewMockApplication(b)
	mockApp.On("Container").Return(container).Maybe()

	config := createTestConfigForProvider()
	container.Instance("log", NewManager(config))
	provider := NewServiceProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.Boot(mockApp)
	}
}

// BenchmarkServiceProviderRequires đo hiệu suất Requires method
func BenchmarkServiceProvider_Requires(b *testing.B) {
	provider := NewServiceProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		requires := provider.Requires()
		_ = requires
	}
}

// BenchmarkServiceProviderProviders đo hiệu suất Providers method
func BenchmarkServiceProvider_Providers(b *testing.B) {
	provider := NewServiceProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		providers := provider.Providers()
		_ = providers
	}
}

// BenchmarkContainerMakeLog đo hiệu suất resolve log service từ container
func BenchmarkContainer_MakeLog(b *testing.B) {
	container := di.New()
	config := createTestConfigForProvider()
	manager := NewManager(config)
	container.Instance("log", manager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logManager, err := container.Make("log")
		if err != nil {
			b.Fatal("Failed to make log:", err)
		}
		_ = logManager
	}

	// Cleanup
	_ = manager.Close()
}
