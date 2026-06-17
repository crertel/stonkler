package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

func runDomainHistory(
	ctx context.Context,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
	getenv getenvFunc,
	domain string,
	help func(io.Writer),
	normalize func(string) string,
) int {
	if len(args) == 0 {
		help(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		help(stdout)
		return 0
	}

	options, ok := parseHistoryOptions(args, stderr)
	if !ok {
		return 2
	}
	if normalize != nil {
		options.symbol = normalize(options.symbol)
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	prices, err := client.PriceHistory(ctx, fmp.StockHistoryRequest{
		Symbol: options.symbol,
		From:   options.from,
		To:     options.to,
	})
	if err != nil {
		fmt.Fprintf(stderr, "%s history failed: %v\n", domain, err)
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
