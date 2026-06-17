package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunSearchFundsMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runSearch(context.Background(), []string{"funds", "spy"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runSearch() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestRunSearchScreenerMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runSearch(context.Background(), []string{"screener", "--sector", "Technology"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runSearch() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestParseMarketCapAmount(t *testing.T) {
	got, err := parseMarketCapAmount("100B")
	if err != nil {
		t.Fatalf("parseMarketCapAmount() error = %v", err)
	}
	if got != 100_000_000_000 {
		t.Fatalf("parseMarketCapAmount() = %v, want 100000000000", got)
	}
}
