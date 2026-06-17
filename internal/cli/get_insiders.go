package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type getInsidersClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	InsiderTrades(context.Context, string, int) ([]fmp.InsiderTrade, error)
}

func runGetInsiders(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetInsidersHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetInsidersHelp(stdout)
		return 0
	}

	options, ok := parseInsidersOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getInsidersSymbol(ctx, client, options.symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get insiders failed: %v\n", err)
		return 1
	}

	trades, err := client.InsiderTrades(ctx, symbol, options.limit)
	if err != nil {
		fmt.Fprintf(stderr, "get insiders failed: %v\n", err)
		return 1
	}

	if err := writeInsiderTrades(stdout, trades, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetInsidersHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch insider transactions, resolving a name query when needed.

Usage:
  stonk get insiders <symbol|name> [flags]

Flags:
  --limit <n>  Maximum transactions to request
  --json       Write JSON output
  --csv        Write CSV output
`)
}

func getInsidersSymbol(ctx context.Context, client getInsidersClient, query string) (string, error) {
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
