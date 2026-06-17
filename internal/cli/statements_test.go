package cli

import (
	"bytes"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

func TestParseStatementsOptions(t *testing.T) {
	var stderr bytes.Buffer

	options, ok := parseStatementsOptions([]string{"AAPL", "income", "--period", "quarter", "--limit", "2", "--csv"}, &stderr)

	if !ok {
		t.Fatalf("parseStatementsOptions() ok = false, stderr = %q", stderr.String())
	}
	if options.symbol != "AAPL" {
		t.Fatalf("symbol = %q, want AAPL", options.symbol)
	}
	if options.statement != fmp.StatementIncome {
		t.Fatalf("statement = %q, want income", options.statement)
	}
	if options.period != "quarter" {
		t.Fatalf("period = %q, want quarter", options.period)
	}
	if options.limit != 2 {
		t.Fatalf("limit = %d, want 2", options.limit)
	}
	if options.format != outputCSV {
		t.Fatalf("format = %q, want csv", options.format)
	}
}

func TestParseStatementsOptionsRejectsUnknownType(t *testing.T) {
	var stderr bytes.Buffer

	_, ok := parseStatementsOptions([]string{"AAPL", "nonsense"}, &stderr)

	if ok {
		t.Fatalf("parseStatementsOptions() ok = true, want false")
	}
}
