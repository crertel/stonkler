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
	case "quote", "quotes":
		return runDomainQuote(ctx, args[1:], stdout, stderr, getenv, "crypto", writeCryptoQuoteHelp, cryptoQuotes)
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
  quote   Fetch one or more cryptocurrency quotes
  quotes  Alias for quote
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
