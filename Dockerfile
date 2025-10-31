FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN apk add ca-certificates

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/server ./cmd/server/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/migrator ./cmd/migration/main.go

# Target untuk server
FROM alpine:latest AS server
WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/build/server .

CMD ["./server"]

# Target untuk migrator
FROM alpine:latest AS migrator
WORKDIR /app

# COPY --from=builder /app/internal/app/database/migration internal/app/database/migration

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/build/migrator .

ENTRYPOINT ["./migrator"]
