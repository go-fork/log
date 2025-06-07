package log

import (
	"fmt"
	"os"
	"path/filepath"

	"go.fork.vn/log/handler"
)

// Config định nghĩa cấu hình cho log package.
//
// Struct này chứa các thiết lập cần thiết để khởi tạo và cấu hình
// logging system bao gồm level và các handler configurations.
type Config struct {
	// Level xác định mức độ log tối thiểu sẽ được ghi.
	// Các giá trị hợp lệ: DebugLevel, InfoLevel, WarningLevel, ErrorLevel, FatalLevel
	Level handler.Level `mapstructure:"level" yaml:"level" json:"level"`

	// Console cấu hình cho console handler
	Console ConsoleConfig `mapstructure:"console" yaml:"console" json:"console"`

	// File cấu hình cho file handler
	File FileConfig `mapstructure:"file" yaml:"file" json:"file"`

	// Stack cấu hình cho stack handler
	Stack StackConfig `mapstructure:"stack" yaml:"stack" json:"stack"`
}

// ConsoleConfig định nghĩa cấu hình cho console handler.
type ConsoleConfig struct {
	// Enabled bật/tắt console handler
	Enabled bool `mapstructure:"enabled" yaml:"enabled" json:"enabled"`

	// Colored bật/tắt màu sắc cho console output
	Colored bool `mapstructure:"colored" yaml:"colored" json:"colored"`
}

// FileConfig định nghĩa cấu hình cho file handler.
type FileConfig struct {
	// Enabled bật/tắt file handler
	Enabled bool `mapstructure:"enabled" yaml:"enabled" json:"enabled"`

	// Path đường dẫn file log
	Path string `mapstructure:"path" yaml:"path" json:"path"`

	// MaxSize kích thước tối đa của file log (bytes) trước khi rotate
	// 0 = không giới hạn
	MaxSize int64 `mapstructure:"max_size" yaml:"max_size" json:"max_size"`
}

// StackConfig định nghĩa cấu hình cho stack handler.
type StackConfig struct {
	// Enabled bật/tắt stack handler
	Enabled bool `mapstructure:"enabled" yaml:"enabled" json:"enabled"`

	// Handlers cấu hình các sub-handlers
	Handlers StackHandlers `mapstructure:"handlers" yaml:"handlers" json:"handlers"`
}

// StackHandlers định nghĩa cấu hình các handler trong stack.
type StackHandlers struct {
	// Console bật/tắt console handler trong stack
	Console bool `mapstructure:"console" yaml:"console" json:"console"`

	// File bật/tắt file handler trong stack
	File bool `mapstructure:"file" yaml:"file" json:"file"`
}

// DefaultConfig trả về cấu hình mặc định cho log package.
//
// Cấu hình mặc định sử dụng:
//   - Level: InfoLevel
//   - Console handler được bật với màu sắc
//   - File handler được tắt (path rỗng)
//   - Stack handler được tắt
//
// Trả về:
//   - *Config: Cấu hình mặc định
func DefaultConfig() *Config {
	return &Config{
		Level: handler.InfoLevel,
		Console: ConsoleConfig{
			Enabled: true,
			Colored: true,
		},
		File: FileConfig{
			Enabled: false,
			Path:    "",               // Empty path - user must set this explicitly
			MaxSize: 10 * 1024 * 1024, // 10MB
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
//   - Thư mục log có tồn tại và có quyền ghi không
//   - Stack handler có ít nhất một sub-handler được bật không
//
// Trả về:
//   - error: Lỗi nếu cấu hình không hợp lệ
func (c *Config) Validate() error {
	// Kiểm tra level hợp lệ
	validLevels := map[handler.Level]bool{
		handler.DebugLevel:   true,
		handler.InfoLevel:    true,
		handler.WarningLevel: true,
		handler.ErrorLevel:   true,
		handler.FatalLevel:   true,
	}

	if !validLevels[c.Level] {
		return &ConfigError{
			Field:   "level",
			Value:   c.Level.String(),
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

	// Validate file handler path - chỉ yêu cầu khi file handler được sử dụng
	needsFilePath := c.File.Enabled || (c.Stack.Enabled && c.Stack.Handlers.File)
	if needsFilePath && c.File.Path == "" {
		return &ConfigError{
			Field:   "file.path",
			Message: "path is required for file handler initialization",
		}
	}

	// Kiểm tra thư mục log nếu có path
	if c.File.Path != "" {
		if err := c.validateAndCreateLogDir(c.File.Path); err != nil {
			return &ConfigError{
				Field:   "file.path",
				Value:   c.File.Path,
				Message: "log directory validation failed: " + err.Error(),
			}
		}
	}

	if c.File.MaxSize < 0 {
		return &ConfigError{
			Field:   "file.max_size",
			Value:   string(rune(c.File.MaxSize)),
			Message: "max_size must be non-negative (0 for unlimited)",
		}
	}

	// Validate file handler - luôn validate path nếu có
	// (không phụ thuộc vào File.Enabled vì chúng ta luôn cần validate)

	// Validate stack handler nếu được bật
	if c.Stack.Enabled {
		if !c.Stack.Handlers.Console && !c.Stack.Handlers.File {
			return &ConfigError{
				Field:   "stack.handlers",
				Message: "stack handler must have at least one sub-handler enabled",
			}
		}

		// File.Path đã được kiểm tra ở trên, không cần kiểm tra lại
	}

	return nil
}

// validateAndCreateLogDir kiểm tra thư mục log có tồn tại và có quyền ghi không.
//
// Phương thức này:
//   - Kiểm tra thư mục cha của file log có tồn tại không
//   - Kiểm tra quyền ghi vào thư mục
//   - KHÔNG tự động tạo thư mục
//
// Tham số:
//   - logPath: Đường dẫn file log
//
// Trả về:
//   - error: Lỗi nếu thư mục không tồn tại hoặc không có quyền ghi
func (c *Config) validateAndCreateLogDir(logPath string) error {
	// Lấy thư mục cha của file log
	logDir := filepath.Dir(logPath)

	// Kiểm tra thư mục có tồn tại không
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return fmt.Errorf("path to folder do not exists: %s", logDir)
	} else if err != nil {
		// Lỗi khác khi stat thư mục
		return fmt.Errorf("cannot access directory: %w", err)
	}

	// Kiểm tra quyền ghi bằng cách tạo file tạm thời
	testFile := filepath.Join(logDir, ".log_write_test")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("directory does not have write permission: %s", logDir)
	}
	file.Close()
	os.Remove(testFile)

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
