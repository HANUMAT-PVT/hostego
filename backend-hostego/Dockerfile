# ---------- Step 1: Build Go binary ----------
FROM golang:latest AS builder

WORKDIR /app

# Only copy go.mod/sum first for caching layer
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the application
COPY . .

# Statically compile Go binary (no CGO, Alpine compatible)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o backend-hostego .

# ---------- Step 2: Minimal runtime ----------
FROM alpine:latest

WORKDIR /root

# Copy the built Go binary
COPY --from=builder /app/backend-hostego .

# Copy .env for runtime config (if needed)
COPY .env .env

# Copy Firebase service account file
# COPY config/firebase-service-account.json config/firebase-service-account.json

# COPY /etc/hostego/firebase.json /etc/hostego/firebase.json


# Add wait-for script for postgres readiness
COPY wait-for.sh /wait-for.sh
RUN chmod +x /wait-for.sh

EXPOSE 8080

# Use wait-for to delay start until DB is ready
ENTRYPOINT [ "/wait-for.sh","172.17.0.1:5432","--","./backend-hostego"]
