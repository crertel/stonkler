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

func TestStockRatiosTTMUsesStableEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/ratios-ttm" {
			t.Fatalf("path = %q, want /ratios-ttm", got)
		}
		if got := r.URL.Query().Get("symbol"); got != "AAPL" {
			t.Fatalf("symbol query = %q, want AAPL", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"AAPL","netProfitMarginTTM":0.27}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	ratios, err := client.StockRatiosTTM(context.Background(), "aapl")
	if err != nil {
		t.Fatalf("StockRatiosTTM() error = %v", err)
	}
	if len(ratios) != 1 {
		t.Fatalf("len(ratios) = %d, want 1", len(ratios))
	}
	if ratios[0]["symbol"] != "AAPL" {
		t.Fatalf("ratios[0][symbol] = %q, want AAPL", ratios[0]["symbol"])
	}
}

func TestStockKeyMetricsTTMUsesStableEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/key-metrics-ttm" {
			t.Fatalf("path = %q, want /key-metrics-ttm", got)
		}
		if got := r.URL.Query().Get("symbol"); got != "AAPL" {
			t.Fatalf("symbol query = %q, want AAPL", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"AAPL","marketCap":4346723008200}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	metrics, err := client.StockKeyMetricsTTM(context.Background(), "aapl")
	if err != nil {
		t.Fatalf("StockKeyMetricsTTM() error = %v", err)
	}
	if len(metrics) != 1 {
		t.Fatalf("len(metrics) = %d, want 1", len(metrics))
	}
	if metrics[0]["symbol"] != "AAPL" {
		t.Fatalf("metrics[0][symbol] = %q, want AAPL", metrics[0]["symbol"])
	}
}

func TestStockPeersUsesStableEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/stock-peers" {
			t.Fatalf("path = %q, want /stock-peers", got)
		}
		if got := r.URL.Query().Get("symbol"); got != "AAPL" {
			t.Fatalf("symbol query = %q, want AAPL", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"MSFT","companyName":"Microsoft Corporation","price":378.91,"mktCap":2814706411300}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	peers, err := client.StockPeers(context.Background(), "aapl")
	if err != nil {
		t.Fatalf("StockPeers() error = %v", err)
	}
	if len(peers) != 1 {
		t.Fatalf("len(peers) = %d, want 1", len(peers))
	}
	if peers[0].Symbol != "MSFT" {
		t.Fatalf("peers[0].Symbol = %q, want MSFT", peers[0].Symbol)
	}
}

func TestStockRatingSnapshotUsesStableEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("apikey"); got != "secret-value" {
			t.Fatalf("apikey header = %q, want secret-value", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "" {
			t.Fatalf("apikey query = %q, want empty", got)
		}
		if got := r.URL.Path; got != "/ratings-snapshot" {
			t.Fatalf("path = %q, want /ratings-snapshot", got)
		}
		if got := r.URL.Query().Get("symbol"); got != "AAPL" {
			t.Fatalf("symbol query = %q, want AAPL", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"symbol":"AAPL","rating":"B","overallScore":3,"discountedCashFlowScore":3,"returnOnEquityScore":5,"returnOnAssetsScore":5,"debtToEquityScore":1,"priceToEarningsScore":2,"priceToBookScore":1}]`)),
		}, nil
	})

	client := NewClient("secret-value", &http.Client{Transport: transport})
	client.baseURL = "https://example.test"

	ratings, err := client.StockRatingSnapshot(context.Background(), "aapl")
	if err != nil {
		t.Fatalf("StockRatingSnapshot() error = %v", err)
	}
	if len(ratings) != 1 {
		t.Fatalf("len(ratings) = %d, want 1", len(ratings))
	}
	if ratings[0].Rating != "B" {
		t.Fatalf("ratings[0].Rating = %q, want B", ratings[0].Rating)
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
