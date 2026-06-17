package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type getSectorWeightingsClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	ETFSectorWeightings(context.Context, string) ([]fmp.ETFSectorWeighting, error)
}

func runGetSectorWeightings(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetSectorWeightingsHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetSectorWeightingsHelp(stdout)
		return 0
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) != 1 {
		fmt.Fprintln(stderr, "get sector-weightings requires exactly one symbol or name")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getSectorWeightingsSymbol(ctx, client, remaining[0])
	if err != nil {
		fmt.Fprintf(stderr, "get sector-weightings failed: %v\n", err)
		return 1
	}

	weightings, err := client.ETFSectorWeightings(ctx, symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get sector-weightings failed: %v\n", err)
		return 1
	}

	if err := writeETFSectorWeightings(stdout, weightings, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetSectorWeightingsHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch ETF or fund sector allocation weights, resolving a name query when needed.

Usage:
  stonk get sector-weightings <symbol|name> [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func getSectorWeightingsSymbol(ctx context.Context, client getSectorWeightingsClient, query string) (string, error) {
	if !shouldResolveGetFundNameQuery(query) {
		return normalizeGetFundTicker(query), nil
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
