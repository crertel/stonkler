package fmp

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSearchNameUsesHeaderAuth(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/search-name" {
			t.Fatalf("path = %q, want /search-name", got)
		}
		if got := r.URL.Query().Get("query"); got != "apple" {
			t.Fatalf("query = %q, want apple", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"AAPL","name":"Apple Inc.","currency":"USD","exchangeShortName":"NASDAQ"}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	results, err := client.SearchName(context.Background(), "apple")
	if err != nil {
		t.Fatalf("SearchName() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].Symbol != "AAPL" {
		t.Fatalf("results[0].Symbol = %q, want AAPL", results[0].Symbol)
	}
}

func TestCompanyScreenerUsesStableEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/company-screener" {
			t.Fatalf("path = %q, want /company-screener", got)
		}
		if got := r.URL.Query().Get("sector"); got != "Technology" {
			t.Fatalf("sector query = %q, want Technology", got)
		}
		if got := r.URL.Query().Get("country"); got != "US" {
			t.Fatalf("country query = %q, want US", got)
		}
		if got := r.URL.Query().Get("marketCapMoreThan"); got != "100000000000" {
			t.Fatalf("marketCapMoreThan query = %q, want 100000000000", got)
		}
		if got := r.URL.Query().Get("limit"); got != "3" {
			t.Fatalf("limit query = %q, want 3", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"NVDA","companyName":"NVIDIA Corporation","marketCap":5023677610000,"sector":"Technology","industry":"Semiconductors","price":204.65,"exchangeShortName":"NASDAQ","country":"US","isActivelyTrading":true}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	results, err := client.CompanyScreener(context.Background(), ScreenerOptions{
		Sector:       "Technology",
		Country:      "US",
		MarketCapMin: 100_000_000_000,
		Limit:        3,
	})
	if err != nil {
		t.Fatalf("CompanyScreener() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].Symbol != "NVDA" {
		t.Fatalf("results[0].Symbol = %q, want NVDA", results[0].Symbol)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
