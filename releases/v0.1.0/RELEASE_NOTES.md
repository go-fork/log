# Release Notes - v0.1.0

## Overview
Initial release of the go-fork/log package, providing comprehensive logging management capabilities for Go applications with dependency injection support.

## What's New
### 🚀 Features
- Comprehensive logging management through Manager interface
- ServiceProvider for dependency injection integration
- Multiple log levels: Debug, Info, Warning, Error, and Fatal
- Multiple handlers: Console, File, and Stack handlers
- Thread-safe concurrent logging capabilities
- Printf-style formatting support

### 🔧 Core Components
- **Manager Interface**: Core logging management with multiple handlers
- **ServiceProvider**: DI container integration for seamless setup
- **Console Handler**: Colored console output with customizable formatting
- **File Handler**: File logging with rotation capabilities
- **Stack Handler**: Multiple destination logging support

## Dependencies
### Added
- go.fork.vn/di: v0.1.0
- go.fork.vn/config: v0.1.0

## Installation
```bash
go get go.fork.vn/log@v0.1.0
```

## Contributors
Thanks to the initial development team for making this first release possible.

---
Release Date: 2025-05-31
