package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type getExposureClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	ETFAssetExposure(context.Context, string) ([]fmp.ETFAssetExposure, error)
}

func runGetExposure(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetExposureHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetExposureHelp(stdout)
		return 0
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) != 1 {
		fmt.Fprintln(stderr, "get exposure requires exactly one asset symbol or name")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getExposureSymbol(ctx, client, remaining[0])
	if err != nil {
		fmt.Fprintf(stderr, "get exposure failed: %v\n", err)
		return 1
	}

	exposures, err := client.ETFAssetExposure(ctx, symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get exposure failed: %v\n", err)
		return 1
	}

	if err := writeETFAssetExposures(stdout, exposures, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetExposureHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch ETF or fund exposure to an asset, resolving a name query when needed.

Usage:
  stonk get exposure <asset-symbol|asset-name> [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func getExposureSymbol(ctx context.Context, client getExposureClient, query string) (string, error) {
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
