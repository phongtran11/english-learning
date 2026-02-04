# Base stage
FROM golang:1.25-alpine AS base
WORKDIR /app
RUN apk add --no-cache git

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Dev stage (with hot reload)
FROM base AS dev
RUN go install github.com/air-verse/air@latest
COPY . .
CMD ["air", "-c", ".air.toml"]

# Builder stage
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# Production stage
FROM alpine:3.21 AS prod
WORKDIR /root/
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs
# Note: In production, .env should not be copied if you use real env vars, 
# but for simplicity we assume env vars are injected or .env is present.

EXPOSE 8080
CMD ["./server"]
