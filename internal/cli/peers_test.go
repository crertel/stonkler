package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunStocksPeersMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runStocksPeers(context.Background(), []string{"AAPL"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runStocksPeers() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stderr = %q, want missing key error", stderr.String())
	}
}
