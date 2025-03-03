package main

import (
	"fmt"
	"os"

	"github.com/reimlima/endoflife_exporter/internal"
)

func main() {
	instance := internal.BuildRootScope(os.Stdout)
	if err := instance.Execute(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
