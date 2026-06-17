package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type quoteFetcher func(context.Context, *fmp.Client, []string) ([]fmp.Quote, error)

func runDomainQuote(
	ctx context.Context,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
	getenv getenvFunc,
	domain string,
	help func(io.Writer),
	fetch quoteFetcher,
) int {
	if len(args) == 0 {
		help(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		help(stdout)
		return 0
	}

	format, symbols, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(symbols) == 0 {
		fmt.Fprintf(stderr, "%s quote requires at least one symbol\n", domain)
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	quotes, err := fetch(ctx, client, symbols)
	if err != nil {
		fmt.Fprintf(stderr, "%s quote failed: %v\n", domain, err)
		return 1
	}

	if err := writeStockQuotes(stdout, quotes, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func runQuoteWatchCommand(
	ctx context.Context,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
	getenv getenvFunc,
	domain string,
	help func(io.Writer),
	fetch quoteFetcher,
) int {
	if len(args) == 0 {
		help(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		help(stdout)
		return 0
	}

	options, ok := parseWatchOptions(args, stderr)
	if !ok {
		return 2
	}
	if len(options.symbols) == 0 {
		fmt.Fprintf(stderr, "%s watch requires at least one symbol\n", domain)
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	return runQuoteWatchLoop(ctx, stdout, stderr, client, options, fetch)
}
