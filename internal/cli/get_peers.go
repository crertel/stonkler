package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type getPeersClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	StockPeers(context.Context, string) ([]fmp.StockPeer, error)
}

func runGetPeers(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetPeersHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetPeersHelp(stdout)
		return 0
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) != 1 {
		fmt.Fprintln(stderr, "get peers requires exactly one symbol or name")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getPeersSymbol(ctx, client, remaining[0])
	if err != nil {
		fmt.Fprintf(stderr, "get peers failed: %v\n", err)
		return 1
	}

	peers, err := client.StockPeers(ctx, symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get peers failed: %v\n", err)
		return 1
	}

	if err := writeStockPeers(stdout, peers, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetPeersHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch peer companies, resolving a name query when needed.

Usage:
  stonk get peers <symbol|name> [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func getPeersSymbol(ctx context.Context, client getPeersClient, query string) (string, error) {
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
