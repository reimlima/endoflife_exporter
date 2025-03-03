# End of Life Exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/reimlima/endoflife_exporter)](https://goreportcard.com/report/github.com/reimlima/endoflife_exporter)
[![Go Version](https://img.shields.io/github/go-mod/go-version/reimlima/endoflife_exporter)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Test Coverage](https://img.shields.io/badge/coverage-77%25-brightgreen.svg)](https://github.com/reimlima/endoflife_exporter/actions)

A Prometheus exporter that collects end-of-life dates for various products using the [endoflife.date](https://endoflife.date) API.

## Features

- Collects end-of-life dates for multiple products
- Configurable through YAML configuration file
- Prometheus metrics format
- Support for various products (Ubuntu, NodeJS, etc.)
- Flexible date handling for different date formats

## Installation

```bash
go install github.com/reimlima/endoflife_exporter@latest
```

Or build from source:

```bash
git clone https://github.com/reimlima/endoflife_exporter.git
cd endoflife_exporter
go build
```

## Usage

1. Create a configuration file (e.g., `config.yaml`):

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

2. Run the exporter:

```bash
endoflife_exporter --config config.yaml
```

Or with a custom config path:

```bash
endoflife_exporter -c /path/to/config.yaml
```

## Configuration

The configuration file supports the following options:

| Field | Type | Description |
|-------|------|-------------|
| port | int | The port number for the metrics server (e.g., 2112) |
| products | array | List of products to monitor |
| products[].name | string | Product name (e.g., ubuntu, nodejs) |
| products[].host | string | Host where the product is running |
| products[].version | string | Product version to monitor |

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

## Development

Requirements:
- Go 1.21 or higher
- Make (optional)

Running tests:
```bash
go test -v ./...
```

Running tests with coverage:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
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
