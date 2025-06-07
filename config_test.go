package log

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.fork.vn/log/handler"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectedErr string
	}{
		{
			name:        "valid_default_config",
			config:      DefaultConfig(),
			expectedErr: "",
		},
		{
			name: "invalid_log_level",
			config: &Config{
				Level: handler.Level(99), // Invalid level
				Console: ConsoleConfig{
					Enabled: true,
				},
			},
			expectedErr: "invalid log level",
		},
		{
			name: "no_handlers_enabled",
			config: &Config{
				Level: handler.InfoLevel,
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: false,
				},
				Stack: StackConfig{
					Enabled: false,
				},
			},
			expectedErr: "at least one handler must be enabled",
		},
		{
			name: "file_handler_enabled_without_path",
			config: &Config{
				Level: handler.InfoLevel,
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: true,
					Path:    "", // Để trống để test validation
				},
				Stack: StackConfig{
					Enabled: false,
				},
			},
			expectedErr: "path is required for file handler initialization",
		},
		{
			name: "file_handler_with_negative_max_size",
			config: &Config{
				Level: handler.InfoLevel,
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: true,
					Path:    "/tmp/logs",
					MaxSize: -1,
				},
				Stack: StackConfig{
					Enabled: false,
				},
			},
			expectedErr: "max_size must be non-negative",
		},
		{
			name: "stack_handler_enabled_without_sub_handlers",
			config: &Config{
				Level: handler.InfoLevel,
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: false,
					Path:    "/tmp/logs", // Cần có path hợp lệ vì luôn được kiểm tra
				},
				Stack: StackConfig{
					Enabled: true,
					Handlers: StackHandlers{
						Console: false,
						File:    false,
					},
				},
			},
			expectedErr: "stack handler must have at least one sub-handler enabled",
		},
		{
			name: "stack_handler_with_file_sub_handler_but_no_file_path",
			config: &Config{
				Level: handler.InfoLevel,
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: false,
					Path:    "", // Để trống để test validation
				},
				Stack: StackConfig{
					Enabled: true,
					Handlers: StackHandlers{
						Console: false,
						File:    true,
					},
				},
			},
			expectedErr: "path is required for file handler initialization",
		},
		{
			name: "valid_config_with_all_features_enabled",
			config: &Config{
				Level: handler.DebugLevel,
				Console: ConsoleConfig{
					Enabled: true,
					Colored: true,
				},
				File: FileConfig{
					Enabled: true,
					Path:    "/tmp/logs/app.log",
					MaxSize: 10485760, // 10MB
				},
				Stack: StackConfig{
					Enabled: true,
					Handlers: StackHandlers{
						Console: true,
						File:    true,
					},
				},
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip test nếu directory không tồn tại (thay vì tạo directory)
			if tt.config.File.Path != "" {
				dir := filepath.Dir(tt.config.File.Path)
				if dir != "." {
					if _, err := os.Stat(dir); os.IsNotExist(err) {
						t.Skipf("Skipping test because directory does not exist: %s", dir)
						return
					}
				}
			}

			err := tt.config.Validate()

			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestConfig_DefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, handler.InfoLevel, config.Level)
	assert.True(t, config.Console.Enabled)
	assert.True(t, config.Console.Colored)
	assert.False(t, config.File.Enabled) // Mặc định File.Enabled = false
	assert.Equal(t, "", config.File.Path) // Default path is empty
	assert.Equal(t, int64(10*1024*1024), config.File.MaxSize) // 10MB
	assert.False(t, config.Stack.Enabled)
	assert.False(t, config.Stack.Handlers.Console)
	assert.False(t, config.Stack.Handlers.File)
}

func TestConfigError_WithValue(t *testing.T) {
	err := &ConfigError{
		Field:   "level",
		Value:   "unknown",
		Message: "invalid log level, must be one of: debug, info, warning, error, fatal",
	}

	expected := "log config error in field 'level' with value 'unknown': invalid log level, must be one of: debug, info, warning, error, fatal"
	assert.Equal(t, expected, err.Error())
}

func TestConfigError_WithoutValue(t *testing.T) {
	err := &ConfigError{
		Field:   "handlers",
		Message: "at least one handler must be enabled",
	}

	expected := "log config error in field 'handlers': at least one handler must be enabled"
	assert.Equal(t, expected, err.Error())
}

func TestConfig_ValidateAllLogLevels(t *testing.T) {
	validLevels := []handler.Level{
		handler.DebugLevel,
		handler.InfoLevel,
		handler.WarningLevel,
		handler.ErrorLevel,
		handler.FatalLevel,
	}

	for _, level := range validLevels {
		t.Run("Valid level: "+level.String(), func(t *testing.T) {
			config := &Config{
				Level: level,
				Console: ConsoleConfig{
					Enabled: true,
				},
				File: FileConfig{
					Enabled: false,
					Path:    "", // Không set path để tránh validation
				},
			}

			err := config.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestConfig_ValidateFileHandlerConfigurations(t *testing.T) {
	tests := []struct {
		name        string
		maxSize     int64
		expectedErr bool
	}{
		{
			name:        "zero_max_size_unlimited",
			maxSize:     0,
			expectedErr: false,
		},
		{
			name:        "positive_max_size",
			maxSize:     1024 * 1024, // 1MB
			expectedErr: false,
		},
		{
			name:        "negative_max_size",
			maxSize:     -100,
			expectedErr: true,
		},
	}

	// Tạo temporary directory cho test
	tempDir := t.TempDir()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Level: handler.InfoLevel,
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: true,
					Path:    filepath.Join(tempDir, "test.log"),
					MaxSize: tt.maxSize,
				},
			}

			err := config.Validate()
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// BenchmarkDefaultConfig đo hiệu suất tạo cấu hình mặc định
func BenchmarkConfig_DefaultConfig(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := DefaultConfig()
		_ = config
	}
}

// BenchmarkConfigValidateValid đo hiệu suất validate cấu hình hợp lệ
func BenchmarkConfig_Validate_Valid(b *testing.B) {
	config := DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := config.Validate()
		_ = err
	}
}

// BenchmarkConfigValidateInvalidLevel đo hiệu suất validate với level không hợp lệ
func BenchmarkConfig_Validate_InvalidLevel(b *testing.B) {
	config := &Config{
		Level: handler.Level(99), // Invalid level
		Console: ConsoleConfig{
			Enabled: true,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := config.Validate()
		_ = err
	}
}

// BenchmarkConfigValidateNoHandlers đo hiệu suất validate khi không có handler
func BenchmarkConfig_Validate_NoHandlers(b *testing.B) {
	config := &Config{
		Level: handler.InfoLevel,
		Console: ConsoleConfig{
			Enabled: false,
		},
		File: FileConfig{
			Enabled: false,
		},
		Stack: StackConfig{
			Enabled: false,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := config.Validate()
		_ = err
	}
}

// BenchmarkConfigValidateFileHandler đo hiệu suất validate file handler
func BenchmarkConfig_Validate_FileHandler(b *testing.B) {
	config := &Config{
		Level: handler.InfoLevel,
		Console: ConsoleConfig{
			Enabled: false,
		},
		File: FileConfig{
			Enabled: true,
			Path:    "/tmp/logs/bench.log",
			MaxSize: 10485760,
		},
		Stack: StackConfig{
			Enabled: false,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := config.Validate()
		_ = err
	}
}

// BenchmarkConfigValidateStackHandler đo hiệu suất validate stack handler
func BenchmarkConfig_Validate_StackHandler(b *testing.B) {
	config := &Config{
		Level: handler.InfoLevel,
		Console: ConsoleConfig{
			Enabled: false,
		},
		File: FileConfig{
			Enabled: false,
			Path:    "/tmp/logs/bench.log",
		},
		Stack: StackConfig{
			Enabled: true,
			Handlers: StackHandlers{
				Console: true,
				File:    true,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := config.Validate()
		_ = err
	}
}

// BenchmarkConfigValidateComplexConfig đo hiệu suất validate cấu hình phức tạp
func BenchmarkConfig_Validate_ComplexConfig(b *testing.B) {
	config := &Config{
		Level: handler.DebugLevel,
		Console: ConsoleConfig{
			Enabled: true,
			Colored: true,
		},
		File: FileConfig{
			Enabled: true,
			Path:    "/tmp/logs/complex_bench.log",
			MaxSize: 20971520, // 20MB
		},
		Stack: StackConfig{
			Enabled: true,
			Handlers: StackHandlers{
				Console: true,
				File:    true,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := config.Validate()
		_ = err
	}
}

// BenchmarkConfigErrorCreation đo hiệu suất tạo ConfigError
func BenchmarkConfigError_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := &ConfigError{
			Field:   "level",
			Value:   "invalid",
			Message: "invalid log level",
		}
		// Use the error to avoid unused write warnings
		_ = err.Field
		_ = err.Value
		_ = err.Message
	}
}

// BenchmarkConfigErrorString đo hiệu suất chuyển ConfigError thành string
func BenchmarkConfigError_String(b *testing.B) {
	err := &ConfigError{
		Field:   "level",
		Value:   "invalid",
		Message: "invalid log level, must be one of: debug, info, warning, error, fatal",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str := err.Error()
		_ = str
	}
}

// BenchmarkConfigValidateWithDifferentLevels đo hiệu suất validate với các level khác nhau
func BenchmarkConfig_Validate_WithDifferentLevels(b *testing.B) {
	levels := []handler.Level{
		handler.DebugLevel,
		handler.InfoLevel,
		handler.WarningLevel,
		handler.ErrorLevel,
		handler.FatalLevel,
	}

	for _, level := range levels {
		b.Run("Level_"+level.String(), func(b *testing.B) {
			config := &Config{
				Level: level,
				Console: ConsoleConfig{
					Enabled: true,
					Colored: true,
				},
				File: FileConfig{
					Enabled: true,
					Path:    "/tmp/logs/bench_" + level.String() + ".log",
					MaxSize: 10485760,
				},
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := config.Validate()
				_ = err
			}
		})
	}
}

// BenchmarkConfigValidateParallel đo hiệu suất validate parallel
func BenchmarkConfig_Validate_Parallel(b *testing.B) {
	config := DefaultConfig()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := config.Validate()
			_ = err
		}
	})
}

// BenchmarkConfigCreationAndValidation đo hiệu suất tạo và validate config
func BenchmarkConfig_CreationAndValidation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := &Config{
			Level: handler.InfoLevel,
			Console: ConsoleConfig{
				Enabled: true,
				Colored: true,
			},
			File: FileConfig{
				Enabled: true,
				Path:    "/tmp/logs/creation_bench.log",
				MaxSize: 5242880, // 5MB
			},
			Stack: StackConfig{
				Enabled: false,
			},
		}

		err := config.Validate()
		_ = err
	}
}

// BenchmarkConfigValidateMemoryUsage đo memory usage của validation
func BenchmarkConfig_Validate_MemoryUsage(b *testing.B) {
	b.ReportAllocs()

	config := &Config{
		Level: handler.InfoLevel,
		Console: ConsoleConfig{
			Enabled: true,
			Colored: true,
		},
		File: FileConfig{
			Enabled: true,
			Path:    "/tmp/logs/memory_bench.log",
			MaxSize: 10485760,
		},
		Stack: StackConfig{
			Enabled: true,
			Handlers: StackHandlers{
				Console: true,
				File:    true,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := config.Validate()
		_ = err
	}
}

// BenchmarkConfigValidateWorstCase đo hiệu suất validate worst case scenario
func BenchmarkConfig_Validate_WorstCase(b *testing.B) {
	// Worst case: invalid config that triggers all validation checks
	config := &Config{
		Level: handler.Level(99), // Invalid level
		Console: ConsoleConfig{
			Enabled: false,
		},
		File: FileConfig{
			Enabled: true,
			Path:    "", // Invalid: empty path
			MaxSize: -1, // Invalid: negative size
		},
		Stack: StackConfig{
			Enabled: true,
			Handlers: StackHandlers{
				Console: false,
				File:    false, // Invalid: no sub-handlers
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := config.Validate()
		_ = err
	}
}

// BenchmarkConfigValidateFilePathVariations đo hiệu suất với các đường dẫn file khác nhau
func BenchmarkConfig_Validate_FilePathVariations(b *testing.B) {
	paths := []string{
		"/tmp/logs/short.log",
		"/very/long/path/to/logs/directory/with/multiple/levels/application.log",
		"/tmp/logs/unicode_测试_файл.log",
		"/tmp/logs/numbers_123456789.log",
		"/tmp/logs/special-chars_@#$%^&*()_+.log",
	}

	for i, path := range paths {
		b.Run("Path_"+string(rune('A'+i)), func(b *testing.B) {
			config := &Config{
				Level: handler.InfoLevel,
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: true,
					Path:    path,
					MaxSize: 10485760,
				},
			}

			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				err := config.Validate()
				_ = err
			}
		})
	}
}

// BenchmarkConfigValidateMaxSizeVariations đo hiệu suất với các max size khác nhau
func BenchmarkConfig_Validate_MaxSizeVariations(b *testing.B) {
	maxSizes := []int64{
		0,          // Unlimited
		1024,       // 1KB
		1048576,    // 1MB
		10485760,   // 10MB
		104857600,  // 100MB
		1073741824, // 1GB
	}

	for _, maxSize := range maxSizes {
		sizeName := ""
		switch {
		case maxSize == 0:
			sizeName = "unlimited"
		case maxSize < 1048576:
			sizeName = "kb"
		case maxSize < 1073741824:
			sizeName = "mb"
		default:
			sizeName = "gb"
		}

		b.Run("max_size_"+sizeName, func(b *testing.B) {
			config := &Config{
				Level: handler.InfoLevel,
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: true,
					Path:    "/tmp/logs/size_bench.log",
					MaxSize: maxSize,
				},
			}

			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				err := config.Validate()
				_ = err
			}
		})
	}
}
