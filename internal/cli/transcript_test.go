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

func TestRunStocksTranscriptLatestMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runStocksTranscript(context.Background(), []string{"AAPL", "--latest"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runStocksTranscript() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestParseTranscriptOptionsLatestRejectsExplicitPeriod(t *testing.T) {
	var stderr bytes.Buffer

	_, ok := parseTranscriptOptions([]string{"AAPL", "--latest", "--year", "2026"}, &stderr)

	if ok {
		t.Fatal("parseTranscriptOptions() ok = true, want false")
	}
	if !strings.Contains(stderr.String(), "--latest cannot be combined") {
		t.Fatalf("stderr = %q, want latest conflict error", stderr.String())
	}
}

func TestLatestTranscriptPeriodPrefersNewestDate(t *testing.T) {
	year, quarter, ok := latestTranscriptPeriod([]fmp.EarningsCallTranscriptDate{
		{Symbol: "AAPL", Year: 2026, Quarter: 1, Date: "2026-01-30"},
		{Symbol: "AAPL", Year: 2025, Quarter: 4, Date: "2025-10-30"},
		{Symbol: "AAPL", Year: 2026, Quarter: 2, Date: "2026-04-30"},
	})

	if !ok {
		t.Fatal("latestTranscriptPeriod() ok = false, want true")
	}
	if year != 2026 || quarter != 2 {
		t.Fatalf("latestTranscriptPeriod() = %d Q%d, want 2026 Q2", year, quarter)
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
