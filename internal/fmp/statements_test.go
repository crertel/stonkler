package fmp

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestStockStatementsUsesIncomeEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/income-statement" {
			t.Fatalf("path = %q, want /income-statement", got)
		}
		if got := r.URL.Query().Get("symbol"); got != "AAPL" {
			t.Fatalf("symbol = %q, want AAPL", got)
		}
		if got := r.URL.Query().Get("period"); got != "annual" {
			t.Fatalf("period = %q, want annual", got)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Fatalf("limit = %q, want 1", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"date":"2025-09-27","symbol":"AAPL","fiscalYear":"2025","period":"FY","revenue":416161000000,"netIncome":112010000000}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	statements, err := client.StockStatements(context.Background(), StockStatementsRequest{
		Symbol:    "aapl",
		Statement: StatementIncome,
		Period:    "annual",
		Limit:     1,
	})
	if err != nil {
		t.Fatalf("StockStatements() error = %v", err)
	}
	if len(statements) != 1 {
		t.Fatalf("len(statements) = %d, want 1", len(statements))
	}
	if statements[0]["symbol"] != "AAPL" {
		t.Fatalf("symbol = %v, want AAPL", statements[0]["symbol"])
	}
}
