package log

// HandlerType định nghĩa loại handler trong hệ thống logging.
type HandlerType string

// Các constants cho các loại handler được hỗ trợ.
var (
	HandlerTypeConsole HandlerType = "console"
	HandlerTypeFile    HandlerType = "file"
	HandlerTypeStack   HandlerType = "stack"
)
