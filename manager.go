package log

import (
	"fmt"
	"sync"

	"go.fork.vn/log/handler"
)

// Manager định nghĩa interface cho hệ thống quản lý handler tập trung.
//
// Interface Manager cung cấp các method để quản lý handlers được chia sẻ
// giữa nhiều logger instances. Điều này đảm bảo rằng các handlers như file,
// console, stack được sử dụng chung và không bị duplicate.
type Manager interface {
	// AddHandler đăng ký một handler mới vào manager.
	//
	// Tham số:
	//   - handlerType: HandlerType - loại handler (console, file, stack)
	//   - handler: handler.Handler - instance của handler cần thêm
	AddHandler(handlerType HandlerType, handler handler.Handler)

	// RemoveHandler hủy đăng ký và đóng một handler.
	//
	// Tham số:
	//   - handlerType: HandlerType - loại handler cần xóa
	RemoveHandler(handlerType HandlerType)

	// GetHandler trả về một handler đã đăng ký theo loại.
	//
	// Tham số:
	//   - handlerType: HandlerType - loại handler cần lấy
	//
	// Trả về:
	//   - handler.Handler - instance của handler hoặc nil nếu không tìm thấy
	GetHandler(handlerType HandlerType) handler.Handler

	// SetHandler thiết lập handler cho logger cụ thể.
	//
	// Tham số:
	//   - loggerContext: string - context của logger
	//   - handlerType: HandlerType - loại handler cần thiết lập
	SetHandler(loggerContext string, handlerType HandlerType)

	// GetLogger trả về logger theo context, tự động tạo mới nếu chưa tồn tại.
	//
	// Method này hoạt động như getOrCreate pattern - nếu logger với context
	// đã tồn tại thì trả về, nếu chưa thì tạo mới và trả về.
	//
	// Tham số:
	//   - context: string - context để xác định nguồn gốc log
	//
	// Trả về:
	//   - Logger: logger cho context đã cho (existing hoặc newly created)
	//
	// Ví dụ:
	//
	//	userLogger := manager.GetLogger("UserService")  // tạo mới nếu chưa có
	//	userLogger2 := manager.GetLogger("UserService") // trả về cái đã tồn tại
	GetLogger(context string) Logger

	// Close đóng tất cả các handlers và giải phóng tài nguyên.
	//
	// Trả về:
	//   - error: một lỗi nếu việc đóng handlers thất bại
	Close() error
}

// manager là triển khai chuẩn của interface Manager.
//
// manager cung cấp quản lý tập trung cho các handlers được chia sẻ
// giữa nhiều logger instances. Nó đảm bảo thread-safety và quản lý
// vòng đời của các handlers và loggers.
//
// Tính năng:
//   - Quản lý handlers thread-safe bằng RWMutex
//   - Chia sẻ handlers giữa nhiều loggers
//   - Quản lý vòng đời handlers (tạo/xóa/đóng)
//   - Thiết lập cấp độ log toàn cục
//   - Quản lý danh sách loggers đã tạo
type manager struct {
	config   *Config                         // Cấu hình manager
	handlers map[HandlerType]handler.Handler // Map các handlers theo loại
	loggers  map[string]Logger               // Map các loggers đã tạo theo context
	mu       sync.RWMutex                    // Mutex để đảm bảo thread-safety
}

// NewManager tạo và trả về một instance manager mới với cấu hình được chỉ định.
//
// Hàm này khởi tạo một manager với cấu hình được cung cấp. Config là bắt buộc
// và phải được cung cấp để xác định handlers nào sẽ được khởi tạo.
//
// Tham số:
//   - config: *Config - cấu hình cho manager (bắt buộc, không thể nil)
//
// Trả về:
//   - Manager: một instance mới của manager triển khai interface Manager.
//
// Ví dụ:
//
//	config := &log.Config{
//		Level: "info",
//		Console: log.ConsoleConfig{Enabled: true, Colored: true},
//	}
//	manager := log.NewManager(config)
//	logger := manager.GetLogger("UserService")
func NewManager(config *Config) Manager {
	if config == nil {
		panic("config cannot be nil")
	}

	m := &manager{
		config:   config,
		handlers: make(map[HandlerType]handler.Handler),
		loggers:  make(map[string]Logger),
	}

	// Khởi tạo handlers theo cấu hình
	m.initializeHandlers()

	return m
}

// AddHandler thêm một handler mới vào manager.
//
// Method này đăng ký một handler với loại đã cho. Nếu một handler với cùng loại
// đã tồn tại, nó sẽ bị thay thế và handler cũ sẽ được đóng. Method này là thread-safe.
// Handler mới cũng sẽ được thêm vào tất cả loggers đã tồn tại.
//
// Tham số:
//   - handlerType: HandlerType - loại handler (console, file, stack)
//   - handler: handler.Handler - triển khai handler cần thêm
//
// Ví dụ:
//
//	// Thêm một file handler
//	fileHandler, _ := handler.NewFileHandler("app.log", 10*1024*1024)
//	manager.AddHandler(HandlerTypeFile, fileHandler)
func (m *manager) AddHandler(handlerType HandlerType, handler handler.Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Nếu handler cũ cùng loại tồn tại, đóng lại để tránh leak resource
	if old, ok := m.handlers[handlerType]; ok {
		old.Close()
	}
	m.handlers[handlerType] = handler

	// Thêm handler vào tất cả loggers đã tồn tại
	for _, logger := range m.loggers {
		logger.AddHandler(handlerType, handler)
	}
}

// RemoveHandler xóa một handler khỏi manager theo loại.
//
// Handler sẽ được đóng đúng cách trước khi xóa để đảm bảo tất cả các tài nguyên
// được giải phóng. Method này là thread-safe.
//
// Tham số:
//   - handlerType: HandlerType - loại handler cần xóa
//
// Nếu handler được chỉ định không tồn tại, thao tác này không làm gì.
//
// Ví dụ:
//
//	manager.RemoveHandler(HandlerTypeFile) // Xóa và đóng file handler
func (m *manager) RemoveHandler(handlerType HandlerType) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Đóng và xóa handler nếu nó tồn tại
	if handler, ok := m.handlers[handlerType]; ok {
		handler.Close()
		delete(m.handlers, handlerType)

		// Xóa handler khỏi tất cả loggers đã tồn tại
		for _, logger := range m.loggers {
			logger.RemoveHandler(handlerType)
		}
	}
}

// GetHandler trả về một handler đã đăng ký theo loại.
//
// Method này trả về một handler theo loại đã cho hoặc nil nếu không tìm thấy.
// Method này là thread-safe.
//
// Tham số:
//   - handlerType: HandlerType - loại handler cần lấy
//
// Trả về:
//   - handler.Handler: instance của handler hoặc nil nếu không tìm thấy
//
// Ví dụ:
//
//	if h := manager.GetHandler(HandlerTypeFile); h != nil {
//	    // Sử dụng handler
//	}
func (m *manager) GetHandler(handlerType HandlerType) handler.Handler {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.handlers[handlerType]
}

// SetHandler thiết lập handler cho logger cụ thể.
//
// Method này thiết lập handler cho logger đã có sẵn.
// Method này là thread-safe.
//
// Tham số:
//   - loggerContext: string - context của logger
//   - handlerType: HandlerType - loại handler cần thiết lập
//
// Ví dụ:
//
//	manager.SetHandler("UserService", HandlerTypeFile)
func (m *manager) SetHandler(loggerContext string, handlerType HandlerType) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Tìm logger theo context
	if logger, exists := m.loggers[loggerContext]; exists {
		// Tìm handler theo loại
		if handler, ok := m.handlers[handlerType]; ok {
			logger.AddHandler(handlerType, handler)
		}
	}
}

// GetLogger trả về logger theo context, tự động tạo mới nếu chưa tồn tại.
//
// Method này hoạt động như getOrCreate pattern - nếu logger với context
// đã tồn tại thì trả về, nếu chưa thì tạo mới và trả về.
// Method này là thread-safe.
//
// Tham số:
//   - context: string - context để xác định nguồn gốc log
//
// Trả về:
//   - Logger: logger cho context đã cho (existing hoặc newly created)
//
// Ví dụ:
//
//	userLogger := manager.GetLogger("UserService")  // tạo mới nếu chưa có
//	userLogger2 := manager.GetLogger("UserService") // trả về cái đã tồn tại
func (m *manager) GetLogger(context string) Logger {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Kiểm tra nếu logger với context này đã tồn tại
	if logger, exists := m.loggers[context]; exists {
		return logger
	}

	// Tạo logger mới
	logger := NewLogger(context)

	// Thiết lập Level từ config
	logger.SetMinLevel(m.config.Level)

	// Bước 1: Luôn thêm Stack Handler nếu được enable
	if m.config.Stack.Enabled {
		if stackHandler := m.handlers[HandlerTypeStack]; stackHandler != nil {
			logger.AddHandler(HandlerTypeStack, stackHandler)
		}
	}

	// Bước 2: Chỉ thêm individual handlers khi cần thiết
	// Console: Chỉ thêm khi Stack không enable HOẶC Stack không có console
	if !m.config.Stack.Enabled || (m.config.Console.Enabled && !m.config.Stack.Handlers.Console) {
		if m.config.Console.Enabled {
			if consoleHandler := m.handlers[HandlerTypeConsole]; consoleHandler != nil {
				logger.AddHandler(HandlerTypeConsole, consoleHandler)
			}
		}
	}

	// File: Chỉ thêm khi Stack không enable HOẶC Stack không có file
	if !m.config.Stack.Enabled || (m.config.File.Enabled && !m.config.Stack.Handlers.File) {
		if m.config.File.Enabled {
			if fileHandler := m.handlers[HandlerTypeFile]; fileHandler != nil {
				logger.AddHandler(HandlerTypeFile, fileHandler)
			}
		}
	}

	// Lưu logger vào danh sách
	m.loggers[context] = logger

	return logger
}

// Close đóng tất cả các handlers đã đăng ký và giải phóng tài nguyên của chúng.
//
// Method này nên được gọi khi ứng dụng đang đóng để đảm bảo
// tất cả các file log được đóng đúng cách và tài nguyên được giải phóng.
//
// Trả về:
//   - error: lỗi đầu tiên gặp phải khi đóng handler, hoặc nil nếu tất cả đều đóng thành công
//
// Ví dụ:
//
//	if err := manager.Close(); err != nil {
//	    fmt.Fprintf(os.Stderr, "Lỗi khi đóng manager: %v\n", err)
//	}
func (m *manager) Close() error {
	// Tạo một bản sao của map handlers để giảm thiểu thời gian giữ lock
	m.mu.Lock()
	handlersCopy := make(map[HandlerType]handler.Handler, len(m.handlers))
	for k, v := range m.handlers {
		handlersCopy[k] = v
	}
	// Xóa tất cả handlers để tránh sử dụng sau khi đóng
	m.handlers = make(map[HandlerType]handler.Handler)
	m.mu.Unlock()

	// Đóng từng handler, theo dõi lỗi đầu tiên
	var firstErr error
	for handlerType, handler := range handlersCopy {
		// Bỏ qua handler nil
		if handler == nil {
			continue
		}
		if err := handler.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("failed to close handler %s: %w", handlerType, err)
		}
	}
	return firstErr
}

// initializeHandlers khởi tạo tất cả 3 handlers bắt buộc.
//
// Method này luôn tạo đầy đủ 3 handlers: console, file và stack theo config.
func (m *manager) initializeHandlers() {
	// Bắt buộc khởi tạo Console Handler
	consoleHandler := handler.NewConsoleHandler(m.config.Console.Colored)
	m.handlers[HandlerTypeConsole] = consoleHandler

	fileHandler, err := handler.NewFileHandler(m.config.File.Path, m.config.File.MaxSize)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file handler: %v", err))
	}

	m.handlers[HandlerTypeFile] = fileHandler
	// Khởi tạo Stack Handler với cấu hình
	stackHandler := handler.NewStackHandler()

	// Chỉ thêm handlers vào stack khi được cấu hình
	if m.config.Stack.Handlers.Console {
		stackHandler.AddHandler(consoleHandler)
	}

	if m.config.Stack.Handlers.File {
		stackHandler.AddHandler(fileHandler)
	}

	m.handlers[HandlerTypeStack] = stackHandler
}
