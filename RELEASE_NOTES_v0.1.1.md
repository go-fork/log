# Log Package v0.1.1 Release Notes

## Summary
Version v0.1.1 của Log Package cung cấp các cải tiến quan trọng về hiệu suất và cập nhật dependencies. Đây là một bản phát hành bảo trì nhỏ không thêm tính năng mới, nhưng cải thiện độ ổn định và khả năng tích hợp với các module khác trong hệ sinh thái Go-Fork.

## Changes

### Dependencies
- Cập nhật `go.fork.vn/config` từ v0.1.0 lên v0.1.1
- Cập nhật `go.fork.vn/di` từ v0.1.0 lên v0.1.1

### Performance
- Cải thiện hiệu suất của các handler logger
- Tối ưu hóa quản lý bộ nhớ khi xử lý log messages với volume cao
- Giảm thiểu độ trễ khi ghi log trong các trường hợp đa luồng

### Documentation
- Cập nhật tham chiếu phiên bản trong tất cả tài liệu
- Thêm file MIGRATION_v0.1.1.md với hướng dẫn nâng cấp

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
Vui lòng báo cáo bất kỳ vấn đề nào tại: https://github.com/go-fork/log/issues
