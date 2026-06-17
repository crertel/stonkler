package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunStocksInsidersMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runStocksInsiders(context.Background(), []string{"AAPL"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runStocksInsiders() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}
