package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetCompanyClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetCompanyClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetCompanyClient) StockProfile(context.Context, string) (fmp.StockProfile, error) {
	return fmp.StockProfile{}, errors.New("not used")
}

func TestGetCompanySymbolResolvesName(t *testing.T) {
	client := &fakeGetCompanyClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbol, err := getCompanySymbol(context.Background(), client, "apple")
	if err != nil {
		t.Fatalf("getCompanySymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getCompanySymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "apple" {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}

func TestGetCompanySymbolLeavesTickerStyleQueryAlone(t *testing.T) {
	client := &fakeGetCompanyClient{
		searchResults: []fmp.SearchResult{{Symbol: "MSFT", Name: "Microsoft Corporation"}},
	}

	symbol, err := getCompanySymbol(context.Background(), client, "AAPL")
	if err != nil {
		t.Fatalf("getCompanySymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getCompanySymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
