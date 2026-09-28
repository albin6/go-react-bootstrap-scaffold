# Stage 1: Build the Go executable binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download dependencies first (leverages Docker cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source code
COPY . .

# Compile a statically linked binary with stripped debug symbols
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -trimpath \
    -o /bin/server \
    ./cmd/server/main.go

# Stage 2: Minimal, secure runtime container
FROM alpine:3.20

# Create a non-root system user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy root CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy compiled binary from builder
COPY --from=builder /bin/server /bin/server

# Execute as non-root user
USER appuser:appgroup

EXPOSE 8080

ENTRYPOINT ["/bin/server"]
