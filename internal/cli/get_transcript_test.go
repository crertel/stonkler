package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetTranscriptClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetTranscriptClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetTranscriptClient) EarningsCallTranscript(context.Context, string, int, int) ([]fmp.EarningsCallTranscript, error) {
	return nil, errors.New("not used")
}

func (c *fakeGetTranscriptClient) EarningsCallTranscriptDates(context.Context, string) ([]fmp.EarningsCallTranscriptDate, error) {
	return nil, errors.New("not used")
}

func TestGetTranscriptSymbolResolvesName(t *testing.T) {
	client := &fakeGetTranscriptClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbol, err := getTranscriptSymbol(context.Background(), client, "apple")
	if err != nil {
		t.Fatalf("getTranscriptSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getTranscriptSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "apple" {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}

func TestGetTranscriptSymbolLeavesTickerStyleQueryAlone(t *testing.T) {
	client := &fakeGetTranscriptClient{
		searchResults: []fmp.SearchResult{{Symbol: "MSFT", Name: "Microsoft Corporation"}},
	}

	symbol, err := getTranscriptSymbol(context.Background(), client, "AAPL")
	if err != nil {
		t.Fatalf("getTranscriptSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getTranscriptSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
