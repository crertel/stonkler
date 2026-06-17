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

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
