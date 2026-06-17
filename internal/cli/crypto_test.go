package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunCryptoHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runCrypto(context.Background(), []string{"--help"}, &stdout, &stderr, func(string) string { return "" })

	if code != 0 {
		t.Fatalf("runCrypto() code = %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "stonk crypto <command>") {
		t.Fatalf("stdout = %q, want crypto help", stdout.String())
	}
}

func TestRunCryptoQuoteMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runCrypto(context.Background(), []string{"quote", "BTCUSD"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runCrypto() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestRunCryptoHistoryMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runCrypto(context.Background(), []string{"history", "BTCUSD"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runCrypto() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestRunCryptoWatchMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runCrypto(context.Background(), []string{"watch", "BTCUSD"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runCrypto() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}
