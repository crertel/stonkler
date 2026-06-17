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

func TestETFSectorWeightingsUsesStableEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/etf/sector-weightings" {
			t.Fatalf("path = %q, want /etf/sector-weightings", got)
		}
		if got := r.URL.Query().Get("symbol"); got != "SPY" {
			t.Fatalf("symbol query = %q, want SPY", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"SPY","sector":"Technology","weightPercentage":39.05}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	weightings, err := client.ETFSectorWeightings(context.Background(), "spy")
	if err != nil {
		t.Fatalf("ETFSectorWeightings() error = %v", err)
	}
	if len(weightings) != 1 {
		t.Fatalf("len(weightings) = %d, want 1", len(weightings))
	}
	if weightings[0].Sector != "Technology" {
		t.Fatalf("weightings[0].Sector = %q, want Technology", weightings[0].Sector)
	}
}

func TestETFCountryWeightingsUsesStableEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/etf/country-weightings" {
			t.Fatalf("path = %q, want /etf/country-weightings", got)
		}
		if got := r.URL.Query().Get("symbol"); got != "VXUS" {
			t.Fatalf("symbol query = %q, want VXUS", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"country":"Japan","weightPercentage":"14.49%"}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	weightings, err := client.ETFCountryWeightings(context.Background(), "vxus")
	if err != nil {
		t.Fatalf("ETFCountryWeightings() error = %v", err)
	}
	if len(weightings) != 1 {
		t.Fatalf("len(weightings) = %d, want 1", len(weightings))
	}
	if weightings[0].Country != "Japan" {
		t.Fatalf("weightings[0].Country = %q, want Japan", weightings[0].Country)
	}
}

func TestETFAssetExposureUsesStableEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/etf/asset-exposure" {
			t.Fatalf("path = %q, want /etf/asset-exposure", got)
		}
		if got := r.URL.Query().Get("symbol"); got != "AAPL" {
			t.Fatalf("symbol query = %q, want AAPL", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"SPY","asset":"AAPL","name":"APPLE INC","sharesNumber":176183317,"weightPercentage":6.81,"marketValue":53568630492}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	exposures, err := client.ETFAssetExposure(context.Background(), "aapl")
	if err != nil {
		t.Fatalf("ETFAssetExposure() error = %v", err)
	}
	if len(exposures) != 1 {
		t.Fatalf("len(exposures) = %d, want 1", len(exposures))
	}
	if exposures[0]["symbol"] != "SPY" {
		t.Fatalf("exposures[0][symbol] = %q, want SPY", exposures[0]["symbol"])
	}
}
