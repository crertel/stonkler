package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

func TestRunStocksRatiosMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runStocksRatios(context.Background(), []string{"AAPL", "--ttm"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runStocksRatios() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestWriteStockRatiosCSV(t *testing.T) {
	var stdout bytes.Buffer

	err := writeStockRatios(&stdout, []fmp.StockRatioRow{{
		"symbol":                   "AAPL",
		"netProfitMarginTTM":       0.27,
		"priceToEarningsRatioTTM":  35.5,
		"debtToEquityRatioTTM":     0.79,
		"grossProfitMarginTTM":     0.47,
		"operatingProfitMarginTTM": 0.32,
		"currentRatioTTM":          1.07,
		"quickRatioTTM":            1.02,
		"priceToBookRatioTTM":      40.8,
		"priceToSalesRatioTTM":     9.6,
		"dividendYieldTTM":         0.003,
	}}, outputCSV)

	if err != nil {
		t.Fatalf("writeStockRatios() error = %v", err)
	}
	if !strings.Contains(stdout.String(), "symbol,grossProfitMarginTTM") {
		t.Fatalf("stdout = %q, want ratios CSV header", stdout.String())
	}
	if !strings.Contains(stdout.String(), "AAPL,0.47") {
		t.Fatalf("stdout = %q, want ratios CSV row", stdout.String())
	}
}
