# ========== Build stage ==========
FROM golang:1.26.0-alpine AS builder

WORKDIR /app

# Install git (needed for some go modules) and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /api ./cmd/api

# ========== Final stage ==========
FROM alpine:3.20

WORKDIR /app

# Add ca-certificates for outbound HTTPS if needed later
RUN apk add --no-cache ca-certificates

# Create a non-root user
RUN adduser -D -g '' appuser
USER appuser

# Copy the binary from the builder stage
COPY --from=builder /api /app/api

# Expose the application port
EXPOSE 8081

# Run the binary
ENTRYPOINT ["/app/api"]