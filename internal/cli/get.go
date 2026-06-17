package cli

import (
	"context"
	"fmt"
	"io"
)

func runGet(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeGetHelp(stdout)
		return 0
	case "company", "profile":
		return runStocksProfile(ctx, args[1:], stdout, stderr, getenv)
	case "etf", "fund", "fund-info":
		return runFundsInfo(ctx, args[1:], stdout, stderr, getenv)
	case "history":
		return runStocksHistory(ctx, args[1:], stdout, stderr, getenv)
	case "holdings":
		return runFundsHoldings(ctx, args[1:], stdout, stderr, getenv)
	case "quote", "quotes":
		return runStocksQuote(ctx, args[1:], stdout, stderr, getenv)
	case "statement", "statements":
		return runStocksStatements(ctx, args[1:], stdout, stderr, getenv)
	default:
		fmt.Fprintf(stderr, "unknown get command %q\n\n", args[0])
		writeGetHelp(stderr)
		return 2
	}
}

func writeGetHelp(w io.Writer) {
	fmt.Fprint(w, `Workflow-oriented shortcuts.

Usage:
  stonk get <command> [flags]

Commands:
  company Fetch company profile data, inferring the stock domain for now
  fund    Fetch ETF or fund profile information
  history Fetch historical prices, inferring the stock domain for now
  holdings Fetch ETF holdings, inferring the funds domain for now
  profile Alias for company
  quote   Fetch one or more quotes, inferring the stock domain for now
  quotes  Alias for quote
  statements Fetch financial statements, inferring the stock domain for now
`)
}
