package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileHandler triển khai một log handler ghi vào file với khả năng xoay vòng.
//
// Tính năng:
//   - Kiểm tra thư mục tồn tại và quyền ghi trước khi khởi tạo
//   - Xoay vòng log dựa trên kích thước
//   - Đặt tên file xoay vòng dựa trên timestamp
//   - Hoạt động thread-safe
//   - Định dạng timestamp chuẩn
//
// Yêu cầu:
//   - Thư mục chứa file log phải tồn tại trước
//   - Thư mục phải có quyền ghi
type FileHandler struct {
	path        string     // Đường dẫn đến file log
	file        *os.File   // File handle hiện tại
	maxSize     int64      // Kích thước file tối đa tính bằng byte trước khi xoay vòng
	currentSize int64      // Kích thước file hiện tại tính bằng byte
	mu          sync.Mutex // Mutex để đảm bảo thread-safety
}

// NewFileHandler tạo một file handler mới cho đường dẫn và kích thước tối đa được chỉ định.
//
// Tham số:
//   - path: string - đường dẫn đến file log
//   - maxSize: int64 - kích thước file tối đa tính bằng byte trước khi xoay vòng (0 để không giới hạn)
//
// Trả về:
//   - *FileHandler: một file handler đã được cấu hình
//   - error: nếu thư mục không tồn tại, không có quyền ghi, hoặc file không thể được mở
//
// Lưu ý:
//   - Thư mục chứa file log phải tồn tại trước khi gọi hàm này
//   - Thư mục phải có quyền ghi
//   - Hàm không tự động tạo thư mục
//
// Ví dụ:
//
//	// Tạo một handler với kích thước file tối đa 10MB
//	handler, err := handler.NewFileHandler("/var/log/app.log", 10*1024*1024)
//	if err != nil {
//	    fmt.Printf("Không thể tạo file log: %v\n", err)
//	}
func NewFileHandler(path string, maxSize int64) (*FileHandler, error) {
	var file *os.File
	var currentSize int64

	// Kiểm tra xem file path có tồn tại không
	if info, err := os.Stat(path); err == nil {
		// 1. Path tồn tại - thử mở file với quyền ghi thêm để kiểm tra có ghi được không
		file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("cannot open existing file for writing: %w", err)
		}
		currentSize = info.Size()
	} else if os.IsNotExist(err) {
		// 2. Path không tồn tại - kiểm tra parent directory
		dir := filepath.Dir(path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return nil, fmt.Errorf("path to folder do not exists: %s", dir)
		} else if err != nil {
			return nil, fmt.Errorf("cannot access parent directory: %w", err)
		}

		// 3. Parent directory tồn tại - thử tạo file với quyền ghi
		file, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			// Kiểm tra nếu lỗi là do permission denied
			if os.IsPermission(err) {
				return nil, fmt.Errorf("directory does not have write permission: %s", dir)
			}
			return nil, fmt.Errorf("cannot create file: %w", err)
		}
		currentSize = 0
	} else {
		// Lỗi khác khi kiểm tra file (có thể là permission denied)
		if os.IsPermission(err) {
			dir := filepath.Dir(path)
			return nil, fmt.Errorf("directory does not have write permission: %s", dir)
		}
		return nil, fmt.Errorf("cannot access file path: %w", err)
	}

	// Khởi tạo handler
	handler := &FileHandler{
		path:        path,
		file:        file,
		maxSize:     maxSize,
		currentSize: currentSize,
	}

	return handler, nil
}

// Log ghi một log entry vào file.
//
// Method này định dạng log entry với timestamp và chỉ báo cấp độ
// và ghi vào file, xoay vòng file nếu nó vượt quá kích thước tối đa.
//
// Tham số:
//   - level: Level - cấp độ nghiêm trọng của log entry
//   - message: string - thông điệp log
//   - args: ...interface{} - tham số định dạng tùy chọn
//
// Trả về:
//   - error: một lỗi nếu ghi vào file thất bại
func (a *FileHandler) Log(level Level, message string, args ...interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Kiểm tra xem file có cần xoay vòng không
	if a.maxSize > 0 && a.currentSize >= a.maxSize {
		if err := a.rotate(); err != nil {
			return fmt.Errorf("không thể xoay vòng file log: %w", err)
		}
	}

	// Định dạng với timestamp và mức độ
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	// Định dạng thông điệp nếu có tham số
	formattedMessage := message
	if len(args) > 0 {
		formattedMessage = fmt.Sprintf(message, args...)
	}
	formattedMessage = fmt.Sprintf("%s [%s] %s\n", timestamp, level.String(), formattedMessage)

	// Ghi vào file
	n, err := a.file.WriteString(formattedMessage)
	if err != nil {
		return fmt.Errorf("không thể ghi vào file log: %w", err)
	}

	// Cập nhật kích thước file hiện tại
	a.currentSize += int64(n)

	return nil
}

// Close đóng file log một cách chính xác.
//
// Phương thức này nên được gọi khi handler không còn cần thiết nữa
// để đảm bảo file được đóng chính xác và tất cả dữ liệu được ghi đệm.
//
// Trả về:
//   - error: một lỗi nếu đóng file thất bại
func (a *FileHandler) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.file != nil {
		if err := a.file.Close(); err != nil {
			return fmt.Errorf("không thể đóng file log: %w", err)
		}
		a.file = nil
	}

	return nil
}

// rotate thực hiện xoay vòng file log khi kích thước file vượt quá giới hạn tối đa.
//
// File hiện tại được đổi tên với hậu tố timestamp, và một file mới được tạo.
//
// Trả về:
//   - error: một lỗi nếu việc xoay vòng thất bại
func (a *FileHandler) rotate() error {
	// Đóng file hiện tại
	if err := a.file.Close(); err != nil {
		return fmt.Errorf("không thể đóng file log hiện tại: %w", err)
	}

	// Tạo tên file sao lưu với timestamp
	backupPath := fmt.Sprintf("%s.%s", a.path, time.Now().Format("20060102150405"))

	// Đổi tên file hiện tại thành file sao lưu
	if err := os.Rename(a.path, backupPath); err != nil {
		return fmt.Errorf("không thể đổi tên file log: %w", err)
	}

	// Mở file log mới
	file, err := os.OpenFile(a.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("không thể mở file log mới: %w", err)
	}

	// Cập nhật trạng thái handler
	a.file = file
	a.currentSize = 0

	return nil
}
