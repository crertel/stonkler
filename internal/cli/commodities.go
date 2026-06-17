package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/crertel/stonkler/internal/fmp"
)

func runCommodities(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeCommoditiesHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeCommoditiesHelp(stdout)
		return 0
	case "quote", "quotes":
		return runDomainQuote(ctx, args[1:], stdout, stderr, getenv, "commodities", writeCommoditiesQuoteHelp, commodityQuotes)
	default:
		fmt.Fprintf(stderr, "unknown commodities command %q\n\n", args[0])
		writeCommoditiesHelp(stderr)
		return 2
	}
}

func commodityQuotes(ctx context.Context, client *fmp.Client, symbols []string) ([]fmp.Quote, error) {
	return client.BatchQuotes(ctx, symbols)
}

func writeCommoditiesHelp(w io.Writer) {
	fmt.Fprint(w, `Commodity market data.

Usage:
  stonk commodities <command> [flags]

Commands:
  quote   Fetch one or more commodity quotes
  quotes  Alias for quote
`)
}

func writeCommoditiesQuoteHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch one or more commodity quotes.

Usage:
  stonk commodities quote <symbol> [symbol...] [flags]
  stonk commodities quotes <symbol> [symbol...] [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}
