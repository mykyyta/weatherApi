# Build stage (x86_64)
FROM --platform=linux/amd64 golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

# Runtime stage (distroless, no shell)
FROM --platform=linux/amd64 gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/main .
COPY templates/ templates/
USER nonroot
EXPOSE 8080
CMD ["/app/main"]