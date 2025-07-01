# Build stage
FROM golang:1.21-alpine AS builder

# Set necessary environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install git and ca-certificates (needed for go mod download)
RUN apk --no-cache add git ca-certificates

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -ldflags="-w -s" -o pmgo ./cmd/pmgo

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 pmgo && \
    adduser -D -s /bin/sh -u 1001 -G pmgo pmgo

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/pmgo .

# Copy configuration files
COPY --from=builder /build/configs ./configs

# Create necessary directories
RUN mkdir -p /app/data /app/logs && \
    chown -R pmgo:pmgo /app

# Switch to non-root user
USER pmgo

# Expose ports
EXPOSE 9876 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ./pmgo status || exit 1

# Set default command
CMD ["./pmgo", "serve", "--config", "configs/config.yaml"]