# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o cardgen-pro \
    ./cmd/cardgen-pro

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/cardgen-pro .

# Create non-root user
RUN addgroup -g 1000 cardgen && \
    adduser -D -u 1000 -G cardgen cardgen

# Switch to non-root user
USER cardgen

# Expose API port (if using serve command)
EXPOSE 8080

# Default command shows help
ENTRYPOINT ["./cardgen-pro"]
CMD ["help"]

# Labels
LABEL maintainer="Felipe Macedo <felipe@example.com>"
LABEL org.opencontainers.image.title="cardgen-pro"
LABEL org.opencontainers.image.description="Card Data & ISO-8583 Test Suite"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.vendor="cardgen-pro"
LABEL org.opencontainers.image.licenses="MIT"
