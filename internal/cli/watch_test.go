package cli

import (
	"bytes"
	"testing"
	"time"
)

func TestParseWatchOptions(t *testing.T) {
	var stderr bytes.Buffer

	options, ok := parseWatchOptions([]string{"AAPL", "--interval", "2s", "MSFT", "--count", "3", "--jsonl"}, &stderr)

	if !ok {
		t.Fatalf("parseWatchOptions() ok = false, stderr = %q", stderr.String())
	}
	if options.interval != 2*time.Second {
		t.Fatalf("interval = %v, want 2s", options.interval)
	}
	if options.count != 3 {
		t.Fatalf("count = %d, want 3", options.count)
	}
	if !options.jsonl {
		t.Fatalf("jsonl = false, want true")
	}
	if got := len(options.symbols); got != 2 {
		t.Fatalf("len(symbols) = %d, want 2", got)
	}
	if options.symbols[0] != "AAPL" || options.symbols[1] != "MSFT" {
		t.Fatalf("symbols = %#v, want AAPL/MSFT", options.symbols)
	}
}

func TestParseWatchOptionsRejectsInvalidInterval(t *testing.T) {
	var stderr bytes.Buffer

	_, ok := parseWatchOptions([]string{"AAPL", "--interval", "0s"}, &stderr)

	if ok {
		t.Fatalf("parseWatchOptions() ok = true, want false")
	}
}
