package cli

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/crertel/stonkler/internal/fmp"
)

func TestParseWatchOptions(t *testing.T) {
	var stderr bytes.Buffer

	options, ok := parseWatchOptions([]string{"AAPL", "--interval", "2s", "MSFT", "--count", "3", "--sort", "-change-percent", "--fields", "symbol,price,change-percent", "--jsonl"}, &stderr)

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
	if options.stream {
		t.Fatalf("stream = true, want false")
	}
	if options.sort != "-change-percent" {
		t.Fatalf("sort = %q, want -change-percent", options.sort)
	}
	if len(options.fields) != 3 || options.fields[0] != "symbol" || options.fields[1] != "price" || options.fields[2] != "change-percent" {
		t.Fatalf("fields = %#v, want symbol/price/change-percent", options.fields)
	}
	if got := len(options.symbols); got != 2 {
		t.Fatalf("len(symbols) = %d, want 2", got)
	}
	if options.symbols[0] != "AAPL" || options.symbols[1] != "MSFT" {
		t.Fatalf("symbols = %#v, want AAPL/MSFT", options.symbols)
	}
}

func TestParseWatchOptionsSupportsStream(t *testing.T) {
	var stderr bytes.Buffer

	options, ok := parseWatchOptions([]string{"AAPL", "--stream"}, &stderr)

	if !ok {
		t.Fatalf("parseWatchOptions() ok = false, stderr = %q", stderr.String())
	}
	if !options.stream {
		t.Fatalf("stream = false, want true")
	}
}

func TestRunFundsWatchRejectsStream(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runFunds(nil, []string{"watch", "SPY", "--stream"}, &stdout, &stderr, func(string) string { return "key" })

	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if got := stderr.String(); got != "funds watch does not support --stream\n" {
		t.Fatalf("stderr = %q", got)
	}
}

func TestParseWatchOptionsRejectsInvalidInterval(t *testing.T) {
	var stderr bytes.Buffer

	_, ok := parseWatchOptions([]string{"AAPL", "--interval", "0s"}, &stderr)

	if ok {
		t.Fatalf("parseWatchOptions() ok = true, want false")
	}
}

func TestParseWatchOptionsRejectsInvalidSort(t *testing.T) {
	var stderr bytes.Buffer

	_, ok := parseWatchOptions([]string{"AAPL", "--sort", "market-cap"}, &stderr)

	if ok {
		t.Fatalf("parseWatchOptions() ok = true, want false")
	}
}

func TestParseWatchOptionsRejectsInvalidFields(t *testing.T) {
	var stderr bytes.Buffer

	_, ok := parseWatchOptions([]string{"AAPL", "--fields", "symbol,pe"}, &stderr)

	if ok {
		t.Fatalf("parseWatchOptions() ok = true, want false")
	}
}

func TestSortWatchQuotes(t *testing.T) {
	quotes := []fmp.Quote{
		{Symbol: "MSFT", ChangePercentage: -1.2, Volume: 20},
		{Symbol: "AAPL", ChangePercentage: 0.8, Volume: 30},
		{Symbol: "NVDA", ChangePercentage: 2.1, Volume: 10},
	}

	sortWatchQuotes(quotes, "-change-percent")

	if quotes[0].Symbol != "NVDA" || quotes[1].Symbol != "AAPL" || quotes[2].Symbol != "MSFT" {
		t.Fatalf("sorted symbols = %s/%s/%s, want NVDA/AAPL/MSFT", quotes[0].Symbol, quotes[1].Symbol, quotes[2].Symbol)
	}
}

func TestWriteWatchQuotesTableUsesSelectedFields(t *testing.T) {
	var stdout bytes.Buffer
	quotes := []fmp.Quote{{
		Symbol:           "AAPL",
		Name:             "Apple Inc.",
		Price:            295.95,
		ChangePercentage: -1.09945,
	}}

	err := writeWatchQuotesTable(&stdout, quotes, []string{"symbol", "price", "change-percent"})
	if err != nil {
		t.Fatalf("writeWatchQuotesTable() error = %v", err)
	}
	got := stdout.String()
	if got != "SYMBOL  PRICE   CHANGE%\nAAPL    295.95  -1.09945\n" {
		t.Fatalf("stdout = %q", got)
	}
}

func TestApplyStreamTradeUpdatesQuote(t *testing.T) {
	quotes := map[string]fmp.Quote{
		"AAPL": {
			Symbol:    "AAPL",
			Name:      "Apple Inc.",
			Price:     100,
			Change:    1,
			Volume:    10,
			Timestamp: 1,
		},
	}

	applyStreamTrade(quotes, fmp.StreamTrade{
		Symbol:    "aapl",
		Price:     102,
		Size:      5,
		Timestamp: 1710000000123,
	})

	quote := quotes["AAPL"]
	if quote.Price != 102 {
		t.Fatalf("price = %v, want 102", quote.Price)
	}
	if quote.Change != 3 {
		t.Fatalf("change = %v, want 3", quote.Change)
	}
	if math.Abs(quote.ChangePercentage-3.0303030303030303) > 0.0000000001 {
		t.Fatalf("change percent = %v, want 3.0303030303030303", quote.ChangePercentage)
	}
	if quote.Volume != 15 {
		t.Fatalf("volume = %v, want 15", quote.Volume)
	}
	if quote.Timestamp != 1710000000 {
		t.Fatalf("timestamp = %d, want 1710000000", quote.Timestamp)
	}
}
