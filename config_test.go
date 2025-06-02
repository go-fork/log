package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectedErr string
	}{
		{
			name:        "Valid default config",
			config:      DefaultConfig(),
			expectedErr: "",
		},
		{
			name: "Invalid log level",
			config: &Config{
				Level: "unknown",
				Console: ConsoleConfig{
					Enabled: true,
				},
			},
			expectedErr: "invalid log level",
		},
		{
			name: "No handlers enabled",
			config: &Config{
				Level: "info",
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
			name: "File handler enabled without path",
			config: &Config{
				Level: "info",
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: true,
					Path:    "",
				},
				Stack: StackConfig{
					Enabled: false,
				},
			},
			expectedErr: "path is required when file handler is enabled",
		},
		{
			name: "File handler with negative max size",
			config: &Config{
				Level: "info",
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
			name: "Stack handler enabled without sub-handlers",
			config: &Config{
				Level: "info",
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: false,
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
			name: "Stack handler with file sub-handler but no file path",
			config: &Config{
				Level: "info",
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: false,
					Path:    "",
				},
				Stack: StackConfig{
					Enabled: true,
					Handlers: StackHandlers{
						Console: false,
						File:    true,
					},
				},
			},
			expectedErr: "path is required when file handler is used in stack",
		},
		{
			name: "Valid config with all features enabled",
			config: &Config{
				Level: "debug",
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

func TestConfigDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "info", config.Level)
	assert.True(t, config.Console.Enabled)
	assert.True(t, config.Console.Colored)
	assert.False(t, config.File.Enabled)
	assert.Empty(t, config.File.Path)
	assert.Equal(t, int64(0), config.File.MaxSize)
	assert.False(t, config.Stack.Enabled)
	assert.False(t, config.Stack.Handlers.Console)
	assert.False(t, config.Stack.Handlers.File)
}

func TestConfigErrorWithValue(t *testing.T) {
	err := &ConfigError{
		Field:   "level",
		Value:   "unknown",
		Message: "invalid log level, must be one of: debug, info, warning, error, fatal",
	}

	expected := "log config error in field 'level' with value 'unknown': invalid log level, must be one of: debug, info, warning, error, fatal"
	assert.Equal(t, expected, err.Error())
}

func TestConfigErrorWithoutValue(t *testing.T) {
	err := &ConfigError{
		Field:   "handlers",
		Message: "at least one handler must be enabled",
	}

	expected := "log config error in field 'handlers': at least one handler must be enabled"
	assert.Equal(t, expected, err.Error())
}

func TestValidateAllLogLevels(t *testing.T) {
	validLevels := []string{"debug", "info", "warning", "error", "fatal"}

	for _, level := range validLevels {
		t.Run("Valid level: "+level, func(t *testing.T) {
			config := &Config{
				Level: level,
				Console: ConsoleConfig{
					Enabled: true,
				},
			}

			err := config.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestValidateFileHandlerConfigurations(t *testing.T) {
	tests := []struct {
		name        string
		maxSize     int64
		expectedErr bool
	}{
		{
			name:        "Zero max size (unlimited)",
			maxSize:     0,
			expectedErr: false,
		},
		{
			name:        "Positive max size",
			maxSize:     1024 * 1024, // 1MB
			expectedErr: false,
		},
		{
			name:        "Negative max size",
			maxSize:     -100,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Level: "info",
				Console: ConsoleConfig{
					Enabled: false,
				},
				File: FileConfig{
					Enabled: true,
					Path:    "/tmp/logs/test.log",
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
func BenchmarkDefaultConfig(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := DefaultConfig()
		_ = config
	}
}

// BenchmarkConfigValidateValid đo hiệu suất validate cấu hình hợp lệ
func BenchmarkConfigValidateValid(b *testing.B) {
	config := DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := config.Validate()
		_ = err
	}
}

// BenchmarkConfigValidateInvalidLevel đo hiệu suất validate với level không hợp lệ
func BenchmarkConfigValidateInvalidLevel(b *testing.B) {
	config := &Config{
		Level: "invalid_level",
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
func BenchmarkConfigValidateNoHandlers(b *testing.B) {
	config := &Config{
		Level: "info",
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
func BenchmarkConfigValidateFileHandler(b *testing.B) {
	config := &Config{
		Level: "info",
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
func BenchmarkConfigValidateStackHandler(b *testing.B) {
	config := &Config{
		Level: "info",
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
func BenchmarkConfigValidateComplexConfig(b *testing.B) {
	config := &Config{
		Level: "debug",
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
func BenchmarkConfigErrorCreation(b *testing.B) {
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
func BenchmarkConfigErrorString(b *testing.B) {
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
func BenchmarkConfigValidateWithDifferentLevels(b *testing.B) {
	levels := []string{"debug", "info", "warning", "error", "fatal"}

	for _, level := range levels {
		b.Run("Level_"+level, func(b *testing.B) {
			config := &Config{
				Level: level,
				Console: ConsoleConfig{
					Enabled: true,
					Colored: true,
				},
				File: FileConfig{
					Enabled: true,
					Path:    "/tmp/logs/bench_" + level + ".log",
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
func BenchmarkConfigValidateParallel(b *testing.B) {
	config := DefaultConfig()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := config.Validate()
			_ = err
		}
	})
}

// BenchmarkConfigCreationAndValidation đo hiệu suất tạo và validate config
func BenchmarkConfigCreationAndValidation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := &Config{
			Level: "info",
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
func BenchmarkConfigValidateMemoryUsage(b *testing.B) {
	b.ReportAllocs()

	config := &Config{
		Level: "info",
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
func BenchmarkConfigValidateWorstCase(b *testing.B) {
	// Worst case: invalid config that triggers all validation checks
	config := &Config{
		Level: "invalid_level",
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
func BenchmarkConfigValidateFilePathVariations(b *testing.B) {
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
				Level: "info",
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
func BenchmarkConfigValidateMaxSizeVariations(b *testing.B) {
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
			sizeName = "Unlimited"
		case maxSize < 1048576:
			sizeName = "KB"
		case maxSize < 1073741824:
			sizeName = "MB"
		default:
			sizeName = "GB"
		}

		b.Run("MaxSize_"+sizeName, func(b *testing.B) {
			config := &Config{
				Level: "info",
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
