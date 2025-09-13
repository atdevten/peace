# Multi-stage build for Go backend
FROM golang:1.24-alpine AS builder

# Install git, ca-certificates, and build tools
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build all services
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-w -s' -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-w -s' -o websocket-server ./cmd/websocket-server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-w -s' -o data-importer ./cmd/data-importer

# Final stage
FROM alpine:latest

# Install ca-certificates, timezone data, and postgresql-client for data-importer
RUN apk --no-cache add ca-certificates tzdata postgresql-client

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/websocket-server .
COPY --from=builder /app/data-importer .

# Copy config files
COPY --from=builder /app/configs ./configs

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8080 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the HTTP server by default
CMD ["./server"] 