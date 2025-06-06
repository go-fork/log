# Tổng Quan Kiến Trúc

Package log được thiết kế đặc biệt cho Fork Framework, cung cấp hệ thống logging hiệu suất cao với kiến trúc shared handlers và contextual loggers.

## Kiến Trúc Tổng Quan

```mermaid
graph TB
    subgraph "Fork Framework Application"
        App[Fork App]
        Container[DI Container]
        Provider[LogServiceProvider]
    end
    
    subgraph "Log Package"
        Config[Config]
        Manager[Manager]
        
        subgraph "Contextual Loggers"
            UserLogger[UserService Logger]
            OrderLogger[OrderService Logger]
            PaymentLogger[PaymentService Logger]
            AuthLogger[AuthService Logger]
        end
        
        subgraph "Shared Handlers"
            Console[Console Handler]
            File[File Handler]
            Stack[Stack Handler]
        end
    end
    
    App --> Container
    Container --> Provider
    Provider --> Manager
    Config --> Manager
    
    Manager --> UserLogger
    Manager --> OrderLogger
    Manager --> PaymentLogger
    Manager --> AuthLogger
    
    UserLogger --> Console
    UserLogger --> File
    UserLogger --> Stack
    
    OrderLogger --> Console
    OrderLogger --> File
    OrderLogger --> Stack
    
    PaymentLogger --> Console
    PaymentLogger --> File
    PaymentLogger --> Stack
    
    AuthLogger --> Console
    AuthLogger --> File
    AuthLogger --> Stack
```

## Các Thành Phần Chính

### 1. Manager (Quản Lý Tập Trung)

Manager là thành phần trung tâm thực hiện pattern **Shared Handlers Architecture**:

- **Shared Handlers**: Tất cả loggers cùng chia sẻ các handler instances
- **GetOrCreate Pattern**: Logger được tạo tự động theo context khi chưa tồn tại
- **Runtime Management**: Quản lý handlers và loggers trong runtime
- **Resource Efficiency**: Tránh duplicate handlers, tiết kiệm tài nguyên

```mermaid
graph LR
    subgraph "Manager Responsibilities"
        A[Handler Management] --> B[Logger Creation]
        B --> C[Resource Sharing]
        C --> D[Lifecycle Management]
    end
```

### 2. Logger (Contextual Logging)

Mỗi logger được gắn với một context cụ thể:

- **Context-Based**: Logger được định danh bởi context (vd: "UserService", "OrderService")
- **Structured Logging**: Hỗ trợ structured logs với key-value pairs
- **Level Filtering**: Logs được filter theo level được cấu hình
- **Handler Delegation**: Ủy thác việc ghi log cho các handlers

```mermaid
graph TD
    Logger[Logger Instance] --> |Debug| Handler1[Handler 1]
    Logger --> |Info| Handler1
    Logger --> |Warning| Handler2[Handler 2]
    Logger --> |Error| Handler2
    Logger --> |Fatal| Handler3[Handler 3]
    
    Handler1 --> Console[Console Output]
    Handler2 --> File[File Output]
    Handler3 --> Stack[Stack Output]
```

### 3. Handlers (Output Processors)

Handlers xử lý việc xuất logs ra các đích khác nhau:

#### Console Handler
- Xuất logs ra stdout/stderr
- Hỗ trợ color coding theo level
- Tối ưu cho development environment

#### File Handler  
- Ghi logs vào file
- Hỗ trợ file rotation theo size
- Tối ưu cho production logging

#### Stack Handler
- Kết hợp nhiều handlers khác
- Cho phép log cùng lúc ra nhiều đích
- Linh hoạt trong cấu hình

```mermaid
graph TB
    subgraph "Handler Types"
        CH[Console Handler]
        FH[File Handler]
        SH[Stack Handler]
    end
    
    subgraph "Outputs"
        Stdout[Standard Output]
        File[Log Files]
        Multiple[Multiple Destinations]
    end
    
    CH --> Stdout
    FH --> File
    SH --> Multiple
```

## Tích Hợp Fork Framework

### Service Provider Pattern

Package log tích hợp với Fork Framework thông qua `LogServiceProvider`:

```mermaid
sequenceDiagram
    participant App as Fork App
    participant Container as DI Container
    participant Provider as LogServiceProvider
    participant Manager as Log Manager
    
    App->>Container: Bootstrap
    Container->>Provider: Register
    Provider->>Container: Bind("log", Manager)
    Container->>Manager: Create Instance
    Manager->>Manager: Initialize Handlers
```

### Dependency Injection

Logger có thể được inject vào các service khác:

```go
type UserService struct {
    logger log.Logger
}

func NewUserService(container *container.Container) *UserService {
    manager := container.Get("log").(log.Manager)
    return &UserService{
        logger: manager.GetLogger("UserService"),
    }
}
```

## Workflow Chuẩn

### 1. Khởi Tạo Application

```mermaid
sequenceDiagram
    participant Config as Configuration
    participant Provider as LogServiceProvider
    participant Manager as Log Manager
    participant Handlers as Shared Handlers
    
    Config->>Provider: Load Config
    Provider->>Manager: Create with Config
    Manager->>Handlers: Initialize All Handlers
    Handlers->>Manager: Register Handlers
```

### 2. Runtime Logging

```mermaid
sequenceDiagram
    participant Service as Business Service
    participant Manager as Log Manager
    participant Logger as Context Logger
    participant Handlers as Shared Handlers
    
    Service->>Manager: GetLogger("ServiceName")
    Manager->>Logger: Create/Return Logger
    Service->>Logger: Log Message
    Logger->>Handlers: Delegate to Handlers
    Handlers->>Handlers: Process & Output
```

## Tính Năng Nâng Cao

### 1. Runtime Handler Management

```go
// Thêm handler mới trong runtime
manager.AddHandler(log.HandlerTypeCustom, customHandler)

// Cấu hình handler cho logger cụ thể
manager.SetHandler("PaymentService", log.HandlerTypeCustom)

// Gỡ bỏ handler
manager.RemoveHandler(log.HandlerTypeCustom)
```

### 2. Context Switching

```go
// Logger có thể thay đổi context
userLogger := manager.GetLogger("UserService")
userLogger.SetContext("UserService::Registration")

// Log sẽ hiển thị context mới
userLogger.Info("User registration started")
// Output: [INFO] [UserService::Registration] User registration started
```

### 3. Performance Optimization

- **Handler Reuse**: Handlers được chia sẻ giữa các loggers
- **Level Filtering**: Logs được filter sớm để tránh xử lý không cần thiết
- **Concurrent Safe**: Thread-safe cho các ứng dụng concurrent
- **Memory Efficient**: Tối ưu memory usage thông qua shared resources

## Patterns Được Hỗ Trợ

### 1. Singleton Pattern (Manager)
- Một Manager instance duy nhất trong application
- Quản lý tất cả handlers và loggers

### 2. Factory Pattern (Logger Creation)
- Manager tạo loggers theo context
- GetOrCreate pattern đảm bảo không duplicate

### 3. Strategy Pattern (Handlers)
- Các handler implementations khác nhau
- Có thể swap handlers trong runtime

### 4. Observer Pattern (Stack Handler)
- Stack handler notify cho multiple sub-handlers
- Loose coupling giữa loggers và output destinations

## So Sánh Với Các Giải Pháp Khác

| Tính Năng | Fork Log | Logrus | Zap | Go Log |
|-----------|----------|--------|-----|---------|
| Shared Handlers | ✅ | ❌ | ❌ | ❌ |
| Contextual Loggers | ✅ | Partial | Partial | ❌ |
| Fork Framework Integration | ✅ | ❌ | ❌ | ❌ |
| Runtime Handler Management | ✅ | Limited | Limited | ❌ |
| Zero Allocation | Partial | ❌ | ✅ | ❌ |
| Structured Logging | ✅ | ✅ | ✅ | ❌ |

Thiết kế này đảm bảo package log phù hợp hoàn hảo với kiến trúc và triết lý của Fork Framework.
