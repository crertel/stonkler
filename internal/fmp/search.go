package fmp

import (
	"context"
	"net/url"
	"strconv"
)

// SearchResult is a normalized security search result from FMP.
type SearchResult struct {
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Currency          string `json:"currency,omitempty"`
	StockExchange     string `json:"stockExchange,omitempty"`
	ExchangeShortName string `json:"exchangeShortName,omitempty"`
}

// ScreenerOptions are supported company screener filters.
type ScreenerOptions struct {
	Sector       string
	Country      string
	MarketCapMin float64
	Limit        int
}

// ScreenerResult is one row returned by FMP's company screener.
type ScreenerResult struct {
	Symbol             string  `json:"symbol"`
	CompanyName        string  `json:"companyName"`
	MarketCap          float64 `json:"marketCap"`
	Sector             string  `json:"sector,omitempty"`
	Industry           string  `json:"industry,omitempty"`
	Beta               float64 `json:"beta,omitempty"`
	Price              float64 `json:"price"`
	LastAnnualDividend float64 `json:"lastAnnualDividend,omitempty"`
	Volume             float64 `json:"volume,omitempty"`
	Exchange           string  `json:"exchange,omitempty"`
	ExchangeShortName  string  `json:"exchangeShortName,omitempty"`
	Country            string  `json:"country,omitempty"`
	IsETF              bool    `json:"isEtf"`
	IsFund             bool    `json:"isFund"`
	IsActivelyTrading  bool    `json:"isActivelyTrading"`
}

// SearchName searches FMP by company or asset name.
func (c *Client) SearchName(ctx context.Context, query string) ([]SearchResult, error) {
	return c.search(ctx, "/search-name", "query", query)
}

// SearchSymbol searches FMP by symbol.
func (c *Client) SearchSymbol(ctx context.Context, query string) ([]SearchResult, error) {
	return c.search(ctx, "/search-symbol", "query", query)
}

// SearchCIK searches FMP by Central Index Key.
func (c *Client) SearchCIK(ctx context.Context, cik string) ([]SearchResult, error) {
	return c.search(ctx, "/search-cik", "cik", cik)
}

// SearchISIN searches FMP by International Securities Identification Number.
func (c *Client) SearchISIN(ctx context.Context, isin string) ([]SearchResult, error) {
	return c.search(ctx, "/search-isin", "isin", isin)
}

// CompanyScreener screens companies using supported FMP filters.
func (c *Client) CompanyScreener(ctx context.Context, options ScreenerOptions) ([]ScreenerResult, error) {
	query := url.Values{}
	if options.Sector != "" {
		query.Set("sector", options.Sector)
	}
	if options.Country != "" {
		query.Set("country", options.Country)
	}
	if options.MarketCapMin > 0 {
		query.Set("marketCapMoreThan", strconv.FormatFloat(options.MarketCapMin, 'f', -1, 64))
	}
	if options.Limit > 0 {
		query.Set("limit", strconv.Itoa(options.Limit))
	}

	var results []ScreenerResult
	if err := c.get(ctx, "/company-screener", query, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Client) search(ctx context.Context, path string, key string, value string) ([]SearchResult, error) {
	var results []SearchResult
	query := url.Values{}
	query.Set(key, value)
	if err := c.get(ctx, path, query, &results); err != nil {
		return nil, err
	}
	return results, nil
}
