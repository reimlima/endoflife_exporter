package eolapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FlexibleString is a type that can unmarshal both string and boolean values
type FlexibleString string

// UnmarshalJSON implements the json.Unmarshaler interface
func (fs *FlexibleString) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		*fs = FlexibleString(v)
	case bool:
		*fs = FlexibleString(fmt.Sprintf("%t", v))
	case float64:
		*fs = FlexibleString(fmt.Sprintf("%g", v))
	case nil:
		*fs = ""
	default:
		return fmt.Errorf("unexpected type for FlexibleString: %T", v)
	}
	return nil
}

// FlexibleDate can handle both date strings and booleans
type FlexibleDate string

func (fd *FlexibleDate) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		*fd = FlexibleDate(v)
	case bool:
		*fd = FlexibleDate(fmt.Sprintf("%t", v))
	case nil:
		*fd = ""
	default:
		return fmt.Errorf("unexpected type for FlexibleDate: %T", v)
	}
	return nil
}

type EOLData struct {
	Cycle        FlexibleString `json:"cycle"`        // can be number or string
	ReleaseDate  FlexibleDate   `json:"releaseDate"`  // string<date>
	EOL          FlexibleString `json:"eol"`          // string or boolean
	Latest       FlexibleString `json:"latest"`       // string
	Link         FlexibleString `json:"link"`         // string or null
	LTS          interface{}    `json:"lts"`          // boolean or string
	Support      FlexibleDate   `json:"support"`      // string<date> or boolean
	Discontinued FlexibleDate   `json:"discontinued"` // string<date> or boolean
}

type Client struct {
	HTTPClient *http.Client
	BaseURL    string // Base URL for the API
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{},
		BaseURL:    "https://endoflife.date/api",
	}
}

func (c *Client) FetchEOLData(product string) ([]EOLData, error) {
	url := fmt.Sprintf("%s/%s.json", c.BaseURL, product)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var data []EOLData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return data, nil
}
