package eolapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
	if client.HTTPClient == nil {
		t.Fatal("Expected non-nil HTTP client")
	}
	if client.BaseURL != "https://endoflife.date/api" {
		t.Errorf("Expected default base URL, got %s", client.BaseURL)
	}
}

func TestFlexibleStringUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"string value", `"test"`, "test", false},
		{"boolean true", `true`, "true", false},
		{"boolean false", `false`, "false", false},
		{"number", `42`, "42", false},
		{"null", `null`, "", false},
		{"invalid json", `{`, "", true},
		{"invalid type", `{"key": "value"}`, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fs FlexibleString
			err := json.Unmarshal([]byte(tt.input), &fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlexibleString.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(fs) != tt.expected {
				t.Errorf("FlexibleString.UnmarshalJSON() = %v, want %v", string(fs), tt.expected)
			}
		})
	}
}

func TestFlexibleDateUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"date string", `"2023-01-01"`, "2023-01-01", false},
		{"boolean true", `true`, "true", false},
		{"boolean false", `false`, "false", false},
		{"null", `null`, "", false},
		{"invalid json", `{`, "", true},
		{"invalid type", `{"key": "value"}`, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fd FlexibleDate
			err := json.Unmarshal([]byte(tt.input), &fd)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlexibleDate.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(fd) != tt.expected {
				t.Errorf("FlexibleDate.UnmarshalJSON() = %v, want %v", string(fd), tt.expected)
			}
		})
	}
}

func TestFetchEOLData(t *testing.T) {
	tests := []struct {
		name     string
		handler  http.HandlerFunc
		wantErr  bool
		validate func(*testing.T, []EOLData)
	}{
		{
			name: "successful request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if !strings.HasSuffix(r.URL.Path, "spring-framework.json") {
					t.Errorf("Expected path to end with spring-framework.json, got %s", r.URL.Path)
				}
				json.NewEncoder(w).Encode([]EOLData{
					{
						Cycle:        FlexibleString("21.04"),
						LTS:          "false",
						ReleaseDate:  FlexibleDate("2021-04-22"),
						Support:      FlexibleDate("2022-01-01"),
						EOL:          FlexibleString("2022-01-01"),
						Latest:       FlexibleString("21.04"),
						Link:         FlexibleString("https://wiki.ubuntu.com/HirsuteHippo/ReleaseNotes/"),
						Discontinued: FlexibleDate("2022-01-01"),
					},
				})
			},
			wantErr: false,
			validate: func(t *testing.T, data []EOLData) {
				if len(data) != 1 {
					t.Errorf("Expected 1 item in data, got %d", len(data))
					return
				}
				if string(data[0].Cycle) != "21.04" {
					t.Errorf("Expected cycle 21.04, got %s", data[0].Cycle)
				}
			},
		},
		{
			name: "non-200 response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr: true,
		},
		{
			name: "invalid JSON response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("invalid json"))
			},
			wantErr: true,
		},
		{
			name: "timeout response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(50 * time.Millisecond)
				w.WriteHeader(http.StatusRequestTimeout)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := &Client{
				HTTPClient: &http.Client{
					Timeout: 25 * time.Millisecond,
				},
				BaseURL: server.URL,
			}

			data, err := client.FetchEOLData("spring-framework")
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchEOLData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, data)
			}
		})
	}
}

func TestFetchEOLDataTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Simulate a delay
		w.WriteHeader(http.StatusRequestTimeout)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{
			Timeout: 1 * time.Second, // Set a timeout
		},
	}

	_, err := client.FetchEOLData("spring-framework")
	if err == nil {
		t.Fatal("Expected timeout error, got none")
	}
}
