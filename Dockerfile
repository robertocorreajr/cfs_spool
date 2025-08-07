# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies for CGO (required for ebfe/scard)
RUN apk add --no-cache gcc musl-dev pkgconfig pcsc-lite-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the web server
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o cfs-spool-web cmd/web-server/main.go

# Build the CLI tool
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o cfs-spool-cli cmd/cfs-spool/main.go

# Runtime stage
FROM alpine:latest

# Install PC/SC lite for RFID communication
RUN apk add --no-cache pcsc-lite pcsc-lite-libs

WORKDIR /app

# Copy built binaries
COPY --from=builder /app/cfs-spool-web .
COPY --from=builder /app/cfs-spool-cli .
COPY --from=builder /app/web ./web/

EXPOSE 8080

# Start PC/SC daemon and web server
CMD ["sh", "-c", "pcscd && ./cfs-spool-web"]
