package log

// Config định nghĩa cấu hình cho log package.
//
// Struct này chứa các thiết lập cần thiết để khởi tạo và cấu hình
// logging system bao gồm level và các handler configurations.
type Config struct {
	// Level xác định mức độ log tối thiểu sẽ được ghi.
	// Các giá trị hợp lệ: debug, info, warning, error, fatal
	Level string `yaml:"level" json:"level"`

	// Console cấu hình cho console handler
	Console ConsoleConfig `yaml:"console" json:"console"`

	// File cấu hình cho file handler
	File FileConfig `yaml:"file" json:"file"`

	// Stack cấu hình cho stack handler
	Stack StackConfig `yaml:"stack" json:"stack"`
}

// ConsoleConfig định nghĩa cấu hình cho console handler.
type ConsoleConfig struct {
	// Enabled bật/tắt console handler
	Enabled bool `yaml:"enabled" json:"enabled"`

	// Colored bật/tắt màu sắc cho console output
	Colored bool `yaml:"colored" json:"colored"`
}

// FileConfig định nghĩa cấu hình cho file handler.
type FileConfig struct {
	// Enabled bật/tắt file handler
	Enabled bool `yaml:"enabled" json:"enabled"`

	// Path đường dẫn file log
	Path string `yaml:"path" json:"path"`

	// MaxSize kích thước tối đa của file log (bytes) trước khi rotate
	// 0 = không giới hạn
	MaxSize int64 `yaml:"max_size" json:"max_size"`
}

// StackConfig định nghĩa cấu hình cho stack handler.
type StackConfig struct {
	// Enabled bật/tắt stack handler
	Enabled bool `yaml:"enabled" json:"enabled"`

	// Handlers cấu hình các sub-handlers
	Handlers StackHandlers `yaml:"handlers" json:"handlers"`
}

// StackHandlers định nghĩa cấu hình các handler trong stack.
type StackHandlers struct {
	// Console bật/tắt console handler trong stack
	Console bool `yaml:"console" json:"console"`

	// File bật/tắt file handler trong stack
	File bool `yaml:"file" json:"file"`
}

// DefaultConfig trả về cấu hình mặc định cho log package.
//
// Cấu hình mặc định sử dụng:
//   - Level: info
//   - Console handler được bật với màu sắc
//   - File và Stack handlers được tắt
//
// Trả về:
//   - *Config: Cấu hình mặc định
func DefaultConfig() *Config {
	return &Config{
		Level: "info",
		Console: ConsoleConfig{
			Enabled: true,
			Colored: true,
		},
		File: FileConfig{
			Enabled: false,
			Path:    "",
			MaxSize: 0,
		},
		Stack: StackConfig{
			Enabled: false,
			Handlers: StackHandlers{
				Console: false,
				File:    false,
			},
		},
	}
}

// Validate kiểm tra tính hợp lệ của cấu hình.
//
// Phương thức này xác minh:
//   - Level có hợp lệ không
//   - Các handler có được cấu hình đúng không
//   - File handler có path hợp lệ không khi được bật
//   - Stack handler có ít nhất một sub-handler được bật không
//
// Trả về:
//   - error: Lỗi nếu cấu hình không hợp lệ
func (c *Config) Validate() error {
	// Kiểm tra level hợp lệ
	validLevels := map[string]bool{
		"debug":   true,
		"info":    true,
		"warning": true,
		"error":   true,
		"fatal":   true,
	}

	if !validLevels[c.Level] {
		return &ConfigError{
			Field:   "level",
			Value:   c.Level,
			Message: "invalid log level, must be one of: debug, info, warning, error, fatal",
		}
	}

	// Kiểm tra có ít nhất một handler được bật
	if !c.Console.Enabled && !c.File.Enabled && !c.Stack.Enabled {
		return &ConfigError{
			Field:   "handlers",
			Message: "at least one handler must be enabled",
		}
	}

	// Validate file handler nếu được bật
	if c.File.Enabled {
		if c.File.Path == "" {
			return &ConfigError{
				Field:   "file.path",
				Message: "path is required when file handler is enabled",
			}
		}
		if c.File.MaxSize < 0 {
			return &ConfigError{
				Field:   "file.max_size",
				Value:   string(rune(c.File.MaxSize)),
				Message: "max_size must be non-negative (0 for unlimited)",
			}
		}
	}

	// Validate stack handler nếu được bật
	if c.Stack.Enabled {
		if !c.Stack.Handlers.Console && !c.Stack.Handlers.File {
			return &ConfigError{
				Field:   "stack.handlers",
				Message: "stack handler must have at least one sub-handler enabled",
			}
		}

		// Nếu stack có file handler được bật, file handler chính cũng phải có cấu hình hợp lệ
		if c.Stack.Handlers.File {
			if c.File.Path == "" {
				return &ConfigError{
					Field:   "file.path",
					Message: "path is required when file handler is used in stack",
				}
			}
		}
	}

	return nil
}

// ConfigError represent lỗi cấu hình log.
//
// Error type này cung cấp thông tin chi tiết về lỗi cấu hình
// bao gồm field nào bị lỗi và lý do.
type ConfigError struct {
	Field   string
	Value   string
	Message string
}

// Error implement error interface.
//
// Trả về:
//   - string: Error message với thông tin chi tiết
func (e *ConfigError) Error() string {
	if e.Value != "" {
		return "log config error in field '" + e.Field + "' with value '" + e.Value + "': " + e.Message
	}
	return "log config error in field '" + e.Field + "': " + e.Message
}
