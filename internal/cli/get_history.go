package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type getHistoryClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	StockHistory(context.Context, fmp.StockHistoryRequest) ([]fmp.StockPrice, error)
}

func runGetHistory(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetHistoryHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetHistoryHelp(stdout)
		return 0
	}

	options, ok := parseHistoryOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getHistorySymbol(ctx, client, options.symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get history failed: %v\n", err)
		return 1
	}
	options.symbol = symbol

	prices, err := client.StockHistory(ctx, fmp.StockHistoryRequest{
		Symbol: options.symbol,
		From:   options.from,
		To:     options.to,
	})
	if err != nil {
		fmt.Fprintf(stderr, "get history failed: %v\n", err)
		return 1
	}
	if options.limit > 0 && len(prices) > options.limit {
		prices = prices[:options.limit]
	}

	if err := writeStockHistory(stdout, prices, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetHistoryHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch historical end-of-day prices, resolving a name query when needed.

Usage:
  stonk get history <symbol|name> [flags]

Flags:
  --from <date>  Start date in YYYY-MM-DD format
  --to <date>    End date in YYYY-MM-DD format
  --limit <n>    Maximum rows to print
  --json         Write JSON output
  --csv          Write CSV output
`)
}

func getHistorySymbol(ctx context.Context, client getHistoryClient, query string) (string, error) {
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
