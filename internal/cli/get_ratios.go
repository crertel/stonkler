package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type getRatiosClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	StockRatiosTTM(context.Context, string) ([]fmp.StockRatioRow, error)
}

func runGetRatios(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetRatiosHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetRatiosHelp(stdout)
		return 0
	}

	options, ok := parseRatiosOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getRatiosSymbol(ctx, client, options.symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get ratios failed: %v\n", err)
		return 1
	}

	ratios, err := client.StockRatiosTTM(ctx, symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get ratios failed: %v\n", err)
		return 1
	}

	if err := writeStockRatios(stdout, ratios, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetRatiosHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch stock ratios, resolving a name query when needed.

Usage:
  stonk get ratios <symbol|name> [flags]

Flags:
  --ttm   Fetch trailing-twelve-month ratios
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func getRatiosSymbol(ctx context.Context, client getRatiosClient, query string) (string, error) {
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
