package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type getMetricsClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	StockKeyMetricsTTM(context.Context, string) ([]fmp.StockMetricRow, error)
}

func runGetMetrics(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetMetricsHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetMetricsHelp(stdout)
		return 0
	}

	options, ok := parseMetricsOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getMetricsSymbol(ctx, client, options.symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get metrics failed: %v\n", err)
		return 1
	}

	metrics, err := client.StockKeyMetricsTTM(ctx, symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get metrics failed: %v\n", err)
		return 1
	}

	if err := writeStockMetrics(stdout, metrics, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetMetricsHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch stock key metrics, resolving a name query when needed.

Usage:
  stonk get metrics <symbol|name> [flags]

Flags:
  --ttm   Fetch trailing-twelve-month metrics
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func getMetricsSymbol(ctx context.Context, client getMetricsClient, query string) (string, error) {
	if !shouldResolveGetNameQuery(query) {
		return query, nil
	}

	results, err := client.SearchName(ctx, query)
	if err != nil {
		return "", err
	}
	symbol, ok := bestSearchSymbol(results)
	if !ok {
		return "", fmt.Errorf("no symbol found for %q", query)
	}
	return symbol, nil
}
