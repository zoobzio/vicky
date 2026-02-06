# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /vicky ./cmd/vicky

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary
COPY --from=builder /vicky /app/vicky

# Copy migrations
COPY --from=builder /app/migrations /app/migrations

# Create non-root user
RUN adduser -D -g '' vicky
USER vicky

EXPOSE 8080

ENTRYPOINT ["/app/vicky"]
