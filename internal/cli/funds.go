package cli

import (
	"context"
	"fmt"
	"io"
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
	case "holdings":
		return runFundsHoldings(ctx, args[1:], stdout, stderr, getenv)
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
  holdings Fetch ETF holdings
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
