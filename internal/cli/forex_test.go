package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunForexHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runForex(context.Background(), []string{"--help"}, &stdout, &stderr, func(string) string { return "" })

	if code != 0 {
		t.Fatalf("runForex() code = %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "stonk forex <command>") {
		t.Fatalf("stdout = %q, want forex help", stdout.String())
	}
}

func TestRunForexQuoteMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runForex(context.Background(), []string{"quote", "EURUSD"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runForex() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestRunForexHistoryMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runForex(context.Background(), []string{"history", "EURUSD"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runForex() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}
