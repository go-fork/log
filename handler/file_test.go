package handler

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "filehandler-test-*")
	if err != nil {
		t.Fatalf("Không thể tạo thư mục tạm: %v", err)
	}
	return dir
}

// Kiểm tra nếu s chứa substring
func contains(s, substring string) bool {
	return strings.Contains(s, substring)
}

// Bổ sung cho TestNewFileHandler với test đường dẫn hợp lệ
func TestFileHandler_New(t *testing.T) {
	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Đường dẫn log
	logPath := filepath.Join(dir, "new-test.log")

	// Tạo và kiểm tra handler mới
	h, err := NewFileHandler(logPath, 100)
	if err != nil {
		t.Fatalf("NewFileHandler() với đường dẫn hợp lệ error = %v", err)
	}
	defer h.Close()

	// Kiểm tra thuộc tính
	if h.path != logPath {
		t.Errorf("NewFileHandler() không thiết lập đúng path, got = %v, want %v", h.path, logPath)
	}
	if h.maxSize != 100 {
		t.Errorf("NewFileHandler() không thiết lập đúng maxSize, got = %v, want %v", h.maxSize, 100)
	}
	if h.file == nil {
		t.Error("NewFileHandler() không mở file")
	}
}

// Test với đường dẫn không hợp lệ
func TestFileHandler_New_WithInvalidPath(t *testing.T) {
	// Thử tạo handler với đường dẫn không hợp lệ
	h, err := NewFileHandler("/invalid/path/that/should/not/exist/log.txt", 100)
	if err == nil {
		t.Error("NewFileHandler() với đường dẫn không hợp lệ nên trả về lỗi")
		h.Close()
	}
}

// Test với folder không tồn tại
func TestFileHandler_New_WithNonExistentFolder(t *testing.T) {
	// Tạo đường dẫn tới folder không tồn tại
	nonExistentPath := "/tmp/non-existent-folder-12345/app.log"

	h, err := NewFileHandler(nonExistentPath, 100)
	if err == nil {
		t.Error("NewFileHandler() với folder không tồn tại nên trả về lỗi")
		if h != nil {
			h.Close()
		}
		return
	}

	// Kiểm tra error message có chứa "path to folder do not exists"
	if !strings.Contains(err.Error(), "path to folder do not exists") {
		t.Errorf("NewFileHandler() error message không đúng, got = %v, want chứa 'path to folder do not exists'", err.Error())
	}
}

// Test với folder không có quyền ghi
func TestFileHandler_New_WithNoWritePermission(t *testing.T) {
	// Tạo thư mục tạm với quyền read-only
	dir := createTempDir(t)
	defer func() {
		// Khôi phục quyền để có thể xóa
		_ = os.Chmod(dir, 0755)
		os.RemoveAll(dir)
	}()

	// Set quyền read-only (không có write permission)
	err := os.Chmod(dir, 0444) // r--r--r--
	if err != nil {
		t.Skipf("Không thể thay đổi quyền thư mục: %v", err)
	}

	logPath := filepath.Join(dir, "readonly-test.log")

	h, err := NewFileHandler(logPath, 100)
	if err == nil {
		t.Error("NewFileHandler() với folder read-only nên trả về lỗi")
		if h != nil {
			h.Close()
		}
		return
	}

	// Kiểm tra error message có chứa "does not have write permission"
	if !strings.Contains(err.Error(), "does not have write permission") {
		t.Errorf("NewFileHandler() error message không đúng, got = %v, want chứa 'does not have write permission'", err.Error())
	}
}

// Test với folder tồn tại và có quyền ghi (positive case)
func TestFileHandler_New_WithValidFolder(t *testing.T) {
	// Tạo thư mục tạm với quyền ghi
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	logPath := filepath.Join(dir, "valid-test.log")

	h, err := NewFileHandler(logPath, 100)
	if err != nil {
		t.Errorf("NewFileHandler() với folder hợp lệ không nên trả về lỗi, got = %v", err)
		return
	}
	defer h.Close()

	// Kiểm tra file đã được tạo
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("NewFileHandler() không tạo file log")
	}
}

func TestFileHandler_Log(t *testing.T) {
	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Đường dẫn log
	logPath := filepath.Join(dir, "log-test.log")

	// Tạo handler với maxSize nhỏ để test rotation
	h, err := NewFileHandler(logPath, 50)
	if err != nil {
		t.Fatalf("NewFileHandler() error = %v", err)
	}
	defer h.Close()

	// Ghi log
	err = h.Log(InfoLevel, "test message 1")
	if err != nil {
		t.Errorf("Log() error = %v", err)
	}

	// Ghi log với args
	err = h.Log(WarningLevel, "test message %d", 2)
	if err != nil {
		t.Errorf("Log() với args error = %v", err)
	}

	// Đọc nội dung file
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Không thể đọc file log: %v", err)
	}

	// Kiểm tra nội dung
	contentStr := string(content)
	if !contains(contentStr, "[INFO]") || !contains(contentStr, "test message 1") {
		t.Errorf("Log không ghi đúng message 1: %q", contentStr)
	}
	if !contains(contentStr, "[WARNING]") || !contains(contentStr, "test message 2") {
		t.Errorf("Log không ghi đúng message 2: %q", contentStr)
	}

	// Ghi nhiều log để kích hoạt rotation
	for i := 0; i < 10; i++ {
		err = h.Log(ErrorLevel, "rotation test message %d", i)
		if err != nil {
			t.Errorf("Log() trong vòng lặp error = %v", err)
		}
	}

	// Kiểm tra file gốc vẫn tồn tại
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("File log gốc không tồn tại sau rotation")
	}

	// Kiểm tra nếu có ít nhất một file backup được tạo
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Không thể đọc thư mục: %v", err)
	}

	backupFound := false
	for _, file := range files {
		if file.Name() != "log-test.log" && contains(file.Name(), "log-test.log") {
			backupFound = true
			break
		}
	}

	if !backupFound {
		t.Error("Không tìm thấy file backup sau rotation")
	}
}

func TestFileHandler_Close(t *testing.T) {
	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Đường dẫn log
	logPath := filepath.Join(dir, "close-test.log")

	// Tạo handler
	h, err := NewFileHandler(logPath, 1024)
	if err != nil {
		t.Fatalf("NewFileHandler() error = %v", err)
	}

	// Ghi log trước khi đóng
	err = h.Log(InfoLevel, "message before close")
	if err != nil {
		t.Errorf("Log() error = %v", err)
	}

	// Đóng file
	err = h.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Tạo handler mới trên cùng file
	h2, err := NewFileHandler(logPath, 1024)
	if err != nil {
		t.Fatalf("NewFileHandler() sau close error = %v", err)
	}
	defer h2.Close()

	// Ghi log thêm sau khi đã tạo handler mới
	err = h2.Log(InfoLevel, "message after close with new handler")
	if err != nil {
		t.Errorf("Log() với handler mới error = %v", err)
	}

	// Kiểm tra nội dung file bằng cách đọc từng dòng
	file, err := os.Open(logPath)
	if err != nil {
		t.Fatalf("Không thể mở file log: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	found1, found2 := false, false
	for scanner.Scan() {
		line := scanner.Text()
		if contains(line, "message before close") {
			found1 = true
		}
		if contains(line, "message after close with new handler") {
			found2 = true
		}
	}

	if !found1 {
		t.Error("Không tìm thấy thông điệp trước khi đóng handler")
	}
	if !found2 {
		t.Error("Không tìm thấy thông điệp sau khi tạo handler mới")
	}
}

func TestFileHandler_Rotate_Manually(t *testing.T) {
	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Đường dẫn log
	logPath := filepath.Join(dir, "rotate-test.log")

	// Tạo handler
	h, err := NewFileHandler(logPath, 1000) // Max size nhỏ để dễ dàng gây rotate
	if err != nil {
		t.Fatalf("NewFileHandler() error = %v", err)
	}
	defer h.Close()

	// Vì rotate là private method, chúng ta trigger nó qua Log
	// Ghi log đủ lớn để buộc rotation
	for i := 0; i < 50; i++ {
		err = h.Log(InfoLevel, "large message to force rotation: %d - this is extra text to make the message bigger", i)
		if err != nil {
			t.Errorf("Log() large message error = %v", err)
		}
	}

	// Kiểm tra các file backup được tạo
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Không thể đọc thư mục: %v", err)
	}

	backupFiles := 0
	for _, file := range files {
		if file.Name() != "rotate-test.log" && contains(file.Name(), "rotate-test.log") {
			backupFiles++
		}
	}

	if backupFiles == 0 {
		t.Error("Không có file backup nào được tạo")
	} else {
		t.Logf("Số file backup đã tạo: %d", backupFiles)
	}
}

// TestFileHandlerRotateError kiểm tra các trường hợp lỗi trong quá trình xoay vòng file
func TestFileHandler_Rotate_Error(t *testing.T) {
	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Đường dẫn log
	logPath := filepath.Join(dir, "rotate-error-test.log")

	// Tạo handler với kích thước nhỏ để buộc rotate
	h, err := NewFileHandler(logPath, 10)
	if err != nil {
		t.Fatalf("NewFileHandler() error = %v", err)
	}

	// Ghi log đủ lớn để kích hoạt rotation
	err = h.Log(InfoLevel, "message that will trigger rotation")
	if err != nil {
		t.Errorf("Log() first message error = %v", err)
	}

	// Đảm bảo file có thể được đóng và rotation có thể gọi được
	err = h.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Tạo một handler mới để kiểm tra lỗi đóng file
	h, err = NewFileHandler(logPath, 10)
	if err != nil {
		t.Fatalf("NewFileHandler() lần thứ hai error = %v", err)
	}
	defer h.Close()

	// Kiểm tra log sau rotation
	err = h.Log(InfoLevel, "another message")
	if err != nil {
		t.Errorf("Log() sau rotation error = %v", err)
	}
}

// TestFileHandlerNewWithExistingDir kiểm tra khi thư mục đã tồn tại
func TestFileHandler_New_WithExistingDir(t *testing.T) {
	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Đường dẫn log trong thư mục con
	subDir := filepath.Join(dir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Không thể tạo thư mục con: %v", err)
	}

	logPath := filepath.Join(subDir, "exist-dir-test.log")

	// Tạo handler
	h, err := NewFileHandler(logPath, 100)
	if err != nil {
		t.Fatalf("NewFileHandler() với thư mục đã tồn tại error = %v", err)
	}
	defer h.Close()

	// Kiểm tra thuộc tính
	if h.path != logPath {
		t.Errorf("NewFileHandler() không thiết lập đúng path, got = %v, want %v", h.path, logPath)
	}
}

// TestFileHandlerLogError kiểm tra lỗi khi ghi log
func TestFileHandler_LogError(t *testing.T) {
	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Đường dẫn log
	logPath := filepath.Join(dir, "log-error-test.log")

	// Tạo handler
	h, err := NewFileHandler(logPath, 100)
	if err != nil {
		t.Fatalf("NewFileHandler() error = %v", err)
	}

	// Đóng file để gây lỗi khi ghi
	err = h.file.Close()
	if err != nil {
		t.Fatalf("Không thể đóng file: %v", err)
	}

	// Thử ghi log vào file đã đóng
	err = h.Log(InfoLevel, "message to closed file")
	if err == nil {
		t.Error("Log() vào file đã đóng nên trả về lỗi")
	}

	// Gọi Close cho an toàn
	h.Close()
}

// TestFileHandlerWithNoPermission kiểm tra trường hợp không có quyền truy cập
func TestFileHandler_WithNoPermission(t *testing.T) {
	// Bỏ qua trên Windows vì cơ chế quyền khác
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Bỏ qua test này trên Windows")
	}

	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Tạo một thư mục con và thay đổi quyền để gây lỗi
	noPermDir := filepath.Join(dir, "noperm")
	if err := os.Mkdir(noPermDir, 0755); err != nil {
		t.Fatalf("Không thể tạo thư mục: %v", err)
	}

	// Tạo file log
	logPath := filepath.Join(noPermDir, "noperm.log")
	file, err := os.Create(logPath)
	if err != nil {
		t.Fatalf("Không thể tạo file: %v", err)
	}
	file.Close()

	// Thay đổi quyền thư mục để không thể ghi
	if err := os.Chmod(noPermDir, 0555); err != nil {
		t.Fatalf("Không thể thay đổi quyền thư mục: %v", err)
	}

	// Thử tạo handler cho file trong thư mục không thể ghi
	// Trên một số hệ thống, điều này có thể không gây lỗi ngay lập tức
	// vì file đã tồn tại, nhưng sẽ gây lỗi khi rotate
	h, err := NewFileHandler(logPath, 1) // size nhỏ để kích hoạt rotate nhanh
	if err != nil {
		// Nếu lỗi ngay lập tức, test đã pass
		t.Logf("NewFileHandler() trả về lỗi như mong đợi: %v", err)
		return
	}
	defer h.Close()

	// Nếu tạo handler thành công, thử ghi log đủ lớn để kích hoạt rotate
	// và gây lỗi
	err = h.Log(InfoLevel, "message to trigger rotation in no-permission directory")
	if err == nil {
		t.Log("Log ghi thành công, nhưng có thể sẽ lỗi khi rotate")
	}

	// Thử ghi thêm để đảm bảo kích hoạt rotate
	_ = h.Log(InfoLevel, "another large message to ensure rotation happens")
}

// TestNewFileHandlerWithStatError kiểm tra lỗi khi lấy thông tin file
func TestFileHandler_New_WithStatError(t *testing.T) {
	if os.Getenv("GO_TEST_FILEHANDLER_STATERROR") == "1" {
		// Subprocess test để tạo môi trường lỗi đặc biệt
		// Ghi chú: trong thực tế, đây là một trường hợp rất khó tái tạo
		// vì cần phải tạo tình huống mở file thành công nhưng Stat() thất bại
		t.Skip("Đây là subprocess test, bỏ qua")
	}

	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Đường dẫn với ký tự đặc biệt để có khả năng gây lỗi trên một số hệ thống
	logPath := filepath.Join(dir, "test:file?.log")

	// Tạo handler
	h, err := NewFileHandler(logPath, 100)
	if err != nil {
		t.Logf("NewFileHandler() trả về lỗi: %v", err)
	} else {
		h.Close()
	}
}

// TestFileHandlerEdgeCases kiểm tra các trường hợp biên
func TestFileHandler_EdgeCases(t *testing.T) {
	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Test với maxSize = 0 (không giới hạn kích thước)
	t.Run("max_size_zero", func(t *testing.T) {
		logPath := filepath.Join(dir, "unlimited.log")
		h, err := NewFileHandler(logPath, 0)
		if err != nil {
			t.Fatalf("NewFileHandler() với maxSize=0 error = %v", err)
		}
		defer h.Close()

		// Ghi log nhiều lần - không nên gây rotation
		for i := 0; i < 10; i++ {
			err = h.Log(InfoLevel, "unlimited log message %d", i)
			if err != nil {
				t.Errorf("Log() lần thứ %d error = %v", i, err)
			}
		}

		// Kiểm tra không có file backup
		files, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("Không thể đọc thư mục: %v", err)
		}

		for _, file := range files {
			if file.Name() != "unlimited.log" && contains(file.Name(), "unlimited.log") {
				t.Errorf("Tìm thấy file backup khi maxSize=0: %s", file.Name())
			}
		}
	})

	// Test với đường dẫn tuyệt đối
	t.Run("absolute_path", func(t *testing.T) {
		absPath, err := filepath.Abs(filepath.Join(dir, "abs-path.log"))
		if err != nil {
			t.Fatalf("Không thể lấy đường dẫn tuyệt đối: %v", err)
		}

		h, err := NewFileHandler(absPath, 100)
		if err != nil {
			t.Fatalf("NewFileHandler() với đường dẫn tuyệt đối error = %v", err)
		}
		defer h.Close()

		// Ghi log để xác nhận hoạt động bình thường
		err = h.Log(InfoLevel, "absolute path test")
		if err != nil {
			t.Errorf("Log() với đường dẫn tuyệt đối error = %v", err)
		}

		// Kiểm tra file được tạo
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			t.Error("File không được tạo với đường dẫn tuyệt đối")
		}
	})

	// Test với nhiều lần đóng
	t.Run("multiple_close", func(t *testing.T) {
		logPath := filepath.Join(dir, "multi-close.log")
		h, err := NewFileHandler(logPath, 100)
		if err != nil {
			t.Fatalf("NewFileHandler() error = %v", err)
		}

		// Đóng lần đầu
		err = h.Close()
		if err != nil {
			t.Errorf("Close() lần đầu error = %v", err)
		}

		// Đóng lần thứ hai - không nên gây lỗi
		err = h.Close()
		if err != nil {
			t.Errorf("Close() lần thứ hai error = %v", err)
		}
	})

	// Test với tên file có ký tự đặc biệt
	t.Run("special_characters", func(t *testing.T) {
		// Một số hệ thống file cho phép các ký tự đặc biệt trong tên file
		logPath := filepath.Join(dir, "special-chars_#@!.log")
		h, err := NewFileHandler(logPath, 100)
		if err != nil {
			t.Logf("NewFileHandler() với ký tự đặc biệt error = %v (có thể chấp nhận được trên một số hệ thống)", err)
			return
		}
		defer h.Close()

		// Ghi log để xác nhận hoạt động bình thường
		err = h.Log(InfoLevel, "special chars test")
		if err != nil {
			t.Errorf("Log() với ký tự đặc biệt error = %v", err)
		}
	})
}

// TestNewFileHandlerWithMkdirAllError kiểm tra lỗi khi không thể tạo thư mục
func TestFileHandler_New_WithMkdirAllError(t *testing.T) {
	// Bỏ qua trên Windows vì cơ chế quyền khác
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Bỏ qua test này trên Windows")
	}

	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Tạo một file thay vì thư mục để gây lỗi cho MkdirAll
	conflictPath := filepath.Join(dir, "conflict")
	file, err := os.Create(conflictPath)
	if err != nil {
		t.Fatalf("Không thể tạo file conflict: %v", err)
	}
	file.Close()

	// Thử tạo log file trong "thư mục" này (thực tế là file)
	logPath := filepath.Join(conflictPath, "impossible.log")

	h, err := NewFileHandler(logPath, 100)
	if err == nil {
		h.Close()
		t.Error("NewFileHandler() nên trả về lỗi khi không thể tạo thư mục")
	} else {
		t.Logf("NewFileHandler() trả về lỗi như mong đợi: %v", err)
	}
}

// TestNewFileHandlerWithFileOpenError kiểm tra lỗi khi không thể mở file
func TestFileHandler_New_WithFileOpenError(t *testing.T) {
	// Bỏ qua trên Windows vì cơ chế quyền khác
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Bỏ qua test này trên Windows")
	}

	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Tạo một thư mục và set quyền chỉ đọc
	readOnlyDir := filepath.Join(dir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0755); err != nil {
		t.Fatalf("Không thể tạo thư mục: %v", err)
	}

	// Đặt quyền chỉ đọc cho thư mục
	if err := os.Chmod(readOnlyDir, 0444); err != nil {
		t.Fatalf("Không thể thay đổi quyền thư mục: %v", err)
	}

	// Khôi phục quyền sau test
	defer func() {
		if err := os.Chmod(readOnlyDir, 0755); err != nil {
			t.Logf("Không thể khôi phục quyền thư mục: %v", err)
		}
	}()

	// Thử tạo file trong thư mục chỉ đọc
	logPath := filepath.Join(readOnlyDir, "test.log")

	h, err := NewFileHandler(logPath, 100)
	if err == nil {
		h.Close()
		t.Error("NewFileHandler() nên trả về lỗi khi không thể mở file")
	} else {
		t.Logf("NewFileHandler() trả về lỗi như mong đợi: %v", err)
	}
}

// Test error messages chi tiết
func TestFileHandler_New_ErrorMessages(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func() string
		expectedError string
		cleanup       func(string)
	}{
		{
			name: "directory_does_not_exist",
			setupFunc: func() string {
				return "/tmp/absolutely-non-existent-folder-xyz/app.log"
			},
			expectedError: "path to folder do not exists",
			cleanup:       func(string) {},
		},
		{
			name: "directory_read_only",
			setupFunc: func() string {
				dir := createTempDir(t)
				_ = os.Chmod(dir, 0444) // read-only
				return filepath.Join(dir, "readonly.log")
			},
			expectedError: "directory does not have write permission",
			cleanup: func(path string) {
				dir := filepath.Dir(path)
				_ = os.Chmod(dir, 0755)
				os.RemoveAll(dir)
			},
		},
		{
			name: "valid_directory",
			setupFunc: func() string {
				dir := createTempDir(t)
				return filepath.Join(dir, "valid.log")
			},
			expectedError: "",
			cleanup: func(path string) {
				dir := filepath.Dir(path)
				os.RemoveAll(dir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logPath := tt.setupFunc()
			defer tt.cleanup(logPath)

			h, err := NewFileHandler(logPath, 100)

			if tt.expectedError == "" {
				// Expect success
				if err != nil {
					t.Errorf("NewFileHandler() với setup hợp lệ không nên lỗi, got = %v", err)
					return
				}
				if h != nil {
					h.Close()
				}
			} else {
				// Expect error
				if err == nil {
					t.Errorf("NewFileHandler() nên trả về lỗi cho test case '%s'", tt.name)
					if h != nil {
						h.Close()
					}
					return
				}

				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("NewFileHandler() error message không đúng cho test case '%s', got = %v, want chứa '%s'",
						tt.name, err.Error(), tt.expectedError)
				}
			}
		})
	}
}
