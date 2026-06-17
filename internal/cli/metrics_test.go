package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/crertel/stonkler/internal/fmp"
)

func TestRunStocksMetricsMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runStocksMetrics(context.Background(), []string{"AAPL", "--ttm"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runStocksMetrics() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}

func TestWriteStockMetricsCSV(t *testing.T) {
	var stdout bytes.Buffer

	err := writeStockMetrics(&stdout, []fmp.StockMetricRow{{
		"symbol":                     "AAPL",
		"marketCap":                  4346723008200.0,
		"enterpriseValueTTM":         4395106008200.0,
		"evToSalesTTM":               9.7,
		"evToEBITDATTM":              27.4,
		"netDebtToEBITDATTM":         0.3,
		"returnOnAssetsTTM":          0.33,
		"returnOnEquityTTM":          1.46,
		"returnOnInvestedCapitalTTM": 0.49,
		"freeCashFlowYieldTTM":       0.02,
		"cashConversionCycleTTM":     -35.2,
	}}, outputCSV)

	if err != nil {
		t.Fatalf("writeStockMetrics() error = %v", err)
	}
	if !strings.Contains(stdout.String(), "symbol,marketCap") {
		t.Fatalf("stdout = %q, want metrics CSV header", stdout.String())
	}
	if !strings.Contains(stdout.String(), "AAPL,4346723008200") {
		t.Fatalf("stdout = %q, want metrics CSV row", stdout.String())
	}
}
