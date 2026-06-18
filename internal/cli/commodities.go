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
	case "history":
		return runDomainHistory(ctx, args[1:], stdout, stderr, getenv, "commodities", writeCommoditiesHistoryHelp, nil)
	case "quote", "quotes":
		return runDomainQuote(ctx, args[1:], stdout, stderr, getenv, "commodities", writeCommoditiesQuoteHelp, commodityQuotes)
	case "watch":
		return runCommoditiesWatch(ctx, args[1:], stdout, stderr, getenv)
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
  history Fetch historical end-of-day commodity prices
  quote   Fetch one or more commodity quotes
  quotes  Alias for quote
  watch   Refresh commodity quotes in a terminal view
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

func writeCommoditiesHistoryHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch historical end-of-day commodity prices.

Usage:
  stonk commodities history <symbol> [flags]

Flags:
  --from <date>  Start date in YYYY-MM-DD format
  --to <date>    End date in YYYY-MM-DD format
  --limit <n>    Maximum rows to print
  --json         Write JSON output
  --csv          Write CSV output
`)
}

func runCommoditiesWatch(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	return runQuoteWatchCommand(ctx, args, stdout, stderr, getenv, "commodities", writeCommoditiesWatchHelp, commodityQuotes)
}

func writeCommoditiesWatchHelp(w io.Writer) {
	fmt.Fprint(w, `Refresh commodity quotes in a terminal view.

Usage:
  stonk commodities watch <symbol> [symbol...] [flags]

Flags:
  --interval <duration>  Refresh interval, such as 5s or 1m
  --count <n>            Number of refreshes before exiting
  --sort <field>         Sort by symbol, price, change, change-percent, or volume
  --fields <list>        Comma-separated fields to show
  --jsonl                Write newline-delimited JSON updates
`)
}
