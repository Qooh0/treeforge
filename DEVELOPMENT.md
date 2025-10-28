# Developer Quick Reference

## Daily Development Commands

```bash
# Before starting work
make install-tools

# During development
make check          # Run all quality checks
make test           # Run tests only
make format         # Format code

# Before committing (automatic via pre-commit hook)
git add .
git commit -m "your message"  # Runs: format + vet + test + complexity

# Manual pre-commit check
make pre-commit
```

## Troubleshooting

### Pre-commit Hook Issues
```bash
# If pre-commit hook fails, fix issues manually:
make format         # Fix formatting
go vet ./...        # Check for vet issues
go test ./...       # Ensure tests pass
make complexity     # Check complexity

# Or run all checks at once:
make check
```

### Complexity Issues
If complexity check fails:
1. Identify complex functions: `gocyclo -over 15 .`
2. Refactor using TDD approach:
   - Write tests for helper functions first
   - Extract helper functions
   - Refactor original function to use helpers

### Test Failures
```bash
# Clear test cache and re-run
go clean -testcache
go test -v ./...

# Run specific test
go test -v -run TestName
```

## File Structure

```
.
├── .git/hooks/pre-commit     # Auto-runs quality checks
├── .github/workflows/        # CI/CD pipelines
├── .vscode/settings.json     # VS Code auto-format settings
├── Makefile                  # Development commands
├── DEVELOPMENT.md           # This file
├── README.md                # Main documentation
├── main_test.go             # Main function tests
├── parse_tree.go            # Core parsing logic
├── parse_tree_test.go       # Parser tests
└── treeforge.go             # CLI entry point
```

## Code Quality Targets

- **Test Coverage**: 100%
- **Cyclomatic Complexity**: ≤ 15 per function
- **Code Formatting**: gofmt -s compliant
- **Static Analysis**: go vet clean

## Release Process

1. Ensure all tests pass: `make check`
2. Update version in code if needed
3. Commit changes: `git commit -m "..."`
4. Tag release: `git tag v0.x.x`
5. Push: `git push origin main --tags`
6. GitHub Actions automatically creates release