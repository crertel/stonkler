package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/crertel/stonkler/internal/fmp"
)

type getQuoteClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	StockQuotes(context.Context, []string) ([]fmp.StockQuote, error)
}

func runGetQuote(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetQuoteHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetQuoteHelp(stdout)
		return 0
	}

	format, queries, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(queries) == 0 {
		fmt.Fprintln(stderr, "get quote requires at least one symbol or name")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbols, err := getQuoteSymbols(ctx, client, queries)
	if err != nil {
		fmt.Fprintf(stderr, "get quote failed: %v\n", err)
		return 1
	}

	quotes, err := client.StockQuotes(ctx, symbols)
	if err != nil {
		fmt.Fprintf(stderr, "get quote failed: %v\n", err)
		return 1
	}

	if err := writeStockQuotes(stdout, quotes, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetQuoteHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch one or more quotes, resolving a single name query when needed.

Usage:
  stonk get quote <symbol|name> [symbol...] [flags]
  stonk get quotes <symbol|name> [symbol...] [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func getQuoteSymbols(ctx context.Context, client getQuoteClient, queries []string) ([]string, error) {
	if len(queries) != 1 || !shouldResolveGetQuoteQuery(queries[0]) {
		return queries, nil
	}

	results, err := client.SearchName(ctx, queries[0])
	if err != nil {
		return nil, err
	}
	symbol, ok := firstSearchSymbol(results)
	if !ok {
		return nil, fmt.Errorf("no symbol found for %q", queries[0])
	}
	return []string{symbol}, nil
}

func shouldResolveGetQuoteQuery(query string) bool {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return false
	}
	return strings.ContainsAny(trimmed, " \t") || trimmed != strings.ToUpper(trimmed)
}

func firstSearchSymbol(results []fmp.SearchResult) (string, bool) {
	bestSymbol := ""
	bestScore := -1
	for _, result := range results {
		symbol := strings.TrimSpace(result.Symbol)
		if symbol == "" {
			continue
		}
		score := searchQuoteSymbolScore(result)
		if score > bestScore {
			bestSymbol = strings.ToUpper(symbol)
			bestScore = score
		}
	}
	if bestSymbol != "" {
		return bestSymbol, true
	}
	return "", false
}

func searchQuoteSymbolScore(result fmp.SearchResult) int {
	score := 0
	if strings.EqualFold(result.Currency, "USD") {
		score += 10
	}
	if !strings.Contains(result.Symbol, ".") {
		score += 5
	}
	switch strings.ToUpper(result.ExchangeShortName) {
	case "NASDAQ", "NYSE", "AMEX":
		score += 3
	}
	return score
}
