# Migration Guide - v0.1.4

## Overview
This guide helps you migrate from v0.1.3 to v0.1.4, which includes a complete documentation overhaul and enhanced Vietnamese language support.

## Prerequisites
- Go 1.23 or later
- Previous version (v0.1.3 or earlier) installed

## Quick Migration Checklist
- [ ] No code changes required - this is a documentation-only release
- [ ] Update documentation bookmarks to new Vietnamese docs
- [ ] Review new Fork Framework integration patterns
- [ ] Check new configuration examples for your environment
- [ ] Update project documentation references

## Breaking Changes

### Documentation Changes
**Good News**: No breaking API changes! All existing code will continue to work exactly as before.

#### Documentation Language
```markdown
# Old documentation (English)
See docs/README.md for usage instructions

# New documentation (Vietnamese)  
See docs/index.md for comprehensive Vietnamese documentation
```

#### Documentation Structure
```
# Old structure
docs/
├── README.md (basic English docs)
└── examples/ (scattered examples)

# New structure  
docs/
├── index.md         # Main hub with quick start (Vietnamese)
├── overview.md      # Architecture with Mermaid diagrams
├── configuration.md # Environment-specific configs
├── handler.md       # Complete handler documentation
├── logger.md        # Logger patterns and usage
└── workflows.md     # Integration workflows
```

### Enhanced Examples
All code examples have been updated to follow Fork Framework patterns:

```go
// Old examples (basic usage)
logger := log.NewLogger()
logger.Info("message")

// New examples (Fork Framework integration)
type UserService struct {
    logger log.Logger
}

func NewUserService(container *di.Container) *UserService {
    manager := container.Get("log").(log.Manager)
    return &UserService{
        logger: manager.GetLogger("UserService"),
    }
}
```

## Step-by-Step Migration

### Step 1: Update Dependencies (Optional)
```bash
go get go.fork.vn/log@v0.1.4
go mod tidy
```
**Note**: This step is optional as there are no API changes.

### Step 2: Update Documentation References
Update any documentation links in your projects:

```go
// Update comments in your code
// Old: See package docs at docs/README.md
// New: See comprehensive docs at docs/index.md (Vietnamese)
```

### Step 3: Review New Patterns (Optional)
Consider adopting new Fork Framework integration patterns:

```go
// Enhanced service integration pattern
func NewMyService(container *di.Container) *MyService {
    manager := container.Get("log").(log.Manager)
    return &MyService{
        logger: manager.GetLogger("MyService"), // Contextual logger
    }
}
```

### Step 4: Check Configuration Examples
Review new environment-specific configurations in `docs/configuration.md`:

```go
// Development environment
devConfig := &log.Config{
    Level: handler.DebugLevel,
    Console: log.ConsoleConfig{Enabled: true, Colored: true},
    File:    log.FileConfig{Enabled: true, Path: "logs/dev.log"},
}

// Production environment  
prodConfig := &log.Config{  
    Level: handler.InfoLevel,
    Console: log.ConsoleConfig{Enabled: false},
    File:    log.FileConfig{
        Enabled: true, 
        Path: "/var/log/app/app.log", 
        MaxSize: 100*1024*1024,
    },
}
```

### Step 5: Verify Everything Works
```bash
go test ./...
go build ./...
```

## What's New in v0.1.4

### Complete Vietnamese Documentation
- **6 comprehensive documentation files** totaling ~95KB
- **50+ practical code examples** with real-world patterns
- **12+ Mermaid diagrams** for architecture visualization
- **Environment-specific guides** (dev/prod/test/docker)

### Enhanced README.md
- Modern GitHub badges and professional presentation
- Comprehensive quick start examples
- Architecture overview with clear explanations  
- Advanced usage patterns and best practices

### Fork Framework Focus
- Deep integration examples with DI container
- Service Provider pattern implementations
- Contextual logging for microservices
- Performance monitoring workflows

## Common Questions

### Q: Do I need to change my existing code?
**A**: No! All existing APIs remain exactly the same. This is purely a documentation enhancement.

### Q: Are there new features I should use?
**A**: While no new APIs were added, the documentation now shows much better patterns for Fork Framework integration that you might want to adopt.

### Q: What about the English documentation?
**A**: The new Vietnamese documentation is much more comprehensive. English documentation may be added back in future releases based on community feedback.

## Common Issues and Solutions

### Issue 1: Cannot Find New Documentation
**Problem**: Looking for old documentation structure  
**Solution**: Use new Vietnamese docs in `docs/index.md` as the starting point

### Issue 2: Want English Documentation
**Problem**: Need English documentation for international teams  
**Solution**: The comprehensive Vietnamese docs can be translated, or refer to inline code comments which remain in English

### Issue 3: New Configuration Examples Don't Work
**Problem**: Copy-pasted new configuration examples cause issues  
**Solution**: New examples are enhanced but optional. Your existing configuration continues to work unchanged.

## Getting Help
- Check the [comprehensive Vietnamese documentation](../../docs/index.md)
- Review [architecture overview](../../docs/overview.md) with Mermaid diagrams
- Search [existing issues](https://github.com/go-fork/log/issues)
- Create a [new issue](https://github.com/go-fork/log/issues/new) for support

## Rollback Instructions
If you need to rollback to v0.1.3:

```bash
go get go.fork.vn/log@v0.1.3
go mod tidy
```

**Note**: Rollback is unlikely to be needed since no APIs changed.

## Getting Help
- Check the [documentation](https://pkg.go.dev/go.fork.vn/log@v0.1.4)
- Search [existing issues](https://github.com/go-fork/log/issues)
- Create a [new issue](https://github.com/go-fork/log/issues/new) if needed

## Rollback Instructions
If you need to rollback:

```bash
go get go.fork.vn/log@previous-version
go mod tidy
```

Replace `previous-version` with your previous version tag.

---
**Need Help?** Feel free to open an issue or discussion on GitHub.
