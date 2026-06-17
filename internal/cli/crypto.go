package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/crertel/stonkler/internal/fmp"
)

func runCrypto(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeCryptoHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeCryptoHelp(stdout)
		return 0
	case "history":
		return runDomainHistory(ctx, args[1:], stdout, stderr, getenv, "crypto", writeCryptoHistoryHelp, nil)
	case "quote", "quotes":
		return runDomainQuote(ctx, args[1:], stdout, stderr, getenv, "crypto", writeCryptoQuoteHelp, cryptoQuotes)
	case "watch":
		return runCryptoWatch(ctx, args[1:], stdout, stderr, getenv)
	default:
		fmt.Fprintf(stderr, "unknown crypto command %q\n\n", args[0])
		writeCryptoHelp(stderr)
		return 2
	}
}

func cryptoQuotes(ctx context.Context, client *fmp.Client, symbols []string) ([]fmp.Quote, error) {
	return client.BatchQuotes(ctx, symbols)
}

func writeCryptoHelp(w io.Writer) {
	fmt.Fprint(w, `Cryptocurrency market data.

Usage:
  stonk crypto <command> [flags]

Commands:
  history Fetch historical end-of-day crypto prices
  quote   Fetch one or more cryptocurrency quotes
  quotes  Alias for quote
  watch   Refresh cryptocurrency quotes in a terminal view
`)
}

func writeCryptoQuoteHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch one or more cryptocurrency quotes.

Usage:
  stonk crypto quote <symbol> [symbol...] [flags]
  stonk crypto quotes <symbol> [symbol...] [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func writeCryptoHistoryHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch historical end-of-day crypto prices.

Usage:
  stonk crypto history <symbol> [flags]

Flags:
  --from <date>  Start date in YYYY-MM-DD format
  --to <date>    End date in YYYY-MM-DD format
  --limit <n>    Maximum rows to print
  --json         Write JSON output
  --csv          Write CSV output
`)
}

func runCryptoWatch(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	return runQuoteWatchCommand(ctx, args, stdout, stderr, getenv, "crypto", writeCryptoWatchHelp, cryptoQuotes)
}

func writeCryptoWatchHelp(w io.Writer) {
	fmt.Fprint(w, `Refresh cryptocurrency quotes in a terminal view.

Usage:
  stonk crypto watch <symbol> [symbol...] [flags]

Flags:
  --interval <duration>  Refresh interval, such as 5s or 1m
  --count <n>            Number of refreshes before exiting
  --sort <field>         Sort by symbol, price, change, change-percent, or volume
  --jsonl                Write newline-delimited JSON updates
`)
}
