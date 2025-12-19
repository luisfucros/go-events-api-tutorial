FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/api

FROM alpine:3.19 AS final

# Set environment variable to ensure Gin runs in release mode for performance/security
ENV GIN_MODE release

RUN apk add --no-cache \
    ca-certificates \
    curl \
 && adduser -D -g '' appuser

RUN apk add --no-cache curl ca-certificates \
&& curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz \
| tar xz \
&& mv migrate /usr/local/bin/migrate

 # Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app ./
COPY ./cmd/migrate/migrations ./cmd/migrate/migrations

# Expose the port your Gin application listens on (default is 8080)
EXPOSE 8080

USER appuser

# Command to run the executable when the container starts
ENTRYPOINT ["./app"]