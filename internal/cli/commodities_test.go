package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunCommoditiesHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runCommodities(context.Background(), []string{"--help"}, &stdout, &stderr, func(string) string { return "" })

	if code != 0 {
		t.Fatalf("runCommodities() code = %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "stonk commodities <command>") {
		t.Fatalf("stdout = %q, want commodities help", stdout.String())
	}
}

func TestRunCommoditiesQuoteMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runCommodities(context.Background(), []string{"quote", "GCUSD"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runCommodities() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestRunCommoditiesHistoryMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runCommodities(context.Background(), []string{"history", "GCUSD"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runCommodities() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}
