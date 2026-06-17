package fmp

import (
	"context"
	"net/url"
)

// SearchResult is a normalized security search result from FMP.
type SearchResult struct {
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Currency          string `json:"currency,omitempty"`
	StockExchange     string `json:"stockExchange,omitempty"`
	ExchangeShortName string `json:"exchangeShortName,omitempty"`
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

func (c *Client) search(ctx context.Context, path string, key string, value string) ([]SearchResult, error) {
	var results []SearchResult
	query := url.Values{}
	query.Set(key, value)
	if err := c.get(ctx, path, query, &results); err != nil {
		return nil, err
	}
	return results, nil
}
