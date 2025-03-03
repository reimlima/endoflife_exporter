package internal

import (
	"io"

	"github.com/reimlima/endoflife_exporter/internal/config"
	"github.com/reimlima/endoflife_exporter/internal/exporter"
	"github.com/spf13/cobra"
)

func BuildRootCommand(out io.Writer, factory ...exporter.CommandFactory) *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "endoflife_exporter",
		Short: "Prometheus exporter for end-of-life API",
		Long: `A Prometheus exporter that collects end-of-life dates for various products.
		
Example usage:
  endoflife_exporter --config config.yaml

The config file should be in YAML format with the following structure:
  port: 2112
  products:
    - ubuntu:
        host: localhost
        version: "22.04"
    - nodejs:
        host: localhost
        version: "16"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.SetConfig(configPath)
			return exporter.StartExporter(cfg)
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to the config file")
	return cmd
}

func BuildRootScope(out io.Writer) *cobra.Command {
	scope := make([]exporter.CommandFactory, 0)
	scope = append(scope, func(out io.Writer) *cobra.Command {
		return BuildRootCommand(out)
	})
	return BuildRootCommand(out, scope...)
}
