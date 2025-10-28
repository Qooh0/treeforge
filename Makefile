# Makefile for treeforge

.PHONY: help check format test build clean install-tools pre-commit

# Default target - show help
help:
	@echo "🛠️  Treeforge Development Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make install-tools  Install development dependencies (gocyclo, etc.)"
	@echo ""
	@echo "Development:"
	@echo "  make check          Run all quality checks (format + vet + test + complexity)"
	@echo "  make format         Format code with gofmt -s"
	@echo "  make vet           Run go vet static analysis"
	@echo "  make test          Run test suite"
	@echo "  make complexity    Check cyclomatic complexity"
	@echo "  make pre-commit    Run pre-commit checks"
	@echo ""
	@echo "Build:"
	@echo "  make build         Build the binary"
	@echo "  make clean         Clean build artifacts and test cache"
	@echo ""
	@echo "💡 Tip: Pre-commit hooks automatically run 'make check' on every commit"

.PHONY: check format test build clean install-tools pre-commit

# Install development tools
install-tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

# Format code
format:
	gofmt -s -w .

# Run all checks
check: format vet test complexity

# Run go vet
vet:
	go vet ./...

# Run tests
test:
	go test -v ./...

# Check cyclomatic complexity
complexity:
	@if command -v gocyclo >/dev/null 2>&1; then \
		if [ "$$(gocyclo -over 15 . | wc -l)" -gt 0 ]; then \
			echo "High cyclomatic complexity detected:"; \
			gocyclo -over 15 .; \
			exit 1; \
		else \
			echo "Cyclomatic complexity check passed"; \
		fi \
	else \
		echo "gocyclo not found. Run 'make install-tools' first"; \
		exit 1; \
	fi

# Pre-commit checks (run before committing)
pre-commit: check
	@echo "All pre-commit checks passed!"

# Build the application
build:
	go build -o treeforge .

# Clean build artifacts and test cache
clean:
	go clean
	go clean -testcache
	rm -f treeforge

# Default target
all: help