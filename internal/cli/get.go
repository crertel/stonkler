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
	case "quote", "quotes":
		return runStocksQuote(ctx, args[1:], stdout, stderr, getenv)
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
  quote   Fetch one or more quotes, inferring the stock domain for now
  quotes  Alias for quote
`)
}
