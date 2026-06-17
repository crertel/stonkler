package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/crertel/stonkler/internal/fmp"
)

func runForex(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeForexHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeForexHelp(stdout)
		return 0
	case "quote", "quotes":
		return runDomainQuote(ctx, args[1:], stdout, stderr, getenv, "forex", writeForexQuoteHelp, forexQuotes)
	default:
		fmt.Fprintf(stderr, "unknown forex command %q\n\n", args[0])
		writeForexHelp(stderr)
		return 2
	}
}

func forexQuotes(ctx context.Context, client *fmp.Client, symbols []string) ([]fmp.Quote, error) {
	return client.BatchQuotes(ctx, symbols)
}

func writeForexHelp(w io.Writer) {
	fmt.Fprint(w, `Foreign exchange market data.

Usage:
  stonk forex <command> [flags]

Commands:
  quote   Fetch one or more forex quotes
  quotes  Alias for quote
`)
}

func writeForexQuoteHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch one or more forex quotes.

Usage:
  stonk forex quote <symbol> [symbol...] [flags]
  stonk forex quotes <symbol> [symbol...] [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}
