package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetAnalystClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetAnalystClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetAnalystClient) StockRatingSnapshot(context.Context, string) ([]fmp.StockRatingSnapshot, error) {
	return nil, errors.New("not used")
}

func TestGetAnalystSymbolResolvesName(t *testing.T) {
	client := &fakeGetAnalystClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbol, err := getAnalystSymbol(context.Background(), client, "apple")
	if err != nil {
		t.Fatalf("getAnalystSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getAnalystSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "apple" {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}

func TestGetAnalystSymbolLeavesTickerStyleQueryAlone(t *testing.T) {
	client := &fakeGetAnalystClient{
		searchResults: []fmp.SearchResult{{Symbol: "MSFT", Name: "Microsoft Corporation"}},
	}

	symbol, err := getAnalystSymbol(context.Background(), client, "AAPL")
	if err != nil {
		t.Fatalf("getAnalystSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getAnalystSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
