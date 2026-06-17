package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/crertel/stonkler/internal/fmp"
)

type getHoldingsClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	ETFHoldings(context.Context, string) ([]fmp.ETFHolding, error)
}

func runGetHoldings(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetHoldingsHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetHoldingsHelp(stdout)
		return 0
	}

	options, ok := parseHoldingsOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getHoldingsSymbol(ctx, client, options.symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get holdings failed: %v\n", err)
		return 1
	}

	holdings, err := client.ETFHoldings(ctx, symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get holdings failed: %v\n", err)
		return 1
	}
	if options.limit > 0 && len(holdings) > options.limit {
		holdings = holdings[:options.limit]
	}

	if err := writeETFHoldings(stdout, holdings, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetHoldingsHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch ETF holdings, resolving a name query when needed.

Usage:
  stonk get holdings <symbol|name> [flags]

Flags:
  --limit <n>  Maximum holdings to print
  --json       Write JSON output
  --csv        Write CSV output
`)
}

func getHoldingsSymbol(ctx context.Context, client getHoldingsClient, query string) (string, error) {
	if !shouldResolveGetFundNameQuery(query) {
		return strings.ToUpper(strings.TrimSpace(query)), nil
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

func shouldResolveGetFundNameQuery(query string) bool {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return false
	}
	if strings.ContainsAny(trimmed, " \t") {
		return true
	}
	return len(trimmed) > 5 && trimmed != strings.ToUpper(trimmed)
}
