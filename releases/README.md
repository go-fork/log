# Release Archive

This directory contains release documentation organized by version, including work-in-progress documentation.

## Structure

```
releases/
├── next/                    # ← Work-in-progress documentation
│   ├── MIGRATION.md         # Current development migration guide
│   ├── RELEASE_NOTES.md     # Current development release notes
│   └── RELEASE_SUMMARY.md   # Current development release tracking
├── v0.1.0/                  # v0.1.0 release documentation
│   ├── MIGRATION.md         # Migration guide for v0.1.0
│   ├── RELEASE_NOTES.md     # Release notes for v0.1.0
│   └── RELEASE_SUMMARY.md   # Release summary for v0.1.0
├── v0.1.1/                  # v0.1.1 release documentation
│   ├── MIGRATION.md         # Migration guide for v0.1.1
│   ├── RELEASE_NOTES.md     # Release notes for v0.1.1
│   └── RELEASE_SUMMARY.md   # Release summary for v0.1.1
└── v0.1.2/                  # v0.1.2 release documentation (current)
    ├── MIGRATION.md         # Migration guide for v0.1.2
    ├── RELEASE_NOTES.md     # Release notes for v0.1.2
    └── RELEASE_SUMMARY.md   # Release summary for v0.1.2
```

## Current Version Documentation

The current released version documentation can be found in the latest `vX.X.X/` directory.

## Current Development Documentation

The work-in-progress documentation is located in `releases/next/`:
- [releases/next/MIGRATION.md](next/MIGRATION.md) - Current development migration guide
- [releases/next/RELEASE_NOTES.md](next/RELEASE_NOTES.md) - Current development release notes
- [releases/next/RELEASE_SUMMARY.md](next/RELEASE_SUMMARY.md) - Current development release tracking
- [CHANGELOG.md](../CHANGELOG.md) - Always in root (complete change history)

## Release Workflow

### During Development
1. Work on documentation in `releases/next/`
2. Update `MIGRATION.md`, `RELEASE_NOTES.md`, and `RELEASE_SUMMARY.md` as needed
3. Track progress in `RELEASE_SUMMARY.md`

### When Ready to Release (vX.X.X)
1. **Archive Current Release**:
   ```bash
   # Create version directory
   mkdir releases/vX.X.X
   
   # Move and rename files from next/ to versioned directory
   mv releases/next/MIGRATION.md releases/vX.X.X/MIGRATION_vX.X.X.md
   mv releases/next/RELEASE_NOTES.md releases/vX.X.X/RELEASE_NOTES_vX.X.X.md
   mv releases/next/RELEASE_SUMMARY.md releases/vX.X.X/RELEASE_SUMMARY.md
   ```

2. **Create New Development Cycle**:
   ```bash
   # Create fresh templates for next version
   ./scripts/create_release_templates.sh
   ```

3. **Update CHANGELOG.md** with the released version

## Benefits of This Workflow

### For Developers
- **Clean Structure**: No confusing symlinks, clear separation of concerns
- **Direct Access**: Work directly in `releases/next/` for current development
- **Clear History**: Easy access to both current and historical documentation

### For Maintainers  
- **Streamlined Releases**: Simple workflow to archive and create new release cycle
- **No Lost Work**: All development work preserved during release process
- **Consistent Structure**: Standardized approach for all releases

### For Users
- **Clear Paths**: Direct links to documentation, no broken symlinks
- **Access to History**: Easy access to both current and historical documentation
- **Clear Versioning**: Obvious which documentation applies to which version

## Automation Scripts

You can create helper scripts in `scripts/` directory:

```bash
# scripts/create_release_templates.sh
#!/bin/bash
# Creates fresh templates in releases/next/ for new development cycle

# scripts/archive_release.sh vX.X.X  
#!/bin/bash
# Archives current release and creates new development cycle
```

This keeps the root directory clean while preserving all historical documentation and providing a smooth development workflow.
