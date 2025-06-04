package log

import (
	"errors"
	"fmt"
	"testing"

	"go.fork.vn/log/handler"
)

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

func TestManager_New(t *testing.T) {
	m := NewManager()
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
	if defaultManager.minLevel != handler.InfoLevel {
		t.Errorf("Manager mới không đặt minLevel mặc định là InfoLevel, got %v", defaultManager.minLevel)
	}
	if len(defaultManager.handlers) != 0 {
		t.Errorf("Manager mới không có handlers trống, got %d handlers", len(defaultManager.handlers))
	}
}

func TestManager_AddHandler(t *testing.T) {
	m := NewManager().(*manager)
	h := &MockHandler{}

	// Test thêm handler
	m.AddHandler("test", h)

	// Kiểm tra handler được thêm đúng cách
	if _, ok := m.handlers["test"]; !ok {
		t.Errorf("AddHandler không thêm handler vào map")
	}

	// Thêm handler thứ hai
	h2 := &MockHandler{}
	m.AddHandler("test2", h2)

	// Kiểm tra cả hai handler tồn tại
	if _, ok := m.handlers["test"]; !ok {
		t.Errorf("Handler đầu tiên không còn tồn tại sau khi thêm handler thứ hai")
	}
	if _, ok := m.handlers["test2"]; !ok {
		t.Errorf("Handler thứ hai không được thêm vào map")
	}

	// Ghi đè lên handler cũ - kiểm tra handler cũ được đóng
	h3 := &MockHandler{}
	m.AddHandler("test", h3)

	// Xác minh rằng h3 thay thế h trong map
	handler, ok := m.handlers["test"]
	if !ok {
		t.Error("Handler 'test' không tồn tại sau khi ghi đè")
	}
	if handler != h3 {
		t.Error("Handler không được ghi đè đúng cách")
	}
	// Kiểm tra handler cũ đã được đóng
	if !h.CloseCalled {
		t.Error("Handler cũ không được đóng khi bị ghi đè bởi AddHandler")
	}
}

func TestManager_RemoveHandler(t *testing.T) {
	m := NewManager()
	h := &MockHandler{}

	// Thêm handler
	m.AddHandler("test", h)

	// Xóa handler
	m.RemoveHandler("test")

	// Kiểm tra handler được gọi Close
	if !h.CloseCalled {
		t.Error("RemoveHandler không gọi Close() trên handler")
	}

	// Kiểm tra handler đã bị xóa khỏi map
	defaultManager := m.(*manager)
	if _, ok := defaultManager.handlers["test"]; ok {
		t.Error("RemoveHandler không xóa handler khỏi map")
	}

	// Xóa handler không tồn tại không gây lỗi
	m.RemoveHandler("nonexistent")

	// Thêm handler có lỗi khi close
	h2 := &MockHandler{ShouldError: true}
	m.AddHandler("test2", h2)

	// RemoveHandler vẫn nên xóa handler ngay cả khi Close() gây lỗi
	m.RemoveHandler("test2")

	// Kiểm tra handler đã bị xóa khỏi map dù Close() lỗi
	if _, ok := defaultManager.handlers["test2"]; ok {
		t.Error("RemoveHandler không xóa handler khỏi map khi Close() lỗi")
	}
}

func TestManager_SetMinLevel(t *testing.T) {
	m := NewManager().(*manager)

	// Kiểm tra mức mặc định
	if m.minLevel != handler.InfoLevel {
		t.Errorf("Mức mặc định không phải InfoLevel, got %v", m.minLevel)
	}

	// Đặt mức mới
	m.SetMinLevel(handler.WarningLevel)

	// Kiểm tra mức được đặt đúng
	if m.minLevel != handler.WarningLevel {
		t.Errorf("SetMinLevel không đặt minLevel đúng, got %v, want %v",
			m.minLevel, handler.WarningLevel)
	}

	// Đặt mức thấp nhất
	m.SetMinLevel(handler.DebugLevel)
	if m.minLevel != handler.DebugLevel {
		t.Errorf("SetMinLevel không đặt minLevel thành DebugLevel, got %v", m.minLevel)
	}

	// Đặt mức cao nhất
	m.SetMinLevel(handler.FatalLevel)
	if m.minLevel != handler.FatalLevel {
		t.Errorf("SetMinLevel không đặt minLevel thành FatalLevel, got %v", m.minLevel)
	}
}

func TestManager_LogMethods(t *testing.T) {
	m := NewManager()
	h := &MockHandler{}

	// Set min level to DebugLevel to ensure Debug logs are processed
	m.SetMinLevel(handler.DebugLevel)
	m.AddHandler("test", h)

	tests := []struct {
		name    string
		logFunc func(message string, args ...interface{})
		level   handler.Level
		message string
		args    []interface{}
	}{
		{"Debug", m.Debug, handler.DebugLevel, "debug message", []interface{}{1, 2}},
		{"Info", m.Info, handler.InfoLevel, "info message", []interface{}{3, 4}},
		{"Warning", m.Warning, handler.WarningLevel, "warn message", []interface{}{5, 6}},
		{"Error", m.Error, handler.ErrorLevel, "error message", []interface{}{7, 8}},
		{"Fatal", m.Fatal, handler.FatalLevel, "fatal message", []interface{}{9, 10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.LogCalled = false
			h.LogMessage = ""

			// Gọi method với message có định dạng
			formattedMsg := fmt.Sprintf("%s %%d %%d", tt.message)
			tt.logFunc(formattedMsg, tt.args...)

			// Kiểm tra handler được gọi với tham số đúng
			if !h.LogCalled {
				t.Errorf("%s không gọi handler", tt.name)
			}
			if h.LogLevel != tt.level {
				t.Errorf("%s truyền sai log level: got %v, want %v", tt.name, h.LogLevel, tt.level)
			}

			// Kiểm tra message được định dạng đúng
			expectedMsg := fmt.Sprintf(formattedMsg, tt.args...)
			if h.LogMessage != expectedMsg {
				t.Errorf("%s không định dạng message đúng: got %q, want %q",
					tt.name, h.LogMessage, expectedMsg)
			}
		})
	}
}

func TestManager_Close(t *testing.T) {
	m := NewManager()
	h1 := &MockHandler{}
	h2 := &MockHandler{}

	m.AddHandler("h1", h1)
	m.AddHandler("h2", h2)

	// Test Close method
	err := m.Close()
	if err != nil {
		t.Errorf("Close trả về lỗi: %v", err)
	}

	// Kiểm tra cả hai handlers được đóng
	if !h1.CloseCalled {
		t.Error("Close không gọi Close() trên handler đầu tiên")
	}
	if !h2.CloseCalled {
		t.Error("Close không gọi Close() trên handler thứ hai")
	}

	// Lưu ý: Theo hiện thực hiện tại, Close() không xóa các handlers khỏi map
	// Nó chỉ đóng các handlers nhưng vẫn giữ chúng trong map
	defaultManager := m.(*manager)
	if len(defaultManager.handlers) == 0 {
		t.Error("Không mong đợi map handlers trống sau khi close, hiện thực chỉ đóng handlers")
	}
}

func TestManager_Close_WithError(t *testing.T) {
	m := NewManager()
	h1 := &MockHandler{}
	h2 := &MockHandler{ShouldError: true}
	h3 := &MockHandler{}

	m.AddHandler("h1", h1)
	m.AddHandler("h2", h2)
	m.AddHandler("h3", h3)

	// Test Close method với handler trả về lỗi
	err := m.Close()
	if err == nil {
		t.Error("Close với handler lỗi không trả về lỗi")
	}

	// Kiểm tra tất cả handlers được đóng dù có lỗi
	if !h1.CloseCalled || !h2.CloseCalled || !h3.CloseCalled {
		t.Error("Close không gọi Close() trên tất cả handlers")
	}
}

func TestManager_LogFilteringByMinLevel(t *testing.T) {
	m := NewManager()
	h := &MockHandler{}

	m.AddHandler("test", h)

	// Đặt mức tối thiểu là ErrorLevel
	m.SetMinLevel(handler.ErrorLevel)

	// Log ở mức thấp hơn không nên gọi handler
	m.Debug("debug message")
	if h.LogCalled {
		t.Error("Debug log không bị lọc khi dưới ngưỡng")
	}

	m.Info("info message")
	if h.LogCalled {
		t.Error("Info log không bị lọc khi dưới ngưỡng")
	}

	m.Warning("warning message")
	if h.LogCalled {
		t.Error("Warning log không bị lọc khi dưới ngưỡng")
	}

	// Reset
	h.LogCalled = false

	// Log ở mức bằng hoặc cao hơn nên gọi handler
	m.Error("error message")
	if !h.LogCalled {
		t.Error("Error log bị lọc khi bằng hoặc trên ngưỡng")
	}

	// Reset
	h.LogCalled = false

	m.Fatal("fatal message")
	if !h.LogCalled {
		t.Error("Fatal log bị lọc khi bằng hoặc trên ngưỡng")
	}
}

// TestLogWithFormatting kiểm tra định dạng thông điệp
func TestManager_LogWithFormatting(t *testing.T) {
	m := NewManager()
	h := &MockHandler{}

	m.AddHandler("test", h)

	// Kiểm tra log với định dạng khác nhau
	tests := []struct {
		name   string
		format string
		args   []interface{}
		want   string
	}{
		{"String", "Hello %s", []interface{}{"world"}, "Hello world"},
		{"Number", "Number: %d", []interface{}{42}, "Number: 42"},
		{"Multiple", "%s: %d, %f", []interface{}{"Test", 123, 45.67}, "Test: 123, 45.670000"},
		{"Empty", "No args", nil, "No args"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.LogCalled = false
			h.LogMessage = ""

			// Gọi phương thức Info với định dạng
			m.Info(tt.format, tt.args...)

			// Kiểm tra định dạng đúng
			if !h.LogCalled {
				t.Error("Log không gọi handler")
			}
			if h.LogMessage != tt.want {
				t.Errorf("Định dạng không đúng: got %q, want %q", h.LogMessage, tt.want)
			}
		})
	}
}

func TestManager_GetHandler(t *testing.T) {
	m := NewManager()
	h1 := &MockHandler{}
	h2 := &MockHandler{}

	// Thêm các handlers
	m.AddHandler("handler1", h1)
	m.AddHandler("handler2", h2)

	// Lấy handler đã đăng ký
	handlerResult := m.GetHandler("handler1")
	if handlerResult == nil {
		t.Error("GetHandler trả về nil cho handler đã đăng ký")
	}

	if handlerResult != h1 {
		t.Errorf("GetHandler không trả về đúng handler, got %v, want %v", handlerResult, h1)
	}

	// Lấy handler thứ hai
	handlerResult = m.GetHandler("handler2")
	if handlerResult != h2 {
		t.Errorf("GetHandler không trả về đúng handler cho key thứ hai, got %v, want %v", handlerResult, h2)
	}

	// Lấy handler không tồn tại
	handlerResult = m.GetHandler("nonexistent")
	if handlerResult != nil {
		t.Errorf("GetHandler không trả về nil cho handler không tồn tại, got %v", handlerResult)
	}

	// Kiểm tra thread-safety (thông qua kiểm tra chức năng cơ bản)
	// Thêm handler mới trong khi đang thực hiện GetHandler
	go func() {
		m.AddHandler("handler3", &MockHandler{})
	}()

	// Xóa handler trong khi đang thực hiện GetHandler
	go func() {
		m.RemoveHandler("handler1")
	}()

	// Lấy handler2 một lần nữa sau các thao tác đồng thời
	handlerResult = m.GetHandler("handler2")
	if handlerResult != h2 {
		t.Errorf("GetHandler không hoạt động đúng sau các thao tác đồng thời, got %v, want %v", handlerResult, h2)
	}
}

func TestManager_LogWithErrorHandler(t *testing.T) {
	m := NewManager()
	h := &MockHandler{ShouldError: true}

	m.AddHandler("test", h)

	// Log với handler trả về lỗi
	// Không cần kiểm tra gì vì lỗi chỉ được in ra stderr
	// Nhưng log call vẫn hoàn thành không bị panic
	m.Info("this will error")

	// Kiểm tra handler vẫn được gọi dù trả về lỗi
	if !h.LogCalled {
		t.Error("Log không gọi handler khi biết handler sẽ lỗi")
	}
}

// TestManagerConcurrency kiểm tra tính năng đồng thời
func TestManager_Concurrency(t *testing.T) {
	m := NewManager()
	h := &MockHandler{}
	m.AddHandler("test", h)

	// Tạo nhiều goroutines để test thread safety
	done := make(chan bool, 3)

	// Goroutine 1: thêm/xóa handlers
	go func() {
		for i := 0; i < 10; i++ {
			handler := &MockHandler{}
			m.AddHandler(fmt.Sprintf("concurrent_%d", i), handler)
			m.RemoveHandler(fmt.Sprintf("concurrent_%d", i))
		}
		done <- true
	}()

	// Goroutine 2: thay đổi log level
	go func() {
		levels := []handler.Level{
			handler.DebugLevel,
			handler.InfoLevel,
			handler.WarningLevel,
			handler.ErrorLevel,
			handler.FatalLevel,
		}
		for i := 0; i < 10; i++ {
			m.SetMinLevel(levels[i%len(levels)])
		}
		done <- true
	}()

	// Goroutine 3: log messages
	go func() {
		for i := 0; i < 10; i++ {
			m.Info("Concurrent log message %d", i)
		}
		done <- true
	}()

	// Đợi tất cả goroutines hoàn thành
	for i := 0; i < 3; i++ {
		<-done
	}

	// Kiểm tra manager vẫn hoạt động bình thường
	m.Info("Final test message")
}

// TestManagerWithNilHandler kiểm tra xử lý khi handler là nil
func TestManager_WithNilHandler(t *testing.T) {
	m := NewManager()

	// Thêm nil handler - không nên panic
	m.AddHandler("nil", nil)

	// Log message - không nên panic
	m.Info("Test with nil handler")

	// Close - không nên panic
	err := m.Close()
	if err != nil {
		t.Errorf("Close với nil handler trả về lỗi: %v", err)
	}
}

// TestManagerEdgeCases kiểm tra các edge cases
func TestManager_EdgeCases(t *testing.T) {
	t.Run("Empty handler name", func(t *testing.T) {
		m := NewManager()
		h := &MockHandler{}

		// Thêm handler với tên rỗng
		m.AddHandler("", h)

		// Lấy handler với tên rỗng
		retrieved := m.GetHandler("")
		if retrieved != h {
			t.Error("Không thể lấy handler với tên rỗng")
		}

		// Xóa handler với tên rỗng
		m.RemoveHandler("")
		if m.GetHandler("") != nil {
			t.Error("Handler với tên rỗng không bị xóa")
		}
	})

	t.Run("Overwrite existing handler", func(t *testing.T) {
		m := NewManager()
		h1 := &MockHandler{}
		h2 := &MockHandler{}

		// Thêm handler đầu tiên
		m.AddHandler("test", h1)

		// Ghi đè với handler thứ hai
		m.AddHandler("test", h2)

		// Kiểm tra handler thứ hai được sử dụng
		retrieved := m.GetHandler("test")
		if retrieved != h2 {
			t.Error("Handler không bị ghi đè")
		}
		if retrieved == h1 {
			t.Error("Handler cũ vẫn còn")
		}
	})

	t.Run("Remove non-existent handler", func(t *testing.T) {
		m := NewManager()

		// Xóa handler không tồn tại - không nên panic
		m.RemoveHandler("nonexistent")

		// Kiểm tra vẫn có thể hoạt động bình thường
		m.Info("Test after removing non-existent handler")
	})

	t.Run("Log with special characters", func(t *testing.T) {
		m := NewManager()
		h := &MockHandler{}
		m.AddHandler("test", h)

		specialChars := "Special chars: 你好 🚀 \n\t\r\\\"'"

		m.Info("Message: %s", specialChars)

		if !h.LogCalled {
			t.Error("Log với ký tự đặc biệt không gọi handler")
		}

		expectedMsg := fmt.Sprintf("Message: %s", specialChars)
		if h.LogMessage != expectedMsg {
			t.Error("Message với ký tự đặc biệt không đúng")
		}
	})
}

// BenchmarkNewManager đo hiệu suất tạo Manager mới
func BenchmarkManager_New(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager := NewManager()
		_ = manager
	}
}

// BenchmarkManagerAddHandler đo hiệu suất thêm handler
func BenchmarkManager_AddHandler(b *testing.B) {
	manager := NewManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		handler := &MockHandler{}
		b.StartTimer()

		manager.AddHandler("test", handler)

		b.StopTimer()
		manager.RemoveHandler("test")
		b.StartTimer()
	}
}

// BenchmarkManagerAddHandlerMultiple đo hiệu suất thêm nhiều handlers
func BenchmarkManager_AddHandler_Multiple(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		manager := NewManager()
		b.StartTimer()

		for j := 0; j < 10; j++ {
			handler := &MockHandler{}
			manager.AddHandler(string(rune('a'+j)), handler)
		}

		b.StopTimer()
		manager.Close()
		b.StartTimer()
	}
}

// BenchmarkManagerRemoveHandler đo hiệu suất xóa handler
func BenchmarkManager_RemoveHandler(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		manager := NewManager()
		handler := &MockHandler{}
		manager.AddHandler("test", handler)
		b.StartTimer()

		manager.RemoveHandler("test")
	}
}

// BenchmarkManagerGetHandler đo hiệu suất lấy handler
func BenchmarkManager_GetHandler(b *testing.B) {
	manager := NewManager()
	handler := &MockHandler{}
	manager.AddHandler("test", handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.GetHandler("test")
	}
}

// BenchmarkManagerGetHandlerNotFound đo hiệu suất lấy handler không tồn tại
func BenchmarkManager_GetHandler_NotFound(b *testing.B) {
	manager := NewManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.GetHandler("nonexistent")
	}
}

// BenchmarkManagerSetMinLevel đo hiệu suất đặt mức log tối thiểu
func BenchmarkManager_SetMinLevel(b *testing.B) {
	manager := NewManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.SetMinLevel(handler.Level(i % 5))
	}
}

// BenchmarkManagerDebug đo hiệu suất log Debug
func BenchmarkManager_Debug(b *testing.B) {
	manager := NewManager()
	mockHandler := &MockHandler{}
	manager.AddHandler("test", mockHandler)
	manager.SetMinLevel(handler.DebugLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Debug("Debug message %d", i)
	}
}

// BenchmarkManagerInfo đo hiệu suất log Info
func BenchmarkManager_Info(b *testing.B) {
	manager := NewManager()
	handler := &MockHandler{}
	manager.AddHandler("test", handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Info("Info message %d", i)
	}
}

// BenchmarkManagerWarning đo hiệu suất log Warning
func BenchmarkManager_Warning(b *testing.B) {
	manager := NewManager()
	handler := &MockHandler{}
	manager.AddHandler("test", handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Warning("Warning message %d", i)
	}
}

// BenchmarkManagerError đo hiệu suất log Error
func BenchmarkManager_Error(b *testing.B) {
	manager := NewManager()
	handler := &MockHandler{}
	manager.AddHandler("test", handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Error("Error message %d", i)
	}
}

// BenchmarkManagerFatal đo hiệu suất log Fatal
func BenchmarkManager_Fatal(b *testing.B) {
	manager := NewManager()
	handler := &MockHandler{}
	manager.AddHandler("test", handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Fatal("Fatal message %d", i)
	}
}

// BenchmarkManagerLogFiltering đo hiệu suất lọc log theo level
func BenchmarkManager_LogFiltering(b *testing.B) {
	manager := NewManager()
	mockHandler := &MockHandler{}
	manager.AddHandler("test", mockHandler)
	manager.SetMinLevel(handler.ErrorLevel) // Chỉ cho phép Error và Fatal

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Các log này sẽ bị lọc
		manager.Debug("Debug message %d", i)
		manager.Info("Info message %d", i)
		manager.Warning("Warning message %d", i)
	}
}

// BenchmarkManagerLogWithMultipleHandlers đo hiệu suất log với nhiều handlers
func BenchmarkManager_LogWithMultipleHandlers(b *testing.B) {
	manager := NewManager()

	// Thêm 5 handlers
	for i := 0; i < 5; i++ {
		handler := &MockHandler{}
		manager.AddHandler(string(rune('a'+i)), handler)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Info("Message to multiple handlers %d", i)
	}
}

// BenchmarkManagerClose đo hiệu suất đóng manager
func BenchmarkManager_Close(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		manager := NewManager()
		for j := 0; j < 3; j++ {
			handler := &MockHandler{}
			manager.AddHandler(string(rune('a'+j)), handler)
		}
		b.StartTimer()

		_ = manager.Close()
	}
}

// BenchmarkManagerCloseWithError đo hiệu suất đóng manager với handler lỗi
func BenchmarkManager_Close_WithError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		manager := NewManager()
		handler1 := &MockHandler{}
		handler2 := &MockHandler{ShouldError: true}
		handler3 := &MockHandler{}
		manager.AddHandler("h1", handler1)
		manager.AddHandler("h2", handler2)
		manager.AddHandler("h3", handler3)
		b.StartTimer()

		_ = manager.Close()
	}
}

// BenchmarkManagerLogComplexMessage đo hiệu suất log với message phức tạp
func BenchmarkManager_LogComplexMessage(b *testing.B) {
	manager := NewManager()
	handler := &MockHandler{}
	manager.AddHandler("test", handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Info("Complex message with %s, %d, %f, %t, %v",
			"string", 123, 45.67, true, map[string]int{"key": i})
	}
}

// BenchmarkManagerConcurrentAddRemove đo hiệu suất concurrent add/remove handlers
func BenchmarkManager_ConcurrentAddRemove(b *testing.B) {
	manager := NewManager()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			handlerName := string(rune('a' + (i % 26)))
			handler := &MockHandler{}

			manager.AddHandler(handlerName, handler)
			manager.RemoveHandler(handlerName)
			i++
		}
	})
}

// BenchmarkManagerConcurrentLog đo hiệu suất concurrent logging
func BenchmarkManager_ConcurrentLog(b *testing.B) {
	manager := NewManager()
	handler := &MockHandler{}
	manager.AddHandler("test", handler)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			manager.Info("Concurrent log message %d", i)
			i++
		}
	})
}

// BenchmarkManagerConcurrentMixed đo hiệu suất mixed concurrent operations
func BenchmarkManager_ConcurrentMixed(b *testing.B) {
	manager := NewManager()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 4 {
			case 0:
				mockHandler := &MockHandler{}
				manager.AddHandler(string(rune('a'+(i%26))), mockHandler)
			case 1:
				manager.Info("Mixed operation log %d", i)
			case 2:
				_ = manager.GetHandler("test")
			case 3:
				manager.SetMinLevel(handler.Level(i % 5))
			}
			i++
		}
	})
}

// BenchmarkManagerLogWithDifferentLevels đo hiệu suất log với các level khác nhau
func BenchmarkManagerLogWithDifferentLevels(b *testing.B) {
	levels := []struct {
		name     string
		logFunc  func(Manager, string, ...interface{})
		minLevel handler.Level
	}{
		{"Debug", func(m Manager, msg string, args ...interface{}) { m.Debug(msg, args...) }, handler.DebugLevel},
		{"Info", func(m Manager, msg string, args ...interface{}) { m.Info(msg, args...) }, handler.InfoLevel},
		{"Warning", func(m Manager, msg string, args ...interface{}) { m.Warning(msg, args...) }, handler.WarningLevel},
		{"Error", func(m Manager, msg string, args ...interface{}) { m.Error(msg, args...) }, handler.ErrorLevel},
		{"Fatal", func(m Manager, msg string, args ...interface{}) { m.Fatal(msg, args...) }, handler.FatalLevel},
	}

	for _, level := range levels {
		b.Run("Level_"+level.name, func(b *testing.B) {
			manager := NewManager()
			handler := &MockHandler{}
			manager.AddHandler("test", handler)
			manager.SetMinLevel(level.minLevel)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				level.logFunc(manager, "Message for level %s: %d", level.name, i)
			}
		})
	}
}

// BenchmarkManagerMemoryUsage đo memory footprint của manager operations
func BenchmarkManagerMemoryUsage(b *testing.B) {
	b.ReportAllocs()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager := NewManager()

		// Thêm handlers
		for j := 0; j < 5; j++ {
			handler := &MockHandler{}
			manager.AddHandler(string(rune('a'+j)), handler)
		}

		// Log messages
		for j := 0; j < 10; j++ {
			manager.Info("Memory test message %d-%d", i, j)
		}

		// Clean up
		manager.Close()
	}
}

// BenchmarkManagerConcurrent kiểm tra hiệu năng đồng thời
func BenchmarkManagerConcurrent(b *testing.B) {
	manager := NewManager()
	handler := &MockHandler{}
	manager.AddHandler("test", handler)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			manager.Info("Concurrent benchmark message")
		}
	})
}

// BenchmarkManagerMultipleHandlers kiểm tra hiệu năng với nhiều handlers
func BenchmarkManagerMultipleHandlers(b *testing.B) {
	manager := NewManager()

	// Thêm nhiều handlers
	for i := 0; i < 10; i++ {
		handler := &MockHandler{}
		manager.AddHandler(fmt.Sprintf("handler_%d", i), handler)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Info("Benchmark message %d", i)
	}
}
