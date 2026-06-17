package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type getCountryWeightingsClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	ETFCountryWeightings(context.Context, string) ([]fmp.ETFCountryWeighting, error)
}

func runGetCountryWeightings(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetCountryWeightingsHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetCountryWeightingsHelp(stdout)
		return 0
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) != 1 {
		fmt.Fprintln(stderr, "get country-weightings requires exactly one symbol or name")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getCountryWeightingsSymbol(ctx, client, remaining[0])
	if err != nil {
		fmt.Fprintf(stderr, "get country-weightings failed: %v\n", err)
		return 1
	}

	weightings, err := client.ETFCountryWeightings(ctx, symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get country-weightings failed: %v\n", err)
		return 1
	}

	if err := writeETFCountryWeightings(stdout, weightings, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetCountryWeightingsHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch ETF or fund country allocation weights, resolving a name query when needed.

Usage:
  stonk get country-weightings <symbol|name> [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func getCountryWeightingsSymbol(ctx context.Context, client getCountryWeightingsClient, query string) (string, error) {
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
