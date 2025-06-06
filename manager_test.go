package log

import (
	"errors"
	"testing"

	"go.fork.vn/log/handler"
)

// TestHandlerType định nghĩa handler type cho testing
const TestHandlerType HandlerType = "test"

// createTestConfig tạo config cho testing với các handlers minimal
func createTestConfig() *Config {
	return &Config{
		Level: handler.InfoLevel,
		Console: ConsoleConfig{
			Enabled: true,
			Colored: false,
		},
		File: FileConfig{
			Enabled: true,
			Path:    "/tmp/test_manager.log", // Provide a valid path
			MaxSize: 1024 * 1024,             // 1MB
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

// MockHandler triển khai interface handler.Handler để kiểm tra
type MockHandler struct {
	LogCalled   bool
	CloseCalled bool
	ShouldError bool
	LogLevel    handler.Level
	LogMessage  string
	LogArgs     []interface{}
}

func (m *MockHandler) Log(level handler.Level, message string, args ...interface{}) error {
	m.LogCalled = true
	m.LogLevel = level
	m.LogMessage = message
	m.LogArgs = args
	if m.ShouldError {
		return errors.New("mock log error")
	}
	return nil
}

func (m *MockHandler) Close() error {
	m.CloseCalled = true
	if m.ShouldError {
		return errors.New("mock close error")
	}
	return nil
}

func TestNewManager(t *testing.T) {
	config := createTestConfig()
	m := NewManager(config)

	if m == nil {
		t.Fatal("NewManager() trả về nil")
	}

	// Kiểm tra kiểu đúng
	_, ok := m.(*manager)
	if !ok {
		t.Errorf("NewManager() không trả về *manager, got %T", m)
	}

	// Kiểm tra các thuộc tính mặc định
	defaultManager := m.(*manager)
	if defaultManager.config.Level != handler.InfoLevel {
		t.Errorf("Manager mới không đặt config.Level mặc định là InfoLevel, got %v", defaultManager.config.Level)
	}
	if len(defaultManager.handlers) != 3 { // Console, File, Stack handlers luôn được khởi tạo
		t.Errorf("Manager mới không có đúng số handlers, got %d handlers", len(defaultManager.handlers))
	}
}

func TestManager_AddHandler(t *testing.T) {
	config := createTestConfig()
	m := NewManager(config)
	h := &MockHandler{}

	// Thêm handler mới
	m.AddHandler(TestHandlerType, h)

	// Kiểm tra handler đã được thêm
	retrievedHandler := m.GetHandler(TestHandlerType)
	if retrievedHandler != h {
		t.Errorf("AddHandler không thêm handler đúng, got %v, want %v", retrievedHandler, h)
	}
}

func TestManager_RemoveHandler(t *testing.T) {
	config := createTestConfig()
	m := NewManager(config)
	h := &MockHandler{}

	// Thêm handler trước
	m.AddHandler(TestHandlerType, h)

	// Xóa handler
	m.RemoveHandler(TestHandlerType)

	// Kiểm tra handler đã bị xóa
	retrievedHandler := m.GetHandler(TestHandlerType)
	if retrievedHandler != nil {
		t.Errorf("RemoveHandler không xóa handler, vẫn có %v", retrievedHandler)
	}

	// Kiểm tra Close() được gọi
	if !h.CloseCalled {
		t.Error("RemoveHandler không gọi Close() trên handler")
	}
}

func TestManager_GetHandler(t *testing.T) {
	config := createTestConfig()
	m := NewManager(config)

	// Test handler không tồn tại
	handler := m.GetHandler(TestHandlerType)
	if handler != nil {
		t.Errorf("GetHandler trả về non-nil cho handler không tồn tại, got %v", handler)
	}

	// Test handler tồn tại (console handler được tạo mặc định)
	consoleHandler := m.GetHandler(HandlerTypeConsole)
	if consoleHandler == nil {
		t.Error("GetHandler trả về nil cho console handler đã được khởi tạo")
	}
}

func TestManager_GetLogger(t *testing.T) {
	config := createTestConfig()
	m := NewManager(config)

	// Test tạo logger mới
	logger1 := m.GetLogger("TestService")
	if logger1 == nil {
		t.Fatal("GetLogger trả về nil")
	}

	// Test get-or-create pattern
	logger2 := m.GetLogger("TestService")
	if logger1 != logger2 {
		t.Error("GetLogger không trả về cùng instance cho cùng context")
	}

	// Test logger khác nhau cho context khác nhau
	logger3 := m.GetLogger("AnotherService")
	if logger1 == logger3 {
		t.Error("GetLogger trả về cùng instance cho context khác nhau")
	}
}

func TestManager_GetLogger_Logging(t *testing.T) {
	config := createTestConfig()
	config.Level = handler.DebugLevel // Đặt level thấp để test tất cả levels
	m := NewManager(config)

	// Thêm mock handler để test
	mockHandler := &MockHandler{}
	m.AddHandler(TestHandlerType, mockHandler)

	// Lấy logger và test logging
	logger := m.GetLogger("TestService")

	// Set handler cho logger sau khi logger đã được tạo
	m.SetHandler("TestService", TestHandlerType)

	// Test Debug log
	logger.Debug("debug message")
	if !mockHandler.LogCalled {
		t.Error("Debug log không được gọi")
	}
	if mockHandler.LogLevel != handler.DebugLevel {
		t.Errorf("Debug log level sai, got %v, want %v", mockHandler.LogLevel, handler.DebugLevel)
	}

	// Reset mock
	mockHandler.LogCalled = false

	// Test Info log
	logger.Info("info message")
	if !mockHandler.LogCalled {
		t.Error("Info log không được gọi")
	}
	if mockHandler.LogLevel != handler.InfoLevel {
		t.Errorf("Info log level sai, got %v, want %v", mockHandler.LogLevel, handler.InfoLevel)
	}
}

func TestManager_Close(t *testing.T) {
	config := createTestConfig()
	m := NewManager(config)

	h1 := &MockHandler{}
	h2 := &MockHandler{}

	m.AddHandler(TestHandlerType, h1)
	m.AddHandler("custom2", h2)

	// Close manager
	err := m.Close()
	if err != nil {
		t.Errorf("Close() trả về lỗi: %v", err)
	}

	// Kiểm tra tất cả handlers đã được close
	if !h1.CloseCalled {
		t.Error("Handler 1 không được close")
	}
	if !h2.CloseCalled {
		t.Error("Handler 2 không được close")
	}
}

func TestManager_CloseWithError(t *testing.T) {
	config := createTestConfig()
	m := NewManager(config)

	h := &MockHandler{ShouldError: true}
	m.AddHandler(TestHandlerType, h)

	// Close manager - nên trả về lỗi
	err := m.Close()
	if err == nil {
		t.Error("Close() không trả về lỗi khi handler.Close() lỗi")
	}
}

func TestManager_SetHandler(t *testing.T) {
	config := createTestConfig()
	m := NewManager(config)

	// Test SetHandler với logger context
	m.SetHandler("TestService", HandlerTypeConsole)

	// Lấy logger và kiểm tra nó có thể log
	logger := m.GetLogger("TestService")
	logger.Info("test message") // Không nên panic
}

func TestManager_ConcurrentAccess(t *testing.T) {
	config := createTestConfig()
	m := NewManager(config)

	done := make(chan bool, 10)

	// Concurrent GetLogger calls
	for i := 0; i < 10; i++ {
		go func(id int) {
			logger := m.GetLogger("Service" + string(rune('A'+id)))
			logger.Info("Concurrent message from %d", id)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Close không nên panic
	err := m.Close()
	if err != nil {
		t.Errorf("Concurrent close failed: %v", err)
	}
}

// Benchmarks

func BenchmarkManager_GetLogger(b *testing.B) {
	config := createTestConfig()
	m := NewManager(config)
	defer m.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.GetLogger("BenchService")
	}
}

func BenchmarkManager_GetLoggerDifferentContexts(b *testing.B) {
	config := createTestConfig()
	m := NewManager(config)
	defer m.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		context := "Service" + string(rune('A'+(i%26)))
		_ = m.GetLogger(context)
	}
}

func BenchmarkManager_LoggerWithMockHandler(b *testing.B) {
	config := createTestConfig()
	m := NewManager(config)
	defer m.Close()

	mockHandler := &MockHandler{}
	m.AddHandler(TestHandlerType, mockHandler)

	logger := m.GetLogger("BenchService")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message %d", i)
	}
}

func BenchmarkManager_ConcurrentGetLogger(b *testing.B) {
	config := createTestConfig()
	m := NewManager(config)
	defer m.Close()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			context := "Service" + string(rune('A'+(i%26)))
			logger := m.GetLogger(context)
			logger.Info("Concurrent benchmark message %d", i)
			i++
		}
	})
}
