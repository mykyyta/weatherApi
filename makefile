# Format all Go files with goimports and gofumpt
fmt:
	goimports -w .
	gofumpt -w .

# Run golangci-lint static code analysis
lint:
	golangci-lint run ./...

# Run all tests with verbose output
test:
	go test -v ./...

# Run tests with filtered output to reduce noise
# Removes GIN framework logs, 'record not found' errors,
# and standard test execution messages
test-quiet:
	@echo "==> Running tests quietly..."
	@go test ./... -v 2>&1 | \
		grep -v -e '^\[GIN\]' \
		        -e 'record not found' \
		        -e '^=== RUN' \

# Run all checks: formatting, linting and tests		
check: fmt lint test-quiet