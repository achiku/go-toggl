package toggl

import (
	"context"
	"testing"
)

// TestNewConfig creates test config
func TestNewConfig(host string) *Config {
	return &Config{
		APIKey: "testapikey",
		Debug:  true,
		Host:   host,
	}
}

// TestNewClient creates test client
func TestNewClient(t *testing.T, host string) (*Client, context.Context) {
	cfg := &Config{
		APIKey: "testapikey",
		Debug:  true,
		Host:   host,
	}
	ctx := context.Background()
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return client, ctx
}
