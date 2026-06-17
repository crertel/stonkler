package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

func TestRunStocksTranscriptMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runStocksTranscript(context.Background(), []string{"AAPL", "--year", "2026", "--quarter", "1"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runStocksTranscript() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestWriteTranscriptsCSV(t *testing.T) {
	var stdout bytes.Buffer

	err := writeTranscripts(&stdout, []fmp.EarningsCallTranscript{{
		"symbol":  "AAPL",
		"year":    2026.0,
		"quarter": 1.0,
		"date":    "2026-01-30",
		"title":   "Apple Q1 2026 Earnings Call",
		"content": "Prepared remarks",
	}}, outputCSV)

	if err != nil {
		t.Fatalf("writeTranscripts() error = %v", err)
	}
	if !strings.Contains(stdout.String(), "symbol,year,quarter,date,title,content") {
		t.Fatalf("stdout = %q, want transcript CSV header", stdout.String())
	}
	if !strings.Contains(stdout.String(), "AAPL,2026,1,2026-01-30") {
		t.Fatalf("stdout = %q, want transcript CSV row", stdout.String())
	}
}
