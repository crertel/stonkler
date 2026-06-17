package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetRatiosClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetRatiosClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetRatiosClient) StockRatiosTTM(context.Context, string) ([]fmp.StockRatioRow, error) {
	return nil, errors.New("not used")
}

func TestGetRatiosSymbolResolvesName(t *testing.T) {
	client := &fakeGetRatiosClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbol, err := getRatiosSymbol(context.Background(), client, "apple")
	if err != nil {
		t.Fatalf("getRatiosSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getRatiosSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "apple" {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}

func TestGetRatiosSymbolLeavesTickerStyleQueryAlone(t *testing.T) {
	client := &fakeGetRatiosClient{
		searchResults: []fmp.SearchResult{{Symbol: "MSFT", Name: "Microsoft Corporation"}},
	}

	symbol, err := getRatiosSymbol(context.Background(), client, "AAPL")
	if err != nil {
		t.Fatalf("getRatiosSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getRatiosSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
