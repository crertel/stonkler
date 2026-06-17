package fmp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL   = "https://financialmodelingprep.com/stable"
	defaultV3BaseURL = "https://financialmodelingprep.com/api/v3"
	defaultV4BaseURL = "https://financialmodelingprep.com/api/v4"
)

// Client calls the Financial Modeling Prep stable API.
type Client struct {
	apiKey     string
	baseURL    string
	v3BaseURL  string
	v4BaseURL  string
	httpClient *http.Client
}

// NewClient returns a Financial Modeling Prep API client.
func NewClient(apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		v3BaseURL:  defaultV3BaseURL,
		v4BaseURL:  defaultV4BaseURL,
		httpClient: httpClient,
	}
}

func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	endpoint, err := url.Parse(c.baseURL + path)
	if err != nil {
		return err
	}
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("fmp returned HTTP %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		return err
	}
	return nil
}

func (c *Client) getV3(ctx context.Context, path string, query url.Values, out any) error {
	return c.getFromBase(ctx, c.v3BaseURL, path, query, out)
}

func (c *Client) getV4(ctx context.Context, path string, query url.Values, out any) error {
	return c.getFromBase(ctx, c.v4BaseURL, path, query, out)
}

func (c *Client) getFromBase(ctx context.Context, baseURL, path string, query url.Values, out any) error {
	endpoint, err := url.Parse(baseURL + path)
	if err != nil {
		return err
	}
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("fmp returned HTTP %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		return err
	}
	return nil
}
