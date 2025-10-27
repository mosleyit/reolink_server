# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/bin/reolink-server cmd/server/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 reolink && \
    adduser -D -u 1000 -G reolink reolink

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/reolink-server /app/reolink-server

# Copy configuration and web files
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/web /app/web

# Change ownership
RUN chown -R reolink:reolink /app

# Switch to non-root user
USER reolink

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
ENTRYPOINT ["/app/reolink-server"]

