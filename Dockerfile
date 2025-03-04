# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

# Install build dependencies and UPX
RUN apk add --no-cache git upx

# Build the binary with optimization flags and compress it
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o endoflife_exporter && \
    upx --ultra-brute -qq endoflife_exporter && \
    upx -t endoflife_exporter

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/endoflife_exporter .

# Create config directory and set permissions
RUN mkdir -p /app/config && \
    chown -R 1000:1000 /app

# Create a non-root user
RUN adduser -D -u 1000 exporter
USER exporter

VOLUME /app/config

EXPOSE 2112

ENTRYPOINT ["/app/endoflife_exporter"]
CMD ["--config", "/app/config/config.yaml"] 