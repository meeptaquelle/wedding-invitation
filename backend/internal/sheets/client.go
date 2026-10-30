package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Client struct {
	svc     *sheets.Service
	sheetID string
}

var (
	clientOnce sync.Once
	clientRef  *Client
	clientErr  error
)

// NewClient creates a Sheets client. Credentials should be provided as raw strings.
func NewClient(ctx context.Context, email, privateKey, sheetID string) (*Client, error) {
	credsMap := map[string]string{
		"type":         "service_account",
		"client_email": email,
		"private_key":  privateKey,
		"token_uri":    "https://oauth2.googleapis.com/token",
	}
	credsJSON, err := json.Marshal(credsMap)
	if err != nil {
		return nil, fmt.Errorf("marshal credentials: %w", err)
	}

	cfg, err := google.JWTConfigFromJSON(credsJSON, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("jwt config: %w", err)
	}

	svc, err := sheets.NewService(ctx, option.WithHTTPClient(cfg.Client(ctx)))
	if err != nil {
		return nil, fmt.Errorf("new sheets service: %w", err)
	}

	return &Client{svc: svc, sheetID: sheetID}, nil
}

func (c *Client) Read(ctx context.Context, rng string) ([][]interface{}, error) {
	resp, err := c.svc.Spreadsheets.Values.Get(c.sheetID, rng).Do()
	if err != nil {
		return nil, fmt.Errorf("sheets read %s: %w", rng, err)
	}
	return resp.Values, nil
}

func (c *Client) Write(ctx context.Context, rng string, values [][]interface{}) error {
	vr := &sheets.ValueRange{Values: values}
	_, err := c.svc.Spreadsheets.Values.Update(c.sheetID, rng, vr).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		return fmt.Errorf("sheets write %s: %w", rng, err)
	}
	return nil
}

func (c *Client) Append(ctx context.Context, rng string, values [][]interface{}) error {
	vr := &sheets.ValueRange{Values: values}
	_, err := c.svc.Spreadsheets.Values.Append(c.sheetID, rng, vr).
		ValueInputOption("RAW").
		InsertDataOption("INSERT_ROWS").
		Do()
	if err != nil {
		return fmt.Errorf("sheets append %s: %w", rng, err)
	}
	return nil
}

// EnsureRowCount is a helper used to confirm writes matched expected row counts.
// Currently a no-op placeholder; reserved for future use.
func (c *Client) EnsureRowCount(n int) {}

// WarmCache is a no-op that documents where a warm cache could be added.
func WarmCache(ctx context.Context, c *Client) error { return nil }

var _ = time.Now // silence unused import when helpers change
