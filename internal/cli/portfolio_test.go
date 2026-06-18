package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPortfolioShowUsesBasisFlag(t *testing.T) {
	path := writeTestPortfolioFile(t)
	var stdout, stderr bytes.Buffer

	code := runPortfolio(nil, []string{"show", "--basis", path}, &stdout, &stderr, func(string) string { return "" })

	if code != 0 {
		t.Fatalf("code = %d, stderr = %q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	got := stdout.String()
	if !strings.Contains(got, "stocks") || !strings.Contains(got, "AAPL") || !strings.Contains(got, "110") {
		t.Fatalf("stdout = %q, want stocks AAPL basis", got)
	}
}

func TestRunPortfolioShowUsesConfiguredPortfolioFile(t *testing.T) {
	path := writeTestPortfolioFile(t)
	var stdout, stderr bytes.Buffer

	code := runPortfolio(nil, []string{"show", "--csv"}, &stdout, &stderr, func(key string) string {
		if key == "STONK_PORTFOLIO_FILE" {
			return path
		}
		return ""
	})

	if code != 0 {
		t.Fatalf("code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "stocks,AAPL,2,110,3,330,2024-01-01") {
		t.Fatalf("stdout = %q, want CSV summary", stdout.String())
	}
}

func TestRunPortfolioShowRequiresBasisFile(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runPortfolio(nil, []string{"show"}, &stdout, &stderr, func(string) string { return "" })

	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "--basis is required") {
		t.Fatalf("stderr = %q, want basis error", stderr.String())
	}
}

func TestParseBasisOutputOptions(t *testing.T) {
	var stderr bytes.Buffer

	options, ok := parseBasisOutputOptions([]string{"AAPL", "--basis", "portfolio.json", "--json"}, &stderr)

	if !ok {
		t.Fatalf("parseBasisOutputOptions() ok = false, stderr = %q", stderr.String())
	}
	if options.format != outputJSON {
		t.Fatalf("format = %q, want json", options.format)
	}
	if options.basisPath != "portfolio.json" {
		t.Fatalf("basisPath = %q, want portfolio.json", options.basisPath)
	}
	if len(options.remaining) != 1 || options.remaining[0] != "AAPL" {
		t.Fatalf("remaining = %#v, want AAPL", options.remaining)
	}
}

func writeTestPortfolioFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "portfolio.json")
	content := `{
		"version": 1,
		"stocks": {
			"AAPL": {
				"lots": [
					{"basis": 100, "quantity": 2, "acquired_on": "2024-02-01"},
					{"basis": 130, "quantity": 1, "acquired_on": "2024-01-01"}
				]
			}
		}
	}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}
