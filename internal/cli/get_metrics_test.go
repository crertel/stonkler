package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetMetricsClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetMetricsClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetMetricsClient) StockKeyMetricsTTM(context.Context, string) ([]fmp.StockMetricRow, error) {
	return nil, errors.New("not used")
}

func TestGetMetricsSymbolResolvesName(t *testing.T) {
	client := &fakeGetMetricsClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbol, err := getMetricsSymbol(context.Background(), client, "apple")
	if err != nil {
		t.Fatalf("getMetricsSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getMetricsSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "apple" {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}

func TestGetMetricsSymbolLeavesTickerStyleQueryAlone(t *testing.T) {
	client := &fakeGetMetricsClient{
		searchResults: []fmp.SearchResult{{Symbol: "MSFT", Name: "Microsoft Corporation"}},
	}

	symbol, err := getMetricsSymbol(context.Background(), client, "AAPL")
	if err != nil {
		t.Fatalf("getMetricsSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getMetricsSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
