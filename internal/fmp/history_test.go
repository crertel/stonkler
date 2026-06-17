package fmp

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestStockHistoryUsesEODEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/historical-price-eod/full" {
			t.Fatalf("path = %q, want /historical-price-eod/full", got)
		}
		if got := r.URL.Query().Get("symbol"); got != "AAPL" {
			t.Fatalf("symbol = %q, want AAPL", got)
		}
		if got := r.URL.Query().Get("from"); got != "2026-06-10" {
			t.Fatalf("from = %q, want 2026-06-10", got)
		}
		if got := r.URL.Query().Get("to"); got != "2026-06-12" {
			t.Fatalf("to = %q, want 2026-06-12", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"AAPL","date":"2026-06-12","open":296.03,"high":297.14,"low":289.62,"close":291.13,"volume":38784789,"change":-4.9,"changePercent":-1.66,"vwap":293.48}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	prices, err := client.StockHistory(context.Background(), StockHistoryRequest{
		Symbol: "aapl",
		From:   "2026-06-10",
		To:     "2026-06-12",
	})
	if err != nil {
		t.Fatalf("StockHistory() error = %v", err)
	}
	if len(prices) != 1 {
		t.Fatalf("len(prices) = %d, want 1", len(prices))
	}
	if prices[0].Date != "2026-06-12" {
		t.Fatalf("prices[0].Date = %q, want 2026-06-12", prices[0].Date)
	}
}
