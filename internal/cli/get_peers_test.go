package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetPeersClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetPeersClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetPeersClient) StockPeers(context.Context, string) ([]fmp.StockPeer, error) {
	return nil, errors.New("not used")
}

func TestGetPeersSymbolResolvesName(t *testing.T) {
	client := &fakeGetPeersClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbol, err := getPeersSymbol(context.Background(), client, "apple")
	if err != nil {
		t.Fatalf("getPeersSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getPeersSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "apple" {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}

func TestGetPeersSymbolLeavesTickerStyleQueryAlone(t *testing.T) {
	client := &fakeGetPeersClient{
		searchResults: []fmp.SearchResult{{Symbol: "MSFT", Name: "Microsoft Corporation"}},
	}

	symbol, err := getPeersSymbol(context.Background(), client, "AAPL")
	if err != nil {
		t.Fatalf("getPeersSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getPeersSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
