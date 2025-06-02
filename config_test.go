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
