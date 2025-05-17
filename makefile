# Format all Go code using goimports (fix imports) and gofumpt (stricter formatting)
fmt:
	goimports -w .
	gofumpt -w .

# Run static analysis and linting using golangci-lint
# Make sure .golangci.yml is configured properly
lint:
	golangci-lint run ./...

# Run all Go tests with verbose output
test:
	go test -v ./...

# Combined command: format, lint, and test (used for pre-commit or CI)
check: fmt lint test