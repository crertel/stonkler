package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/crertel/stonkler/internal/fmp"
)

func runFunds(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeFundsHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeFundsHelp(stdout)
		return 0
	case "exposure":
		return runFundsExposure(ctx, args[1:], stdout, stderr, getenv)
	case "country-weightings":
		return runFundsCountryWeightings(ctx, args[1:], stdout, stderr, getenv)
	case "holdings":
		return runFundsHoldings(ctx, args[1:], stdout, stderr, getenv)
	case "info":
		return runFundsInfo(ctx, args[1:], stdout, stderr, getenv)
	case "sector-weightings":
		return runFundsSectorWeightings(ctx, args[1:], stdout, stderr, getenv)
	case "watch":
		return runFundsWatch(ctx, args[1:], stdout, stderr, getenv)
	default:
		fmt.Fprintf(stderr, "unknown funds command %q\n\n", args[0])
		writeFundsHelp(stderr)
		return 2
	}
}

func writeFundsHelp(w io.Writer) {
	fmt.Fprint(w, `ETF and mutual fund data.

Usage:
  stonk funds <command> [flags]

Commands:
  country-weightings Fetch ETF or fund country allocation weights
  exposure          Fetch ETF or fund exposure to an asset
  holdings          Fetch ETF holdings
  info              Fetch ETF or fund profile information
  sector-weightings Fetch ETF or fund sector allocation weights
  watch             Refresh ETF or fund quotes in a terminal view
`)
}

func writeFundsHoldingsHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch ETF holdings.

Usage:
  stonk funds holdings <symbol> [flags]

Flags:
  --limit <n>  Maximum holdings to print
  --json       Write JSON output
  --csv        Write CSV output
`)
}

func writeFundsInfoHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch ETF or fund profile information.

Usage:
  stonk funds info <symbol> [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func runFundsWatch(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	return runQuoteWatchCommand(ctx, args, stdout, stderr, getenv, "funds", writeFundsWatchHelp, fundQuotes)
}

func fundQuotes(ctx context.Context, client *fmp.Client, symbols []string) ([]fmp.Quote, error) {
	return client.BatchQuotes(ctx, symbols)
}

func writeFundsWatchHelp(w io.Writer) {
	fmt.Fprint(w, `Refresh ETF or fund quotes in a terminal view.

Usage:
  stonk funds watch <symbol> [symbol...] [flags]

Flags:
  --interval <duration>  Refresh interval, such as 5s or 1m
  --count <n>            Number of refreshes before exiting
  --jsonl                Write newline-delimited JSON updates
`)
}
