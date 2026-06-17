package fmp

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestStockProfileUsesProfileEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/profile" {
			t.Fatalf("path = %q, want /profile", got)
		}
		if got := r.URL.Query().Get("symbol"); got != "AAPL" {
			t.Fatalf("symbol = %q, want AAPL", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"AAPL","companyName":"Apple Inc.","price":295.25,"marketCap":4336441859000,"exchange":"NASDAQ","currency":"USD","sector":"Technology","industry":"Consumer Electronics","ceo":"Timothy D. Cook","country":"US","website":"https://www.apple.com","ipoDate":"1980-12-12"}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	profile, err := client.StockProfile(context.Background(), "aapl")
	if err != nil {
		t.Fatalf("StockProfile() error = %v", err)
	}
	if profile.Symbol != "AAPL" {
		t.Fatalf("profile.Symbol = %q, want AAPL", profile.Symbol)
	}
	if profile.CompanyName != "Apple Inc." {
		t.Fatalf("profile.CompanyName = %q, want Apple Inc.", profile.CompanyName)
	}
}
