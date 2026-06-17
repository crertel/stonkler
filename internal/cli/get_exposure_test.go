package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetExposureClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetExposureClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetExposureClient) ETFAssetExposure(context.Context, string) ([]fmp.ETFAssetExposure, error) {
	return nil, errors.New("not used")
}

func TestGetExposureSymbolResolvesName(t *testing.T) {
	client := &fakeGetExposureClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbol, err := getExposureSymbol(context.Background(), client, "apple")
	if err != nil {
		t.Fatalf("getExposureSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getExposureSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "apple" {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}

func TestGetExposureSymbolLeavesTickerStyleQueryAlone(t *testing.T) {
	client := &fakeGetExposureClient{
		searchResults: []fmp.SearchResult{{Symbol: "MSFT", Name: "Microsoft Corporation"}},
	}

	symbol, err := getExposureSymbol(context.Background(), client, "AAPL")
	if err != nil {
		t.Fatalf("getExposureSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getExposureSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
