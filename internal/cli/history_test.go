package cli

import (
	"bytes"
	"testing"
)

func TestParseHistoryOptions(t *testing.T) {
	var stderr bytes.Buffer

	options, ok := parseHistoryOptions([]string{"AAPL", "--from", "2026-06-10", "--to", "2026-06-12", "--limit", "2", "--csv"}, &stderr)

	if !ok {
		t.Fatalf("parseHistoryOptions() ok = false, stderr = %q", stderr.String())
	}
	if options.symbol != "AAPL" {
		t.Fatalf("symbol = %q, want AAPL", options.symbol)
	}
	if options.from != "2026-06-10" || options.to != "2026-06-12" {
		t.Fatalf("date range = %q/%q, want expected range", options.from, options.to)
	}
	if options.limit != 2 {
		t.Fatalf("limit = %d, want 2", options.limit)
	}
	if options.format != outputCSV {
		t.Fatalf("format = %q, want csv", options.format)
	}
}

func TestParseHistoryOptionsRejectsBadDate(t *testing.T) {
	var stderr bytes.Buffer

	_, ok := parseHistoryOptions([]string{"AAPL", "--from", "06-10-2026"}, &stderr)

	if ok {
		t.Fatalf("parseHistoryOptions() ok = true, want false")
	}
}
