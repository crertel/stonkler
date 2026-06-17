package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetHoldingsClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetHoldingsClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetHoldingsClient) ETFHoldings(context.Context, string) ([]fmp.ETFHolding, error) {
	return nil, errors.New("not used")
}

func TestGetHoldingsSymbolResolvesName(t *testing.T) {
	client := &fakeGetHoldingsClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "VTI.MX", Name: "Vanguard Total Stock Market ETF", Currency: "MXN"},
			{Symbol: "VTI", Name: "Vanguard Total Stock Market ETF", Currency: "USD"},
		},
	}

	symbol, err := getHoldingsSymbol(context.Background(), client, "vanguard total stock")
	if err != nil {
		t.Fatalf("getHoldingsSymbol() error = %v", err)
	}
	if symbol != "VTI" {
		t.Fatalf("getHoldingsSymbol() = %q, want VTI", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "vanguard total stock" {
		t.Fatalf("searchQueries = %#v, want vanguard total stock", client.searchQueries)
	}
}

func TestGetHoldingsSymbolUppercasesShortTickerStyleQuery(t *testing.T) {
	client := &fakeGetHoldingsClient{
		searchResults: []fmp.SearchResult{{Symbol: "SPY", Name: "SPDR S&P 500 ETF Trust"}},
	}

	symbol, err := getHoldingsSymbol(context.Background(), client, "spy")
	if err != nil {
		t.Fatalf("getHoldingsSymbol() error = %v", err)
	}
	if symbol != "SPY" {
		t.Fatalf("getHoldingsSymbol() = %q, want SPY", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
