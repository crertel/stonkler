package cli

import (
	"context"
	"fmt"
	"io"
	"os"
)

const version = "dev"

// Run executes the stonk command-line interface.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		writeRootHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeRootHelp(stdout)
		return 0
	case "version", "--version":
		fmt.Fprintf(stdout, "stonk %s\n", version)
		return 0
	case "config":
		return runConfig(args[1:], stdout, stderr, os.Getenv)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		writeRootHelp(stderr)
		return 2
	}
}

func writeRootHelp(w io.Writer) {
	fmt.Fprint(w, `stonk is a domain-first financial data CLI.

Usage:
  stonk <command> [flags]

Commands:
  stocks      Stock quotes, history, fundamentals, and watch views
  funds       ETF and mutual fund data
  crypto      Cryptocurrency market data
  forex       Foreign exchange market data
  commodities Commodity market data
  indexes     Index market data
  search      Discover symbols and securities
  get         Workflow-oriented shortcuts
  config      Configuration and provider diagnostics
  version     Print version information

Use "stonk <command> --help" for command-specific help.
`)
}
