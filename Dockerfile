# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS builder
WORKDIR /src

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/app .

FROM alpine:3.19
# Create non-root user
RUN adduser -D -u 10001 appuser
USER appuser

WORKDIR /app
COPY --from=builder /out/app /app/app

# Set default environment variable for Postgres connection
ENV POSTGRES_DSN="postgres://user:password@localhost:5432/dbname?sslmode=disable"

EXPOSE 8080

ENTRYPOINT ["/app/app"]