package exporter

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/reimlima/endoflife_exporter/internal/config"
	"github.com/reimlima/endoflife_exporter/internal/eolapi"
)

var (
	endOfLifeMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "endoflife_service",
		Help: "End of life date for services",
	}, []string{"service", "host", "version", "cycle"})
)

func RegisterMetrics(cfg config.Config, client *eolapi.Client) error {
	for _, product := range cfg.Products {
		for name, details := range product {
			data, err := client.FetchEOLData(name)
			if err != nil {
				return fmt.Errorf("error fetching EOL data: %v", err)
			}

			for _, item := range data {
				switch v := item.LTS.(type) {
				case bool:
					item.LTS = fmt.Sprintf("%t", v)
				case string:
					item.LTS = v
				default:
					item.LTS = "unknown"
				}

				eolStr := string(item.EOL)
				if eolStr == "true" || eolStr == "false" {
					log.Printf("Skipping non-date EOL value for %s: %s", name, eolStr)
					continue
				}

				// Try to parse the EOL date
				timestamp, err := time.Parse("2006-01-02", eolStr)
				if err != nil {
					log.Printf("Warning: Invalid date format for %s EOL: %s", name, eolStr)
					continue
				}

				// Only set metric if we have a valid date
				endOfLifeMetric.WithLabelValues(
					name,
					details.Host,
					details.Version,
					string(item.Cycle),
				).Set(float64(timestamp.Unix()))
			}
		}
	}
	return nil
}

func StartHTTPServer(port int) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func StartExporter(cfg config.Config) error {
	client := eolapi.NewClient()
	if err := RegisterMetrics(cfg, client); err != nil {
		return err
	}
	return StartHTTPServer(cfg.Port)
}
