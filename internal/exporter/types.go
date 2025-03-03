package exporter

import (
	"io"

	"github.com/spf13/cobra"
)

type CommandFactory func(io.Writer) *cobra.Command
