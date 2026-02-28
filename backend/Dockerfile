# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o budget-tracker ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS and postgresql-client for migrations
RUN apk --no-cache add ca-certificates postgresql-client

# Copy binary from builder
COPY --from=builder /app/budget-tracker .

# Copy entrypoint script
COPY --from=builder /app/entrypoint.sh .
RUN chmod +x /app/entrypoint.sh

# Expose port
EXPOSE 8080

# Run entrypoint script
CMD ["./entrypoint.sh"]
