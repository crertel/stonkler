package cli

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

type fakeGetStatementsClient struct {
	searchResults []fmp.SearchResult
	searchQueries []string
}

func (c *fakeGetStatementsClient) SearchName(_ context.Context, query string) ([]fmp.SearchResult, error) {
	c.searchQueries = append(c.searchQueries, query)
	return c.searchResults, nil
}

func (c *fakeGetStatementsClient) StockStatements(context.Context, fmp.StockStatementsRequest) ([]fmp.FinancialStatement, error) {
	return nil, errors.New("not used")
}

func TestParseGetStatementsOptionsDefaultsToIncome(t *testing.T) {
	var stderr bytes.Buffer

	options, ok := parseGetStatementsOptions([]string{"AAPL", "--csv"}, &stderr)

	if !ok {
		t.Fatalf("parseGetStatementsOptions() ok = false, stderr = %q", stderr.String())
	}
	if options.symbol != "AAPL" {
		t.Fatalf("symbol = %q, want AAPL", options.symbol)
	}
	if options.statement != fmp.StatementIncome {
		t.Fatalf("statement = %q, want income", options.statement)
	}
	if options.period != "annual" {
		t.Fatalf("period = %q, want annual", options.period)
	}
	if options.limit != 5 {
		t.Fatalf("limit = %d, want 5", options.limit)
	}
	if options.format != outputCSV {
		t.Fatalf("format = %q, want csv", options.format)
	}
}

func TestParseGetStatementsOptionsAcceptsStatementType(t *testing.T) {
	var stderr bytes.Buffer

	options, ok := parseGetStatementsOptions([]string{"AAPL", "cash-flow", "--period", "quarter", "--limit", "2"}, &stderr)

	if !ok {
		t.Fatalf("parseGetStatementsOptions() ok = false, stderr = %q", stderr.String())
	}
	if options.statement != fmp.StatementCashFlow {
		t.Fatalf("statement = %q, want cash-flow", options.statement)
	}
	if options.period != "quarter" {
		t.Fatalf("period = %q, want quarter", options.period)
	}
	if options.limit != 2 {
		t.Fatalf("limit = %d, want 2", options.limit)
	}
}

func TestGetStatementsSymbolResolvesName(t *testing.T) {
	client := &fakeGetStatementsClient{
		searchResults: []fmp.SearchResult{
			{Symbol: "AAPL.DE", Name: "Apple Inc.", Currency: "EUR"},
			{Symbol: "AAPL", Name: "Apple Inc.", Currency: "USD"},
		},
	}

	symbol, err := getStatementsSymbol(context.Background(), client, "apple")
	if err != nil {
		t.Fatalf("getStatementsSymbol() error = %v", err)
	}
	if symbol != "AAPL" {
		t.Fatalf("getStatementsSymbol() = %q, want AAPL", symbol)
	}
	if len(client.searchQueries) != 1 || client.searchQueries[0] != "apple" {
		t.Fatalf("searchQueries = %#v, want apple", client.searchQueries)
	}
}
