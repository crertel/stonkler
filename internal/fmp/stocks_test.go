package fmp

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestStockQuotesUsesBatchEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/batch-quote" {
			t.Fatalf("path = %q, want /batch-quote", got)
		}
		if got := r.URL.Query().Get("symbols"); got != "AAPL,MSFT" {
			t.Fatalf("symbols = %q, want AAPL,MSFT", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"AAPL","name":"Apple Inc.","price":200.12,"change":1.2,"changePercentage":0.6,"volume":123.45,"marketCap":3000,"timestamp":1710000000}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	quotes, err := client.StockQuotes(context.Background(), []string{"aapl", "msft"})
	if err != nil {
		t.Fatalf("StockQuotes() error = %v", err)
	}
	if len(quotes) != 1 {
		t.Fatalf("len(quotes) = %d, want 1", len(quotes))
	}
	if quotes[0].Symbol != "AAPL" {
		t.Fatalf("quotes[0].Symbol = %q, want AAPL", quotes[0].Symbol)
	}
}
