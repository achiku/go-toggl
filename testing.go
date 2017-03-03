package toggl

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
)

// TestJSONFileResponse test response
type TestJSONFileResponse struct {
	StatusCode int
	JSONPath   string
}

// TestCreateResponseHandler test handler
func TestCreateResponseHandler(res TestJSONFileResponse, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pwd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		p := path.Join(pwd, res.JSONPath)
		t.Logf("reading %s", p)
		src, err := ioutil.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(res.StatusCode)
		w.Write(src)
		return
	}
}

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
