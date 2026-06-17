package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetInsidersClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetInsidersClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetInsidersClient) InsiderTrades(context.Context, string, int) ([]fmp.InsiderTrade, error) {
	return nil, errors.New("not used")
}

func TestGetInsidersSymbolResolvesName(t *testing.T) {
	client := &fakeGetInsidersClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbol, err := getInsidersSymbol(context.Background(), client, "apple")
	if err != nil {
		t.Fatalf("getInsidersSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getInsidersSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "apple" {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}

func TestGetInsidersSymbolLeavesTickerStyleQueryAlone(t *testing.T) {
	client := &fakeGetInsidersClient{
		searchResults: []fmp.SearchResult{{Symbol: "MSFT", Name: "Microsoft Corporation"}},
	}

	symbol, err := getInsidersSymbol(context.Background(), client, "AAPL")
	if err != nil {
		t.Fatalf("getInsidersSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getInsidersSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
