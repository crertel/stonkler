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
	case "history":
		return runDomainHistory(ctx, args[1:], stdout, stderr, getenv, "forex", writeForexHistoryHelp, nil)
	case "quote", "quotes":
		return runDomainQuote(ctx, args[1:], stdout, stderr, getenv, "forex", writeForexQuoteHelp, forexQuotes)
	case "watch":
		return runForexWatch(ctx, args[1:], stdout, stderr, getenv)
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
  history Fetch historical end-of-day forex prices
  quote   Fetch one or more forex quotes
  quotes  Alias for quote
  watch   Refresh forex quotes in a terminal view
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

func writeForexHistoryHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch historical end-of-day forex prices.

Usage:
  stonk forex history <symbol> [flags]

Flags:
  --from <date>  Start date in YYYY-MM-DD format
  --to <date>    End date in YYYY-MM-DD format
  --limit <n>    Maximum rows to print
  --json         Write JSON output
  --csv          Write CSV output
`)
}

func runForexWatch(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	return runQuoteWatchCommand(ctx, args, stdout, stderr, getenv, "forex", writeForexWatchHelp, forexQuotes)
}

func writeForexWatchHelp(w io.Writer) {
	fmt.Fprint(w, `Refresh forex quotes in a terminal view.

Usage:
  stonk forex watch <symbol> [symbol...] [flags]

Flags:
  --interval <duration>  Refresh interval, such as 5s or 1m
  --count <n>            Number of refreshes before exiting
  --sort <field>         Sort by symbol, price, change, change-percent, or volume
  --jsonl                Write newline-delimited JSON updates
`)
}
