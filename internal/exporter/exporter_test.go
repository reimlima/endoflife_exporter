package exporter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/reimlima/endoflife_exporter/internal/config"
	"github.com/reimlima/endoflife_exporter/internal/eolapi"
)

func TestRegisterMetrics(t *testing.T) {
	tests := []struct {
		name     string
		mockData []eolapi.EOLData
		config   config.Config
		wantErr  bool
	}{
		{
			name: "valid date EOL",
			mockData: []eolapi.EOLData{
				{
					Cycle:        eolapi.FlexibleString("21.04"),
					LTS:          "false",
					ReleaseDate:  eolapi.FlexibleDate("2021-04-22"),
					Support:      eolapi.FlexibleDate("2022-01-01"),
					EOL:          eolapi.FlexibleString("2022-01-01"),
					Latest:       eolapi.FlexibleString("21.04"),
					Link:         eolapi.FlexibleString("https://wiki.ubuntu.com/HirsuteHippo/ReleaseNotes/"),
					Discontinued: eolapi.FlexibleDate("2022-01-01"),
				},
			},
			config: config.Config{
				Port: 2112,
				Products: []map[string]config.Product{
					{
						"spring-framework": {Host: "localhost", Version: "3.3"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "boolean EOL value",
			mockData: []eolapi.EOLData{
				{
					Cycle:        eolapi.FlexibleString("22.04"),
					LTS:          true,
					ReleaseDate:  eolapi.FlexibleDate("2022-04-22"),
					Support:      eolapi.FlexibleDate("2023-01-01"),
					EOL:          eolapi.FlexibleString("true"),
					Latest:       eolapi.FlexibleString("22.04"),
					Link:         eolapi.FlexibleString(""),
					Discontinued: eolapi.FlexibleDate("false"),
				},
			},
			config: config.Config{
				Port: 2112,
				Products: []map[string]config.Product{
					{
						"ubuntu": {Host: "localhost", Version: "22.04"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid date format",
			mockData: []eolapi.EOLData{
				{
					Cycle:       eolapi.FlexibleString("22.04"),
					LTS:         true,
					EOL:         eolapi.FlexibleString("invalid-date"),
					Latest:      eolapi.FlexibleString("22.04"),
					ReleaseDate: eolapi.FlexibleDate("2022-04-22"),
				},
			},
			config: config.Config{
				Port: 2112,
				Products: []map[string]config.Product{
					{
						"ubuntu": {Host: "localhost", Version: "22.04"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple products",
			mockData: []eolapi.EOLData{
				{
					Cycle:       eolapi.FlexibleString("21.04"),
					LTS:         "false",
					EOL:         eolapi.FlexibleString("2022-01-01"),
					Latest:      eolapi.FlexibleString("21.04"),
					ReleaseDate: eolapi.FlexibleDate("2021-04-22"),
				},
			},
			config: config.Config{
				Port: 2112,
				Products: []map[string]config.Product{
					{
						"spring-framework": {Host: "localhost", Version: "3.3"},
						"nodejs":           {Host: "localhost", Version: "16"},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(tt.mockData)
			}))
			defer server.Close()

			client := &eolapi.Client{
				HTTPClient: server.Client(),
				BaseURL:    server.URL,
			}

			err := RegisterMetrics(tt.config, client)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStartHTTPServer(t *testing.T) {
	// Create a test server with a custom mux
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Get the port from the test server URL
	port := 2112

	// Start the server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- StartHTTPServer(port)
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Try to connect to the server
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/metrics", port))
	if err != nil {
		t.Errorf("Failed to connect to server: %v", err)
	} else {
		resp.Body.Close()
	}
}

func TestStartExporter(t *testing.T) {
	// Create a mock server for the API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]eolapi.EOLData{
			{
				Cycle:       eolapi.FlexibleString("1.0"),
				LTS:         "false",
				EOL:         eolapi.FlexibleString("2024-01-01"),
				Latest:      eolapi.FlexibleString("1.0"),
				ReleaseDate: eolapi.FlexibleDate("2023-01-01"),
			},
		})
	}))
	defer server.Close()

	cfg := config.Config{
		Port: 2113, // Use a different port
		Products: []map[string]config.Product{
			{
				"test-product": {Host: "localhost", Version: "1.0"},
			},
		},
	}

	// Create a client with the mock server
	client := &eolapi.Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
	}

	// Register metrics first
	if err := RegisterMetrics(cfg, client); err != nil {
		t.Fatalf("Failed to register metrics: %v", err)
	}

	// Create a test server with a custom mux
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Start the server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- StartHTTPServer(cfg.Port)
	}()

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	// Try to connect to the server
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/metrics", cfg.Port))
	if err != nil {
		t.Errorf("Failed to connect to server: %v", err)
	} else {
		resp.Body.Close()
	}
}
