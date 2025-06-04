# Log Package v0.1.1 Release Notes

## Summary
Version v0.1.1 là một bản phát hành bảo trì tập trung vào cập nhật dependencies và cải thiện độ ổn định cùng hiệu suất của hệ thống logging.

## Changes

### Changed
- Updated dependencies to latest versions (go.fork.vn/config v0.1.0 → v0.1.1, go.fork.vn/di v0.1.0 → v0.1.1)
- Enhanced stability and performance improvements
- Improved logger handler performance and memory optimization for high-volume logging
- Reduced latency in multi-threaded logging scenarios

### Added
- Comprehensive test suite for Config validation and error handling
- Enhanced test coverage and stability
- Better integration with updated dependencies

## Compatibility
Phiên bản này hoàn toàn tương thích ngược với v0.1.0. Không có thay đổi API và không cần thực hiện bất kỳ điều chỉnh nào đối với code hiện có khi nâng cấp.

## Installation
```bash
go get go.fork.vn/log@v0.1.1
```

hoặc cập nhật trong file go.mod:
```go
require go.fork.vn/log v0.1.1
```

## Support
Vui lòng báo cáo bất kỳ vấn đề nào tại: github.com/go-fork/log/issues
