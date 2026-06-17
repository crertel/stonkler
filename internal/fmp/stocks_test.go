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

func TestBatchQuotesHandlesNullMarketCapAndChangeFieldVariants(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"EURUSD","name":"EUR/USD","price":1.14,"change":-0.01,"changePercentage":-1.2,"volume":155095,"marketCap":null,"timestamp":1781727386},{"symbol":"^GSPC","name":"S&P 500","price":7421.76,"changesPercentage":-1.19,"change":-89.59,"volume":2624646000,"marketCap":0,"timestamp":1781726399}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	quotes, err := client.BatchQuotes(context.Background(), []string{"EURUSD", "^GSPC"})
	if err != nil {
		t.Fatalf("BatchQuotes() error = %v", err)
	}
	if quotes[0].MarketCap != 0 {
		t.Fatalf("quotes[0].MarketCap = %v, want 0", quotes[0].MarketCap)
	}
	if quotes[1].ChangePercentage != -1.19 {
		t.Fatalf("quotes[1].ChangePercentage = %v, want -1.19", quotes[1].ChangePercentage)
	}
}

func TestIndexQuotesUsesV3QuoteEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.URL.Path; got != "/quote/^GSPC,^DJI" {
			t.Fatalf("path = %q, want /quote/^GSPC,^DJI", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"^GSPC","name":"S&P 500","price":7421.76,"changesPercentage":-1.19,"change":-89.59}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.v3BaseURL = "https://example.test"

	quotes, err := client.IndexQuotes(context.Background(), []string{"gspc", "^dji"})
	if err != nil {
		t.Fatalf("IndexQuotes() error = %v", err)
	}
	if len(quotes) != 1 {
		t.Fatalf("len(quotes) = %d, want 1", len(quotes))
	}
}
