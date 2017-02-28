package toggl

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Client toggl client
type Client struct {
	config *Config
}

// Config toggl client config
type Config struct {
	HTTPClient *http.Client
	Host       string
	APIKey     string
	Debug      bool
	Logger     *log.Logger
}

// NewClient creates toggl client
func NewClient(cfg *Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("Config.APIKey is not set")
	}
	if cfg.Host == "" {
		return nil, errors.New("Config.Host is not set")
	}
	if cfg.Logger == nil {
		cfg.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	if cfg.HTTPClient == nil {
		hc := &http.Client{
			Timeout: time.Duration(10 * time.Second),
			Transport: &http.Transport{
				MaxIdleConns: 10,
			},
		}
		cfg.HTTPClient = hc
		return &Client{config: cfg}, nil
	}
	return &Client{config: cfg}, nil
}

func (c *Client) call(ctx context.Context, method, pathStr string, req, res interface{}) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "failed to marshal request")
	}
	if c.config.Debug {
		c.config.Logger.Printf("request [%s: %s] %s", method, pathStr, payload)
	}

	endpoint := fmt.Sprintf("%s%s", c.config.Host, pathStr)
	request, err := http.NewRequest(method, endpoint, strings.NewReader(string(payload)))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	request.SetBasicAuth(c.config.APIKey, "api_token")
	request.WithContext(ctx)
	request.Header.Add("Content-type", "application/json")

	response, err := c.config.HTTPClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("status code: %d, body: %s", response.StatusCode, response.Body)
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(res); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if c.config.Debug {
		c.config.Logger.Printf("response [%s: %s] %s", method, pathStr, res)
	}
	return nil
}

// Workspace workspace
type Workspace struct {
	ID                          int       `json:"id"`
	Name                        string    `json:"name"`
	Profile                     int       `json:"profile"`
	Premium                     bool      `json:"premium"`
	Admin                       bool      `json:"admin"`
	DefaultHourlyRate           int       `json:"default_hourly_rate"`
	DefaultCurrency             string    `json:"default_currency"`
	OnlyAdminsMayCreateProjects bool      `json:"only_admins_may_create_projects"`
	OnlyAdminsSeeBillableRates  bool      `json:"only_admins_see_billable_rates"`
	OnlyAdminsSeeTeamDashboard  bool      `json:"only_admins_see_team_dashboard"`
	ProjectsBillableByDefault   bool      `json:"projects_billable_by_default"`
	Rounding                    int       `json:"rounding"`
	RoundingMinutes             int       `json:"rounding_minutes"`
	APIToken                    string    `json:"api_token"`
	At                          time.Time `json:"at"`
	IcalEnabled                 bool      `json:"ical_enabled"`
}

// GetWorkspaces gets workspaces
func (c *Client) GetWorkspaces(ctx context.Context) ([]Workspace, error) {
	pt := "/api/v8/workspaces"
	m := "GET"
	var res []Workspace
	if err := c.call(ctx, m, pt, nil, &res); err != nil {
		return nil, errors.Wrapf(err, "failed call [%s] %s/%s", m, c.config.Host, pt)
	}
	return res, nil
}

// GetWorkspaceByID gets workspaces
func (c *Client) GetWorkspaceByID(ctx context.Context, wid int) (*Workspace, error) {
	pt := fmt.Sprintf("/api/v8/workspaces/%d", wid)
	m := "GET"
	var res Workspace
	if err := c.call(ctx, m, pt, nil, &res); err != nil {
		return nil, errors.Wrapf(err, "failed call [%s] %s/%s", m, c.config.Host, pt)
	}
	return &res, nil
}

// Dashboard dashboard
type Dashboard struct {
	MostActiveUser []struct {
		UserID   int `json:"user_id"`
		Duration int `json:"duration"`
	} `json:"most_active_user"`
	Activity []struct {
		UserID      int         `json:"user_id"`
		ProjectID   int         `json:"project_id"`
		Duration    int         `json:"duration"`
		Description string      `json:"description"`
		Stop        interface{} `json:"stop"`
		Tid         interface{} `json:"tid"`
	} `json:"activity"`
}

// GetDashboardByWorkspaceID gets dashboard
func (c *Client) GetDashboardByWorkspaceID(ctx context.Context, wid int) (*Dashboard, error) {
	pt := fmt.Sprintf("/api/v8/dashboard/%d", wid)
	m := "GET"
	var res Dashboard
	if err := c.call(ctx, m, pt, nil, &res); err != nil {
		return nil, errors.Wrapf(err, "failed call [%s] %s/%s", m, c.config.Host, pt)
	}
	return &res, nil
}

// DetailedReport detailed report
type DetailedReport struct {
	TotalGrand      int `json:"total_grand"`
	TotalBillable   int `json:"total_billable"`
	TotalCount      int `json:"total_count"`
	PerPage         int `json:"per_page"`
	TotalCurrencies []struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	} `json:"total_currencies"`
	Data []struct {
		ID          int         `json:"id"`
		Pid         int         `json:"pid"`
		Tid         interface{} `json:"tid"`
		UID         int         `json:"uid"`
		Description string      `json:"description"`
		Start       string      `json:"start"`
		End         string      `json:"end"`
		Updated     string      `json:"updated"`
		Dur         int         `json:"dur"`
		User        string      `json:"user"`
		UseStop     bool        `json:"use_stop"`
		Client      string      `json:"client"`
		Project     string      `json:"project"`
		Task        interface{} `json:"task"`
		Billable    float64     `json:"billable"`
		IsBillable  bool        `json:"is_billable"`
		Cur         string      `json:"cur"`
		Tags        []string    `json:"tags"`
	} `json:"data"`
}
