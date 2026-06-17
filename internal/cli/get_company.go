package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type getCompanyClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	StockProfile(context.Context, string) (fmp.StockProfile, error)
}

func runGetCompany(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetCompanyHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetCompanyHelp(stdout)
		return 0
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) != 1 {
		fmt.Fprintln(stderr, "get company requires exactly one symbol or name")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getCompanySymbol(ctx, client, remaining[0])
	if err != nil {
		fmt.Fprintf(stderr, "get company failed: %v\n", err)
		return 1
	}

	profile, err := client.StockProfile(ctx, symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get company failed: %v\n", err)
		return 1
	}

	if err := writeStockProfile(stdout, profile, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetCompanyHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch company profile data, resolving a name query when needed.

Usage:
  stonk get company <symbol|name> [flags]
  stonk get profile <symbol|name> [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func getCompanySymbol(ctx context.Context, client getCompanyClient, query string) (string, error) {
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
