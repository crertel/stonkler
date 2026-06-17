package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetHistoryClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetHistoryClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetHistoryClient) StockHistory(context.Context, fmp.StockHistoryRequest) ([]fmp.StockPrice, error) {
	return nil, errors.New("not used")
}

func TestGetHistorySymbolResolvesName(t *testing.T) {
	client := &fakeGetHistoryClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbol, err := getHistorySymbol(context.Background(), client, "apple")
	if err != nil {
		t.Fatalf("getHistorySymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getHistorySymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "apple" {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}

func TestGetHistorySymbolLeavesTickerStyleQueryAlone(t *testing.T) {
	client := &fakeGetHistoryClient{
		searchResults: []fmp.SearchResult{{Symbol: "MSFT", Name: "Microsoft Corporation"}},
	}

	symbol, err := getHistorySymbol(context.Background(), client, "AAPL")
	if err != nil {
		t.Fatalf("getHistorySymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getHistorySymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
