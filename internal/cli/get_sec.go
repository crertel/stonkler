package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type getSECClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	SECFilings(context.Context, string, int) ([]fmp.SECFiling, error)
}

func runGetSEC(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetSECHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetSECHelp(stdout)
		return 0
	}

	options, ok := parseSECOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getSECSymbol(ctx, client, options.symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get sec failed: %v\n", err)
		return 1
	}

	filings, err := client.SECFilings(ctx, symbol, options.limit)
	if err != nil {
		fmt.Fprintf(stderr, "get sec failed: %v\n", err)
		return 1
	}

	if err := writeSECFilings(stdout, filings, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetSECHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch SEC filings, resolving a name query when needed.

Usage:
  stonk get sec <symbol|name> [flags]

Flags:
  --limit <n>  Maximum filings to request
  --json       Write JSON output
  --csv        Write CSV output
`)
}

func getSECSymbol(ctx context.Context, client getSECClient, query string) (string, error) {
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
