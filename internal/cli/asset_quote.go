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

	options, ok := parseBasisOutputOptions(args, stderr)
	if !ok {
		return 2
	}
	format := options.format
	symbols := options.remaining
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

	var writeErr error
	if options.basisPath != "" {
		book, ok := loadBasisPathOption(options.basisPath, stderr, getenv)
		if !ok {
			return 2
		}
		writeErr = writeQuotesWithBasis(stdout, attachBasis(domain, quotes, book), format)
	} else {
		writeErr = writeStockQuotes(stdout, quotes, format)
	}
	if writeErr != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", writeErr)
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
	if options.stream {
		fmt.Fprintf(stderr, "%s watch does not support --stream\n", domain)
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	book, ok := loadOptionalWatchBasis(options.basisPath, stderr, getenv)
	if !ok {
		return 2
	}
	return runQuoteWatchLoop(ctx, stdout, stderr, client, options, domain, book, fetch)
}
