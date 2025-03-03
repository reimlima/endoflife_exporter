package config

import (
	"os"
	"strings"
	"testing"
)

func TestSetConfig(t *testing.T) {
	// Create a temporary config file
	content := []byte(`
port: 2112
products:
  - ubuntu:
      host: localhost
      version: "22.04"
`)
	tmpfile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test valid config
	cfg := SetConfig(tmpfile.Name())
	if cfg.Port != 2112 {
		t.Errorf("Expected port 2112, got %d", cfg.Port)
	}
	if len(cfg.Products) != 1 {
		t.Errorf("Expected 1 product, got %d", len(cfg.Products))
	}
}

func TestSetConfigInvalidFile(t *testing.T) {
	// Create a temporary file with invalid YAML
	content := []byte(`
port: invalid
products: invalid yaml content
`)
	tmpfile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if r := recover(); r != nil {
			// Check if the panic message contains the expected error
			panicMsg, ok := r.(string)
			if !ok || !strings.Contains(panicMsg, "Error reading config file") {
				t.Errorf("Expected panic with 'Error reading config file', got: %v", r)
			}
		} else {
			t.Error("Expected SetConfig to panic with invalid config file")
		}
	}()

	SetConfig(tmpfile.Name())
}
