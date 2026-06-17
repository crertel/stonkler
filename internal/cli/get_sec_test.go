package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetSECClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetSECClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetSECClient) SECFilings(context.Context, string, int) ([]fmp.SECFiling, error) {
	return nil, errors.New("not used")
}

func TestGetSECSymbolResolvesName(t *testing.T) {
	client := &fakeGetSECClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbol, err := getSECSymbol(context.Background(), client, "apple")
	if err != nil {
		t.Fatalf("getSECSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getSECSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "apple" {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}

func TestGetSECSymbolLeavesTickerStyleQueryAlone(t *testing.T) {
	client := &fakeGetSECClient{
		searchResults: []fmp.SearchResult{{Symbol: "MSFT", Name: "Microsoft Corporation"}},
	}

	symbol, err := getSECSymbol(context.Background(), client, "AAPL")
	if err != nil {
		t.Fatalf("getSECSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getSECSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
