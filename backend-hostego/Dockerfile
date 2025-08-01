# ---------- Step 1: Build Go binary ----------
FROM golang:latest AS builder

WORKDIR /app

# Copy go mod files and download deps first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the source code
COPY . .

# Build the Go binary for Alpine (static)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o backend-hostego .

# ---------- Step 2: Minimal Alpine Runtime ----------
FROM alpine:latest

WORKDIR /root

# Install certs for HTTPS (optional but usually required)
RUN apk add --no-cache ca-certificates

# Copy compiled binary from builder
COPY --from=builder /app/backend-hostego .

# Copy wait-for script for database readiness
COPY wait-for.sh /wait-for.sh
RUN chmod +x /wait-for.sh

# Copy environment file if needed (optional, or handled via docker-compose)
COPY .env .env

# Let Firebase JSON be mounted via docker-compose volume, so no need to copy here

# Expose application port
EXPOSE 8080

# ENTRYPOINT (uses service name like postgres, not IP)
ENTRYPOINT [ "/wait-for.sh", "172.17.0.1:5432", "--", "./backend-hostego" ]
