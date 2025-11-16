# ==========================================
# Stage 1: Builder
# ==========================================
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum .
RUN go mod download

# Copy source code
COPY . .

# Install CA certificates for SSL support
RUN apk add --no-cache ca-certificates

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/server ./cmd/server/main.go

# Build the migrator binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/migrator ./cmd/migrator/main.go

# ==========================================
# Stage 2: Server Target
# ==========================================
FROM alpine:latest AS server
WORKDIR /app

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy server binary
COPY --from=builder /app/build/server .

# Run the server
CMD ["./server"]

# ==========================================
# Stage 3: Migrator Target
# ==========================================
FROM alpine:latest AS migrator
WORKDIR /app

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy migrator binary
COPY --from=builder /app/build/migrator .

# Run the migrator
ENTRYPOINT ["./migrator"]
