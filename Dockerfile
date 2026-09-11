# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /zemlygo cmd/main.go

# Stage 2: Runtime
FROM alpine:3.20

WORKDIR /app

# Add ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /zemlygo /app/zemlygo
COPY --from=builder /app/.env.example /app/.env.example

EXPOSE 8080

CMD ["/app/zemlygo"]
