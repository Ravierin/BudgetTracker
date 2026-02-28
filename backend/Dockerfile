# Build stage
FROM golang:1.21-alpine AS builder

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

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/budget-tracker .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy .env.example as template
COPY --from=builder /app/.env.example .env

# Expose port
EXPOSE 8080

# Run application
CMD ["./budget-tracker"]
