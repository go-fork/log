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

func TestNewServiceProvider(t *testing.T) {
	provider := NewServiceProvider()
	assert.NotNil(t, provider, "NewServiceProvider() không được trả về nil")
}

func TestServiceProviderRegister(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo mock config manager
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
		config := args.Get(1).(*Config)
		config.Level = "info"
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
	tests := []struct {
		name        string
		setupMocks  func() di.Application
		expectPanic bool
	}{
		{
			name: "valid application with log binding",
			setupMocks: func() di.Application {
				mockApp, container := setupMockApplication(t)
				container.Instance("log", NewManager())
				return mockApp
			},
			expectPanic: false,
		},
		{
			name: "nil application",
			setupMocks: func() di.Application {
				return nil
			},
			expectPanic: false,
		},
		{
			name: "application with nil container",
			setupMocks: func() di.Application {
				mockApp := diMocks.NewMockApplication(t)
				mockApp.On("Container").Return(nil).Maybe()
				return mockApp
			},
			expectPanic: false,
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

func TestServiceProviderWithConfigError(t *testing.T) {
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

func TestServiceProviderWithInvalidConfig(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo mock config manager với cấu hình không hợp lệ
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
		// Gán cấu hình không hợp lệ cho log config
		config := args.Get(1).(*Config)
		config.Level = "invalid_level" // Level không hợp lệ
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

func TestServiceProviderWithStackHandler(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo mock config manager với cấu hình stack handler
	mockConfigManager := mocks.NewMockManager(t)
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
	mockApp, container := setupMockApplication(t)

	// Tạo mock config manager
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
		config := args.Get(1).(*Config)
		config.Level = "info"
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

	// Kiểm tra số lượng và giá trị cụ thể - chỉ "log" được đăng ký
	expectedServices := []string{"log"}
	assert.ElementsMatch(t, expectedServices, providers,
		"Provider phải đăng ký đúng các services: %v, nhưng nhận được: %v",
		expectedServices, providers)
}

// TestRegisterWithInvalidInputs kiểm tra các trường hợp đầu vào không hợp lệ cho Register
func TestRegisterWithInvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func() (di.Application, di.Container)
		expectPanic bool
		description string
	}{
		{
			name: "nil application",
			setupMocks: func() (di.Application, di.Container) {
				return nil, nil
			},
			expectPanic: true,
			description: "ServiceProvider.Register nên panic khi app là nil",
		},
		{
			name: "application with nil container",
			setupMocks: func() (di.Application, di.Container) {
				mockApp := diMocks.NewMockApplication(t)
				mockApp.On("Container").Return(nil).Once()
				return mockApp, nil
			},
			expectPanic: true,
			description: "ServiceProvider.Register nên panic khi container là nil",
		},
		{
			name: "container without config manager",
			setupMocks: func() (di.Application, di.Container) {
				mockApp, container := setupMockApplication(t)
				return mockApp, container
			},
			expectPanic: true,
			description: "ServiceProvider.Register nên panic khi config manager không tồn tại",
		},
		{
			name: "container with invalid config manager type",
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

// BenchmarkNewServiceProvider đo hiệu suất tạo ServiceProvider mới
func BenchmarkNewServiceProvider(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider := NewServiceProvider()
		_ = provider
	}
}

// BenchmarkServiceProviderRegister đo hiệu suất đăng ký ServiceProvider
func BenchmarkServiceProviderRegister(b *testing.B) {
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
			config.Level = "info"
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

		// Cleanup để tránh memory leak
		b.StopTimer()
		if manager, err := container.Make("log"); err == nil {
			if logManager, ok := manager.(Manager); ok {
				logManager.Close()
			}
		}
	}
}

// BenchmarkServiceProviderRegisterWithStackHandler đo hiệu suất với stack handler
func BenchmarkServiceProviderRegisterWithStackHandler(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		container := di.New()

		mockApp := &diMocks.MockApplication{}
		mockApp.On("Container").Return(container).Once()

		mockConfigManager := &mocks.MockManager{}
		mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
			config := args.Get(1).(*Config)
			config.Level = "info"
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

		// Cleanup
		b.StopTimer()
		if manager, err := container.Make("log"); err == nil {
			if logManager, ok := manager.(Manager); ok {
				logManager.Close()
			}
		}
	}
}

// BenchmarkServiceProviderBoot đo hiệu suất Boot method
func BenchmarkServiceProviderBoot(b *testing.B) {
	// Setup một lần cho Boot benchmark vì Boot không có side effects
	container := di.New()
	mockApp := &diMocks.MockApplication{}
	mockApp.On("Container").Return(container).Maybe()
	container.Instance("log", NewManager())

	provider := NewServiceProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.Boot(mockApp)
	}
}

// BenchmarkServiceProviderRequires đo hiệu suất Requires method
func BenchmarkServiceProviderRequires(b *testing.B) {
	provider := NewServiceProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.Requires()
	}
}

// BenchmarkServiceProviderProviders đo hiệu suất Providers method
func BenchmarkServiceProviderProviders(b *testing.B) {
	provider := NewServiceProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.Providers()
	}
}

// BenchmarkContainerMakeLog đo hiệu suất resolve log service từ container
func BenchmarkContainerMakeLog(b *testing.B) {
	// Setup một lần
	container := di.New()
	mockApp := &diMocks.MockApplication{}
	mockApp.On("Container").Return(container).Once()

	mockConfigManager := &mocks.MockManager{}
	mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
		config := args.Get(1).(*Config)
		config.Level = "info"
		config.Console.Enabled = true
		config.Console.Colored = true
		config.File.Enabled = true
		config.File.Path = filepath.Join(os.TempDir(), "logs", "bench_make.log")
		config.File.MaxSize = 10485760
	}).Return(nil).Once()

	container.Instance("config", mockConfigManager)

	provider := NewServiceProvider()
	provider.Register(mockApp)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager, err := container.Make("log")
		if err != nil {
			b.Fatal("Failed to make log service:", err)
		}
		_ = manager
	}

	// Cleanup after benchmark
	b.StopTimer()
	if manager, err := container.Make("log"); err == nil {
		if logManager, ok := manager.(Manager); ok {
			logManager.Close()
		}
	}
}

// BenchmarkCompleteServiceProviderWorkflow đo hiệu suất toàn bộ workflow
func BenchmarkCompleteServiceProviderWorkflow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Setup cho mỗi iteration
		container := di.New()
		mockApp := &diMocks.MockApplication{}
		mockApp.On("Container").Return(container).Times(2) // Register + Boot

		mockConfigManager := &mocks.MockManager{}
		mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
			config := args.Get(1).(*Config)
			config.Level = "info"
			config.Console.Enabled = true
			config.Console.Colored = true
			config.File.Enabled = true
			config.File.Path = filepath.Join(os.TempDir(), "logs", "bench_complete.log")
			config.File.MaxSize = 10485760
		}).Return(nil).Once()

		container.Instance("config", mockConfigManager)
		provider := NewServiceProvider()
		b.StartTimer()

		// Toàn bộ workflow: Register -> Boot -> Make
		provider.Register(mockApp)
		provider.Boot(mockApp)
		manager, err := container.Make("log")
		if err != nil {
			b.Fatal("Failed to make log service:", err)
		}
		_ = manager

		// Cleanup
		b.StopTimer()
		if logManager, ok := manager.(Manager); ok {
			logManager.Close()
		}
	}
}

// BenchmarkParallelServiceProviderRegister đo hiệu suất với concurrent access
func BenchmarkParallelServiceProviderRegister(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Setup cho mỗi goroutine với clean mocks
			container := di.New()
			mockApp := &diMocks.MockApplication{}
			mockApp.On("Container").Return(container).Once()

			mockConfigManager := &mocks.MockManager{}
			mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
				config := args.Get(1).(*Config)
				config.Level = "info"
				config.Console.Enabled = true
				config.Console.Colored = true
				config.File.Enabled = true
				config.File.Path = filepath.Join(os.TempDir(), "logs", "bench_parallel.log")
				config.File.MaxSize = 10485760
			}).Return(nil).Once()

			container.Instance("config", mockConfigManager)

			provider := NewServiceProvider()
			provider.Register(mockApp)

			// Cleanup
			if manager, err := container.Make("log"); err == nil {
				if logManager, ok := manager.(Manager); ok {
					logManager.Close()
				}
			}
		}
	})
}

// BenchmarkServiceProviderWithDifferentLogLevels đo hiệu suất với các log level khác nhau
func BenchmarkServiceProviderWithDifferentLogLevels(b *testing.B) {
	logLevels := []string{"debug", "info", "warning", "error", "fatal"}

	for _, level := range logLevels {
		b.Run("Level_"+level, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				container := di.New()
				mockApp := &diMocks.MockApplication{}
				mockApp.On("Container").Return(container).Once()

				mockConfigManager := &mocks.MockManager{}
				mockConfigManager.On("UnmarshalKey", "log", mock.AnythingOfType("*log.Config")).Run(func(args mock.Arguments) {
					config := args.Get(1).(*Config)
					config.Level = level
					config.Console.Enabled = true
					config.Console.Colored = true
					config.File.Enabled = true
					config.File.Path = filepath.Join(os.TempDir(), "logs", "bench_"+level+".log")
					config.File.MaxSize = 10485760
				}).Return(nil).Once()

				container.Instance("config", mockConfigManager)
				provider := NewServiceProvider()
				b.StartTimer()

				provider.Register(mockApp)

				// Cleanup
				b.StopTimer()
				if manager, err := container.Make("log"); err == nil {
					if logManager, ok := manager.(Manager); ok {
						logManager.Close()
					}
				}
			}
		})
	}
}
