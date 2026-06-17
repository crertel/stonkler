package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetSectorWeightingsClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetSectorWeightingsClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetSectorWeightingsClient) ETFSectorWeightings(context.Context, string) ([]fmp.ETFSectorWeighting, error) {
	return nil, errors.New("not used")
}

func TestGetSectorWeightingsSymbolResolvesName(t *testing.T) {
	client := &fakeGetSectorWeightingsClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "VTI.MX", Name: "Vanguard Total Stock Market ETF", Currency: "MXN"},
			{Symbol: "VTI", Name: "Vanguard Total Stock Market ETF", Currency: "USD"},
		},
	}

	symbol, err := getSectorWeightingsSymbol(context.Background(), client, "vanguard total stock")
	if err != nil {
		t.Fatalf("getSectorWeightingsSymbol() error = %v", err)
	}
	if symbol != "VTI" {
		t.Fatalf("getSectorWeightingsSymbol() = %q, want VTI", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "vanguard total stock" {
		t.Fatalf("searchQueries = %#v, want vanguard total stock", client.searchQueries)
	}
}

func TestGetSectorWeightingsSymbolUppercasesShortTickerStyleQuery(t *testing.T) {
	client := &fakeGetSectorWeightingsClient{
		searchResults: []fmp.SearchResult{{Symbol: "SPY", Name: "SPDR S&P 500 ETF Trust"}},
	}

	symbol, err := getSectorWeightingsSymbol(context.Background(), client, "spy")
	if err != nil {
		t.Fatalf("getSectorWeightingsSymbol() error = %v", err)
	}
	if symbol != "SPY" {
		t.Fatalf("getSectorWeightingsSymbol() = %q, want SPY", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
