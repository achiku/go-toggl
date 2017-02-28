package toggl

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient(t *testing.T) {
	cases := []struct {
		Config *Config
		Fail   bool
	}{
		{
			Config: &Config{HTTPClient: nil, Host: "http://example.com", APIKey: "apikey"},
			Fail:   false,
		},
		{
			Config: &Config{HTTPClient: nil, Host: "", APIKey: "apikey"},
			Fail:   true,
		},
		{
			Config: &Config{HTTPClient: nil, Host: "http://example.com", APIKey: ""},
			Fail:   true,
		},
	}

	for _, c := range cases {
		client, err := NewClient(c.Config)
		if err != nil {
			if c.Fail {
				continue
			}
			t.Fatal(err)
		}
		if client == nil {
			t.Error("client is nil")
		}
	}
}

func TestGetWorkspaces(t *testing.T) {
	s := httptest.NewServer(TestCreateResponseHandler(TestJSONFileResponse{
		StatusCode: http.StatusOK,
		JSONPath:   "/testdata/response_workspaces.json",
	}, t))

	client, ctx := TestNewClient(t, s.URL)
	res, err := client.GetWorkspaces(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res)
}

func TestGetWorkspaceByID(t *testing.T) {
	s := httptest.NewServer(TestCreateResponseHandler(TestJSONFileResponse{
		StatusCode: http.StatusOK,
		JSONPath:   "/testdata/response_workspace.json",
	}, t))

	client, ctx := TestNewClient(t, s.URL)
	res, err := client.GetWorkspaceByID(ctx, 1111111)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res)
}

func TestGetDashboardByWorkspaceID(t *testing.T) {
	s := httptest.NewServer(TestCreateResponseHandler(TestJSONFileResponse{
		StatusCode: http.StatusOK,
		JSONPath:   "/testdata/response_dashboard.json",
	}, t))

	client, ctx := TestNewClient(t, s.URL)
	res, err := client.GetDashboardByWorkspaceID(ctx, 1111111)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res)
}
