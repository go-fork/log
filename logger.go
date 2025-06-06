package log

import (
	"fmt"
	"sync"

	"go.fork.vn/log/handler"
)

// Logger định nghĩa interface cho hệ thống logging tập trung.
//
// Interface Logger cung cấp các method để ghi log ở nhiều cấp độ nghiêm trọng
// khác nhau và để quản lý các handler.
//
// Các triển khai của interface này cần đảm bảo thread-safe và xử lý
// việc phân phối các log entry đến tất cả các handler đã đăng ký.
type Logger interface {
	// Debug ghi một thông điệp ở cấp độ debug.
	//
	// Tham số:
	//   - message: string - thông điệp log (có thể là chuỗi định dạng)
	//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
	Debug(message string, args ...interface{})

	// Info ghi một thông điệp ở cấp độ info.
	//
	// Tham số:
	//   - message: string - thông điệp log (có thể là chuỗi định dạng)
	//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
	Info(message string, args ...interface{})

	// Warning ghi một thông điệp ở cấp độ warning.
	//
	// Tham số:
	//   - message: string - thông điệp log (có thể là chuỗi định dạng)
	//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
	Warning(message string, args ...interface{})

	// Error ghi một thông điệp ở cấp độ error.
	//
	// Tham số:
	//   - message: string - thông điệp log (có thể là chuỗi định dạng)
	//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
	Error(message string, args ...interface{})

	// Fatal ghi một thông điệp ở cấp độ fatal.
	//
	// Tham số:
	//   - message: string - thông điệp log (có thể là chuỗi định dạng)
	//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
	Fatal(message string, args ...interface{})

	// AddHandler đăng ký một handler mới vào logger.
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

	// SetMinLevel thiết lập ngưỡng cấp độ log tối thiểu.
	//
	// Tham số:
	//   - level: handler.Level - cấp độ tối thiểu để log
	SetMinLevel(level handler.Level)

	// Close đóng logger và tất cả các handler.
	//
	// Trả về:
	//   - error: một lỗi nếu việc đóng handler thất bại
	Close() error
}

// logger là triển khai chuẩn của interface Logger.
//
// logger cung cấp cách quản lý nhiều handler log với bộ lọc
// dựa trên cấp độ thread-safe. Nó được thiết kế cho truy cập đồng thời
// trong môi trường đa goroutine.
//
// Tính năng:
//   - Quản lý handler thread-safe bằng RWMutex
//   - Lọc cấp độ log
//   - Thêm/xóa handler động
//   - Dọn dẹp tài nguyên an toàn khi tắt
//   - Context cố định để xác định nguồn gốc log (immutable sau khi tạo)
type logger struct {
	handlers map[HandlerType]handler.Handler // Map các handler theo loại
	minLevel handler.Level                   // Ngưỡng cấp độ log tối thiểu
	context  string                          // Context cố định để xác định nguồn gốc log (immutable)
	mu       sync.RWMutex                    // Mutex để đảm bảo thread-safety
}

// NewLogger tạo và trả về một instance logger mới với context cố định.
//
// Hàm này khởi tạo một logger không có handler nào và InfoLevel là
// cấp độ log tối thiểu mặc định. Context được thiết lập khi tạo và không thể thay đổi.
//
// Tham số:
//   - context: string - context cố định để xác định nguồn gốc log (VD: UserService, UserController)
//
// Trả về:
//   - Logger: một instance mới của logger triển khai interface Logger.
//
// Ví dụ:
//
//	logger := log.NewLogger("UserService")
//	logger.AddHandler("console", handler.NewConsoleHandler(true))
//	// context "UserService" sẽ không thể thay đổi trong suốt vòng đời của logger
func NewLogger(context string) Logger {
	return &logger{
		handlers: make(map[HandlerType]handler.Handler),
		minLevel: handler.InfoLevel, // Mặc định là InfoLevel
		context:  context,           // Thiết lập context từ tham số
	}
}

// Debug ghi một thông điệp ở cấp độ debug.
//
// Debug logs dành cho thông tin chẩn đoán chi tiết hữu ích trong quá trình
// phát triển hoặc khắc phục sự cố.
//
// Tham số:
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
//
// Ví dụ:
//
//	logger.Debug("Lần thử kết nối %d đến %s", attempt, serverAddress)
func (l *logger) Debug(message string, args ...interface{}) {
	// Ghi log ở cấp độ DebugLevel
	l.log(handler.DebugLevel, message, args...)
}

// Info ghi một thông điệp ở cấp độ info.
//
// Info logs dành cho thông tin hoạt động chung về hành vi
// bình thường của ứng dụng.
//
// Tham số:
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
//
// Ví dụ:
//
//	logger.Info("Máy chủ đã khởi động trên cổng %d", port)
func (l *logger) Info(message string, args ...interface{}) {
	// Ghi log ở cấp độ InfoLevel
	l.log(handler.InfoLevel, message, args...)
}

// Warning ghi một thông điệp ở cấp độ warning.
//
// Warning logs chỉ ra các vấn đề tiềm ẩn hoặc điều kiện không mong đợi
// mà không phải lỗi nhưng có thể cần chú ý.
//
// Tham số:
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
//
// Ví dụ:
//
//	logger.Warning("Sử dụng bộ nhớ cao: %d MB", memoryUsage)
func (l *logger) Warning(message string, args ...interface{}) {
	// Ghi log ở cấp độ WarningLevel
	l.log(handler.WarningLevel, message, args...)
}

// Error ghi một thông điệp ở cấp độ error.
//
// Error logs chỉ ra các lỗi hoặc thất bại ảnh hưởng đến hoạt động bình thường
// nhưng không yêu cầu phải kết thúc ngay lập tức.
//
// Tham số:
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
//
// Ví dụ:
//
//	logger.Error("Xử lý yêu cầu thất bại: %v", err)
func (l *logger) Error(message string, args ...interface{}) {
	// Ghi log ở cấp độ ErrorLevel
	l.log(handler.ErrorLevel, message, args...)
}

// Fatal ghi một thông điệp ở cấp độ fatal.
//
// Fatal logs chỉ ra các lỗi nghiêm trọng thường yêu cầu kết thúc ứng dụng
// hoặc cần sự chú ý ngay lập tức của người quản trị.
//
// Tham số:
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
//
// Ví dụ:
//
//	logger.Fatal("Kết nối database thất bại: %v", err)
func (l *logger) Fatal(message string, args ...interface{}) {
	// Ghi log ở cấp độ FatalLevel
	l.log(handler.FatalLevel, message, args...)
}

// AddHandler thêm một handler log mới vào logger.
//
// Method này đăng ký một handler với loại đã cho. Nếu một handler với cùng loại
// đã tồn tại, nó sẽ bị thay thế mà không đóng handler cũ. Method này là thread-safe.
//
// Tham số:
//   - handlerType: HandlerType - loại handler (console, file, stack)
//   - handler: handler.Handler - triển khai handler cần thêm
//
// Ví dụ:
//
//	// Thêm một file handler
//	fileHandler, _ := handler.NewFileHandler("app.log", 10*1024*1024)
//	logger.AddHandler(HandlerTypeFile, fileHandler)
func (l *logger) AddHandler(handlerType HandlerType, handler handler.Handler) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// Nếu handler cũ cùng loại tồn tại, đóng lại để tránh leak resource
	if old, ok := l.handlers[handlerType]; ok {
		old.Close()
	}
	l.handlers[handlerType] = handler
}

// RemoveHandler xóa một handler khỏi logger theo loại.
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
//	logger.RemoveHandler(HandlerTypeFile) // Xóa và đóng file handler
func (l *logger) RemoveHandler(handlerType HandlerType) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// Đóng và xóa handler nếu nó tồn tại
	if handler, ok := l.handlers[handlerType]; ok {
		handler.Close()
		delete(l.handlers, handlerType)
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
//	if h := logger.GetHandler(HandlerTypeFile); h != nil {
//	    // Sử dụng handler
//	}
func (l *logger) GetHandler(handlerType HandlerType) handler.Handler {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Trả về handler nếu tồn tại hoặc nil nếu không tìm thấy
	return l.handlers[handlerType]
}

// SetMinLevel thiết lập cấp độ log tối thiểu cho logger.
//
// Bất kỳ log entry nào có cấp độ dưới ngưỡng này sẽ bị bỏ qua.
// Method này là thread-safe.
//
// Tham số:
//   - level: handler.Level - cấp độ log tối thiểu cần thiết lập
//
// Ví dụ:
//
//	// Chỉ xử lý log Warning, Error và Fatal
//	logger.SetMinLevel(handler.WarningLevel)
func (l *logger) SetMinLevel(level handler.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.minLevel = level
}

// Close đóng tất cả các handler log đã đăng ký và giải phóng tài nguyên của chúng.
//
// Method này nên được gọi khi ứng dụng đang đóng để đảm bảo
// tất cả các file log được đóng đúng cách và tài nguyên được giải phóng.
//
// Trả về:
//   - error: lỗi đầu tiên gặp phải khi đóng handler, hoặc nil nếu tất cả đều đóng thành công
//
// Ví dụ:
//
//	if err := logger.Close(); err != nil {
//	    fmt.Fprintf(os.Stderr, "Lỗi khi đóng log logger: %v\n", err)
//	}
func (l *logger) Close() error {
	// Tạo một bản sao của map handlers để giảm thiểu thời gian giữ lock
	l.mu.Lock()
	handlersCopy := make(map[HandlerType]handler.Handler, len(l.handlers))
	for k, v := range l.handlers {
		handlersCopy[k] = v
	}
	l.mu.Unlock()

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

// log là method nội bộ để ghi một log entry đến tất cả các handler.
//
// Method này xử lý lọc cấp độ, định dạng thông điệp với context và gửi
// log entry đến tất cả các handler đã đăng ký. Nó được thiết kế để giảm thiểu
// thời gian giữ lock để tăng concurrency.
//
// Tham số:
//   - level: handler.Level - cấp độ log của thông điệp
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - tham số tùy chọn để định dạng thông điệp
func (l *logger) log(level handler.Level, message string, args ...interface{}) {
	// Bỏ qua nếu dưới cấp độ tối thiểu
	if level < l.minLevel {
		return
	}

	// Lấy snapshot của handlers để giảm thiểu thời gian giữ lock
	l.mu.RLock()
	handlersCopy := make(map[HandlerType]handler.Handler, len(l.handlers))
	for k, v := range l.handlers {
		handlersCopy[k] = v
	}
	l.mu.RUnlock()

	// Định dạng thông điệp nếu có tham số
	formattedMessage := message
	if len(args) > 0 {
		formattedMessage = fmt.Sprintf(message, args...)
	}

	// Thêm context vào thông điệp nếu có (context là immutable nên không cần lock)
	if l.context != "" {
		formattedMessage = fmt.Sprintf("[%s] %s", l.context, formattedMessage)
	}

	// Ghi log entry đến tất cả các handler
	for handlerType, handler := range handlersCopy {
		// Bỏ qua handler nil
		if handler == nil {
			continue
		}
		if err := handler.Log(level, formattedMessage); err != nil {
			// Xử lý lỗi logging (ghi ra stderr)
			fmt.Printf("Lỗi khi ghi log đến handler %s: %v\n", handlerType, err)
		}
	}
}
