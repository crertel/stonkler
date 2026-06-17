package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunGetHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runGet(context.Background(), []string{"--help"}, &stdout, &stderr, func(string) string { return "" })

	if code != 0 {
		t.Fatalf("runGet() code = %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "stonk get <command>") {
		t.Fatalf("stdout = %q, want get help", stdout.String())
	}
}

func TestRunGetQuoteUsesStocksQuoteValidation(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runGet(context.Background(), []string{"quote", "AAPL"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runGet() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing FMP key error", stderr.String())
	}
}

func TestRunGetCompanyUsesStocksProfileValidation(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runGet(context.Background(), []string{"company", "AAPL"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runGet() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing FMP key error", stderr.String())
	}
}

func TestRunGetHistoryUsesStocksHistoryValidation(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runGet(context.Background(), []string{"history", "AAPL"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runGet() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing FMP key error", stderr.String())
	}
}

func TestRunGetStatementsUsesStocksStatementsValidation(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runGet(context.Background(), []string{"statements", "AAPL", "income"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runGet() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing FMP key error", stderr.String())
	}
}

func TestRunGetHoldingsUsesFundsHoldingsValidation(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runGet(context.Background(), []string{"holdings", "SPY"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runGet() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing FMP key error", stderr.String())
	}
}

func TestRunGetFundUsesFundsInfoValidation(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runGet(context.Background(), []string{"fund", "SPY"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runGet() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing FMP key error", stderr.String())
	}
}

func TestRunGetCrossAssetQuoteValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "crypto", args: []string{"crypto", "BTCUSD"}},
		{name: "forex", args: []string{"forex", "EURUSD"}},
		{name: "commodity", args: []string{"commodity", "GCUSD"}},
		{name: "index", args: []string{"index", "GSPC"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer

			code := runGet(context.Background(), tt.args, &stdout, &stderr, func(string) string { return "" })

			if code != 1 {
				t.Fatalf("runGet() code = %d, want 1", code)
			}
			if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
				t.Fatalf("stderr = %q, want missing FMP key error", stderr.String())
			}
		})
	}
}

func TestRunGetCrossAssetHistoryValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "crypto-history", args: []string{"crypto-history", "BTCUSD"}},
		{name: "forex-history", args: []string{"forex-history", "EURUSD"}},
		{name: "commodity-history", args: []string{"commodity-history", "GCUSD"}},
		{name: "index-history", args: []string{"index-history", "GSPC"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer

			code := runGet(context.Background(), tt.args, &stdout, &stderr, func(string) string { return "" })

			if code != 1 {
				t.Fatalf("runGet() code = %d, want 1", code)
			}
			if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
				t.Fatalf("stderr = %q, want missing FMP key error", stderr.String())
			}
		})
	}
}
