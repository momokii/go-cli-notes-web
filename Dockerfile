# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# -ldflags="-s -w" strips debug info for smaller binary
# -trimpath removes file system paths from binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -trimpath \
    -o dashboard main.go

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS (if needed for external resources)
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1000 -S dashboard && \
    adduser -u 1000 -S dashboard -G dashboard

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/dashboard .

# Copy templates and static files
COPY --chown=dashboard:dashboard templates ./templates
COPY --chown=dashboard:dashboard static ./static

# Switch to non-root user
USER dashboard

# Expose port (can be overridden at runtime)
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Set default environment variables
ENV PORT=3000
ENV ENV=production

# Run the binary
CMD ["./dashboard"]
