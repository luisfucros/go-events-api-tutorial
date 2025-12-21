FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary with flags to reduce size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" -o /usr/local/bin/app ./cmd/api

FROM alpine:3.19 AS final

# Set environment variable to ensure Gin runs in release mode (performance/security)
ENV GIN_MODE=release
ENV MIGRATE_VERSION=v4.16.0

# Install runtime deps and migrate; keep it to one layer
RUN apk add --no-cache ca-certificates curl mysql-client \
 && adduser -D -g '' appuser \
 && curl -L "https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz" \
    | tar xz && mv migrate /usr/local/bin/migrate

WORKDIR /app

# Copy binary and migrations
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
COPY ./cmd/migrate/migrations /migrations
COPY scripts/entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/entrypoint.sh \
 && chown -R appuser:appuser /migrations /usr/local/bin/app

EXPOSE 8080

USER appuser

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]