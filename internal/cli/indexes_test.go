package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunIndexesHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runIndexes(context.Background(), []string{"--help"}, &stdout, &stderr, func(string) string { return "" })

	if code != 0 {
		t.Fatalf("runIndexes() code = %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "stonk indexes <command>") {
		t.Fatalf("stdout = %q, want indexes help", stdout.String())
	}
}

func TestRunIndexesQuoteMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runIndexes(context.Background(), []string{"quote", "GSPC"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runIndexes() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestRunIndexesHistoryMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runIndexes(context.Background(), []string{"history", "GSPC"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runIndexes() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestNormalizeIndexSymbol(t *testing.T) {
	if got := normalizeIndexSymbol("gspc"); got != "^GSPC" {
		t.Fatalf("normalizeIndexSymbol() = %q, want ^GSPC", got)
	}
	if got := normalizeIndexSymbol("^dji"); got != "^DJI" {
		t.Fatalf("normalizeIndexSymbol() = %q, want ^DJI", got)
	}
}
