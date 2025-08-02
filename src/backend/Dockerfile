# Multi-stage build for Go application
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY ./backend/go.mod ./backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source code
COPY ./backend/ .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o breach-checker-app .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/breach-checker-app .

# Copy frontend files
COPY ./frontend ./frontend

# Create directory for database
RUN mkdir -p /data

# Expose port
EXPOSE 8082

# Set environment variables
ENV PORT=8082
ENV DB_PATH=/data/breach_checker.db

# Run the application
CMD ["./breach-checker-app"]