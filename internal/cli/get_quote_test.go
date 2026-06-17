package cli

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetQuoteClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetQuoteClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetQuoteClient) StockQuotes(context.Context, []string) ([]fmp.StockQuote, error) {
	return nil, errors.New("not used")
}

func TestGetQuoteSymbolsResolvesLowercaseName(t *testing.T) {
	client := &fakeGetQuoteClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbols, err := getQuoteSymbols(context.Background(), client, []string{"apple"})
	if err != nil {
		t.Fatalf("getQuoteSymbols() error = %v", err)
	}
	if !reflect.DeepEqual(symbols, []string{"AAPL"}) {
		t.Fatalf("getQuoteSymbols() = %#v, want AAPL", symbols)
	}
	if !reflect.DeepEqual(client.searchQueries, []string{"apple"}) {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}

func TestGetQuoteSymbolsLeavesTickerStyleArgsAlone(t *testing.T) {
	client := &fakeGetQuoteClient{
		searchResults: []fmp.SearchResult{{Symbol: "MSFT", Name: "Microsoft Corporation"}},
	}

	symbols, err := getQuoteSymbols(context.Background(), client, []string{"AAPL", "MSFT"})
	if err != nil {
		t.Fatalf("getQuoteSymbols() error = %v", err)
	}
	if !reflect.DeepEqual(symbols, []string{"AAPL", "MSFT"}) {
		t.Fatalf("getQuoteSymbols() = %#v, want original symbols", symbols)
	}
	if len(client.searchQueries) != 0 {
		t.Fatalf("searchQueries = %#v, want no search", client.searchQueries)
	}
}

func TestGetQuoteSymbolsReturnsMissingSearchResult(t *testing.T) {
	client := &fakeGetQuoteClient{}

	_, err := getQuoteSymbols(context.Background(), client, []string{"not a company"})
	if err == nil {
		t.Fatal("getQuoteSymbols() error = nil, want error")
	}
}
