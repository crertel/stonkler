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
