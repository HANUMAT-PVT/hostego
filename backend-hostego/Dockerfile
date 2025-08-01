# Step 1: Build the Go binary statically
FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ⛔️ No need to copy .env here — it's for runtime
# COPY .env .env  ❌

# 🔐 Build statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o backend-hostego .

# Step 2: Use minimal base image
FROM alpine:latest

WORKDIR /root

# Copy the binary
COPY --from=builder /app/backend-hostego .



# ✅ Copy .env file into runtime container
COPY .env .env

# Copy firebase service account file for dev

# COPY config/firebase-service-account.json config/firebase-service-account.json

# Expose port
EXPOSE 8000

# Start app
CMD ["./backend-hostego"]
