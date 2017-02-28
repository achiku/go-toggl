package toggl

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestGetDetailedReport(t *testing.T) {
	s := httptest.NewServer(TestCreateResponseHandler(TestJSONFileResponse{
		StatusCode: http.StatusOK,
		JSONPath:   "/testdata/response_detailed_report.json",
	}, t))

	client, ctx := TestNewClient(t, s.URL)
	n := time.Now()
	req := &DetailedReportRequest{
		WorkspaceID: 11111,
		Since:       n.AddDate(0, 0, -3),
		Until:       n,
		UserAgent:   "go-toggl",
	}
	res, err := client.GetDetailedReport(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res)
}
