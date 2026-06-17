package cli

import (
	"bytes"
	"testing"
)

func TestParseHoldingsOptions(t *testing.T) {
	var stderr bytes.Buffer

	options, ok := parseHoldingsOptions([]string{"SPY", "--limit", "10", "--csv"}, &stderr)

	if !ok {
		t.Fatalf("parseHoldingsOptions() ok = false, stderr = %q", stderr.String())
	}
	if options.symbol != "SPY" {
		t.Fatalf("symbol = %q, want SPY", options.symbol)
	}
	if options.limit != 10 {
		t.Fatalf("limit = %d, want 10", options.limit)
	}
	if options.format != outputCSV {
		t.Fatalf("format = %q, want csv", options.format)
	}
}

func TestParseHoldingsOptionsRejectsDuplicateSymbol(t *testing.T) {
	var stderr bytes.Buffer

	_, ok := parseHoldingsOptions([]string{"SPY", "VTI"}, &stderr)

	if ok {
		t.Fatalf("parseHoldingsOptions() ok = true, want false")
	}
}
