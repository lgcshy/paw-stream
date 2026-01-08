# Multi-stage build for minimal image size

# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with CGO disabled for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-s -w" \
    -o api ./cmd/api

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /build/api .

# Copy migrations
COPY migrations ./migrations

# Create directories for data and logs
RUN mkdir -p /app/data /app/logs

# Expose port
EXPOSE 3000

# Run
CMD ["./api"]
