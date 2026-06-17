package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetFundInfoClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetFundInfoClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetFundInfoClient) FundProfile(context.Context, string) (fmp.StockProfile, error) {
	return fmp.StockProfile{}, errors.New("not used")
}

func TestGetFundInfoSymbolResolvesName(t *testing.T) {
	client := &fakeGetFundInfoClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "VTI.MX", Name: "Vanguard Total Stock Market ETF", Currency: "MXN"},
			{Symbol: "VTI", Name: "Vanguard Total Stock Market ETF", Currency: "USD"},
		},
	}

	symbol, err := getFundInfoSymbol(context.Background(), client, "vanguard total stock")
	if err != nil {
		t.Fatalf("getFundInfoSymbol() error = %v", err)
	}
	if symbol != "VTI" {
		t.Fatalf("getFundInfoSymbol() = %q, want VTI", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "vanguard total stock" {
		t.Fatalf("searchQueries = %#v, want vanguard total stock", client.searchQueries)
	}
}

func TestGetFundInfoSymbolUppercasesShortTickerStyleQuery(t *testing.T) {
	client := &fakeGetFundInfoClient{
		searchResults: []fmp.SearchResult{{Symbol: "SPY", Name: "SPDR S&P 500 ETF Trust"}},
	}

	symbol, err := getFundInfoSymbol(context.Background(), client, "spy")
	if err != nil {
		t.Fatalf("getFundInfoSymbol() error = %v", err)
	}
	if symbol != "SPY" {
		t.Fatalf("getFundInfoSymbol() = %q, want SPY", symbol)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}
