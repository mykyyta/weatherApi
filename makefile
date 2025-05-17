fmt:
	goimports -w .
	gofumpt -w .

lint:
	golangci-lint run ./...

test:
	go test -v ./...

check: fmt lint test