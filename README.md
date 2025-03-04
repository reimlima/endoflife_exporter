# End of Life Exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/reimlima/endoflife_exporter)](https://goreportcard.com/report/github.com/reimlima/endoflife_exporter)
[![Go Version](https://img.shields.io/github/go-mod/go-version/reimlima/endoflife_exporter)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Test Coverage](https://img.shields.io/badge/coverage-77%25-brightgreen.svg)](https://github.com/reimlima/endoflife_exporter/actions)
[![Docker Pulls](https://img.shields.io/docker/pulls/reimlima/endoflife_exporter)](https://hub.docker.com/r/reimlima/endoflife_exporter)

A Prometheus exporter that collects end-of-life dates for various products using the [endoflife.date](https://endoflife.date) API.

## Features

- Collects end-of-life dates for multiple products
- Configurable through YAML configuration file
- Prometheus metrics format
- Support for various products (Ubuntu, NodeJS, etc.)
- Flexible date handling for different date formats
- Docker support
- Systemd service support

## Installation

### Using Go

```bash
go install github.com/reimlima/endoflife_exporter@latest
```

### Using Docker

```bash
docker pull reimlima/endoflife_exporter:latest
```

### Building from source

```bash
git clone https://github.com/reimlima/endoflife_exporter.git
cd endoflife_exporter
go build
```

### Using systemd (Linux)

1. Create the prometheus user and group:
```bash
sudo groupadd -r prometheus
sudo useradd -r -g prometheus -s /sbin/nologin prometheus
```

2. Copy the binary to the system:
```bash
sudo cp endoflife_exporter /usr/local/bin/
sudo chown prometheus:prometheus /usr/local/bin/endoflife_exporter
```

3. Create configuration directories:
```bash
sudo mkdir -p /etc/endoflife_exporter
sudo mkdir -p /var/lib/endoflife_exporter
sudo chown -R prometheus:prometheus /etc/endoflife_exporter
sudo chown -R prometheus:prometheus /var/lib/endoflife_exporter
```

4. Copy the systemd service file:
```bash
sudo cp endoflife_exporter.service /etc/systemd/system/
```

5. Enable and start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable endoflife_exporter
sudo systemctl start endoflife_exporter
```

## Usage

### Configuration

Create a configuration file (e.g., `config.yaml`):

```yaml
port: 2112
products:
  - ubuntu:
      host: localhost
      version: "22.04"
  - nodejs:
      host: localhost
      version: "16"
```

### Running with Go binary

```bash
endoflife_exporter --config config.yaml
```

### Running with Docker

1. Create a config directory and add your config file:
```bash
mkdir -p config
cp config.yaml config/
```

2. Run the container:
```bash
docker run -d \
  --name endoflife_exporter \
  -p 2112:2112 \
  -v $(pwd)/config:/app/config \
  reimlima/endoflife_exporter:latest
```

### Running with Docker Compose

The project includes a `docker-compose.yml` that sets up both the exporter and Prometheus for easy deployment:

1. Create the necessary directories and files:
```bash
mkdir -p config prometheus
cp config.yaml config/
```

2. Create `prometheus/prometheus.yml`:
```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'endoflife_exporter'
    static_configs:
      - targets: ['endoflife_exporter:2112']
```

3. Start the services:
```bash
docker-compose up -d
```

This will start:
- End of Life Exporter on port 2112 (metrics endpoint)
- Prometheus on port 9090 (web interface)

You can access:
- Metrics: http://localhost:2112/metrics
- Prometheus UI: http://localhost:9090

### Running with systemd

1. Place your configuration file:
```bash
sudo cp config.yaml /etc/endoflife_exporter/
sudo chown prometheus:prometheus /etc/endoflife_exporter/config.yaml
```

2. Start the service:
```bash
sudo systemctl start endoflife_exporter
```

3. Check the status:
```bash
sudo systemctl status endoflife_exporter
```

## Metrics

The exporter provides the following metrics:

- `endoflife_service`: End of life date for services (Unix timestamp)
  - Labels:
    - `service`: Product name
    - `host`: Host where the product is running
    - `version`: Product version
    - `cycle`: Release cycle

Example:
```
# HELP endoflife_service End of life date for services
# TYPE endoflife_service gauge
endoflife_service{cycle="22.04",host="localhost",service="ubuntu",version="22.04"} 1.682899e+09
```

## Grafana Dashboard

This repository contains a Grafana dashboard to show collected metrics. Take a look at [grafana_dashboard](./grafana_dashboard/dashboard.md) documentation to know more.

## Development

Requirements:
- Go 1.21 or higher
- Make (optional)
- Docker (optional)

Running tests:
```bash
go test -v ./...
```

Running tests with coverage:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Building Docker image locally:
```bash
docker build -t endoflife_exporter .
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [endoflife.date](https://endoflife.date) for providing the API
- [Prometheus](https://prometheus.io) for the monitoring framework
