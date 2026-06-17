package fmp

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestETFHoldingsUsesV3Endpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/etf-holder/SPY" {
			t.Fatalf("path = %q, want /etf-holder/SPY", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"asset":"AAPL","name":"APPLE INC","isin":"US0378331005","cusip":"037833100","sharesNumber":176183317,"weightPercentage":6.81,"marketValue":53568630492,"updated":"2026-06-17 04:06:08"}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.v3BaseURL = "https://example.test"

	holdings, err := client.ETFHoldings(context.Background(), "spy")
	if err != nil {
		t.Fatalf("ETFHoldings() error = %v", err)
	}
	if len(holdings) != 1 {
		t.Fatalf("len(holdings) = %d, want 1", len(holdings))
	}
	if holdings[0].Asset != "AAPL" {
		t.Fatalf("holdings[0].Asset = %q, want AAPL", holdings[0].Asset)
	}
}
