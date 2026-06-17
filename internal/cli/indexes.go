package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/crertel/stonkler/internal/fmp"
)

func runIndexes(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeIndexesHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeIndexesHelp(stdout)
		return 0
	case "history":
		return runDomainHistory(ctx, args[1:], stdout, stderr, getenv, "indexes", writeIndexesHistoryHelp, normalizeIndexSymbol)
	case "quote", "quotes":
		return runDomainQuote(ctx, args[1:], stdout, stderr, getenv, "indexes", writeIndexesQuoteHelp, indexQuotes)
	default:
		fmt.Fprintf(stderr, "unknown indexes command %q\n\n", args[0])
		writeIndexesHelp(stderr)
		return 2
	}
}

func indexQuotes(ctx context.Context, client *fmp.Client, symbols []string) ([]fmp.Quote, error) {
	return client.IndexQuotes(ctx, symbols)
}

func writeIndexesHelp(w io.Writer) {
	fmt.Fprint(w, `Index market data.

Usage:
  stonk indexes <command> [flags]

Commands:
  history Fetch historical end-of-day index prices
  quote   Fetch one or more index quotes
  quotes  Alias for quote
`)
}

func writeIndexesQuoteHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch one or more index quotes.

Usage:
  stonk indexes quote <symbol> [symbol...] [flags]
  stonk indexes quotes <symbol> [symbol...] [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func writeIndexesHistoryHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch historical end-of-day index prices.

Usage:
  stonk indexes history <symbol> [flags]

Flags:
  --from <date>  Start date in YYYY-MM-DD format
  --to <date>    End date in YYYY-MM-DD format
  --limit <n>    Maximum rows to print
  --json         Write JSON output
  --csv          Write CSV output
`)
}

func normalizeIndexSymbol(symbol string) string {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" || strings.HasPrefix(symbol, "^") {
		return symbol
	}
	return "^" + symbol
}
