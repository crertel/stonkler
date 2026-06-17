package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetCountryWeightingsClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetCountryWeightingsClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetCountryWeightingsClient) ETFCountryWeightings(context.Context, string) ([]fmp.ETFCountryWeighting, error) {
	return nil, errors.New("not used")
}

func TestGetCountryWeightingsSymbolResolvesName(t *testing.T) {
	client := &fakeGetCountryWeightingsClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "VXUS.MX", Name: "Vanguard Total International Stock ETF", Currency: "MXN"},
			{Symbol: "VXUS", Name: "Vanguard Total International Stock ETF", Currency: "USD"},
		},
	}

	symbol, err := getCountryWeightingsSymbol(context.Background(), client, "vanguard total international")
	if err != nil {
		t.Fatalf("getCountryWeightingsSymbol() error = %v", err)
	}
	if symbol != "VXUS" {
		t.Fatalf("getCountryWeightingsSymbol() = %q, want VXUS", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "vanguard total international" {
		t.Fatalf("searchQueries = %#v, want vanguard total international", client.searchQueries)
	}
}

func TestGetCountryWeightingsSymbolUppercasesShortTickerStyleQuery(t *testing.T) {
	client := &fakeGetCountryWeightingsClient{
		searchResults: []fmp.SearchResult{{Symbol: "VXUS", Name: "Vanguard Total International Stock ETF"}},
	}

	symbol, err := getCountryWeightingsSymbol(context.Background(), client, "vxus")
	if err != nil {
		t.Fatalf("getCountryWeightingsSymbol() error = %v", err)
	}
	if symbol != "VXUS" {
		t.Fatalf("getCountryWeightingsSymbol() = %q, want VXUS", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
