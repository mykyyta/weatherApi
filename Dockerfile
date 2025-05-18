# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Ensure binary is built for Linux (important on macOS M1/M2)
RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

# Templates (if used by html renderer)
COPY templates/ templates/

EXPOSE 8080

# Ensure entrypoint binary is executable
RUN chmod +x ./main

CMD ["./main"]