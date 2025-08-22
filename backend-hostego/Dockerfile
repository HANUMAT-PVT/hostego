# ---------- Build stage ----------
FROM golang:1.24.3-alpine AS builder

# Set working directory
WORKDIR /app

# Install git (needed for go mod download sometimes)
RUN apk add --no-cache git

# Copy go.mod and go.sum first (to leverage Docker cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN go build -o main .

# ---------- Run stage ----------
FROM alpine:3.19

WORKDIR /app

# Copy only the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8000

# Run the app
CMD ["./main"]
