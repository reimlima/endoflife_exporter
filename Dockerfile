# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

# Install build dependencies
RUN apk add --no-cache git

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o endoflife_exporter

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/endoflife_exporter .

# Create a non-root user
RUN adduser -D -u 1000 exporter
USER exporter

# Create config directory
RUN mkdir -p /app/config
VOLUME /app/config

EXPOSE 2112

ENTRYPOINT ["/app/endoflife_exporter"]
CMD ["--config", "/app/config/config.yaml"] 