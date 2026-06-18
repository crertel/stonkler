package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/crertel/stonkler/internal/fmp"
)

type portfolioQuoteUpdate struct {
	Timestamp string           `json:"timestamp"`
	Quotes    []quoteWithBasis `json:"quotes,omitempty"`
	Error     string           `json:"error,omitempty"`
}

func runPortfolio(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writePortfolioHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		writePortfolioHelp(stdout)
		return 0
	case "show":
		return runPortfolioShow(args[1:], stdout, stderr, getenv)
	case "quote":
		return runPortfolioQuote(ctx, args[1:], stdout, stderr, getenv)
	case "watch":
		return runPortfolioWatch(ctx, args[1:], stdout, stderr, getenv)
	default:
		fmt.Fprintf(stderr, "unknown portfolio command %q\n\n", args[0])
		writePortfolioHelp(stderr)
		return 2
	}
}

func writePortfolioHelp(w io.Writer) {
	fmt.Fprint(w, `Portfolio cost basis and market value views.

Usage:
  stonk portfolio <command> [flags]

Commands:
  show   Print configured portfolio lots and basis summaries
  quote  Fetch current quotes for configured portfolio symbols
  watch  Refresh current quotes for configured portfolio symbols

Flags:
  --basis <path>  Portfolio basis JSON file; defaults to STONK_PORTFOLIO_FILE
  --json          Write JSON output
  --csv           Write CSV output
`)
}

func runPortfolioShow(args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	options, ok := parseBasisOutputOptions(args, stderr)
	if !ok {
		return 2
	}
	if len(options.remaining) != 0 {
		fmt.Fprintln(stderr, "portfolio show does not accept symbols")
		return 2
	}
	book, ok := loadBasisPathOption(options.basisPath, stderr, getenv)
	if !ok {
		return 2
	}
	if err := writeBasisEntries(stdout, book.Entries(), options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func runPortfolioQuote(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	options, ok := parseBasisOutputOptions(args, stderr)
	if !ok {
		return 2
	}
	if len(options.remaining) != 0 {
		fmt.Fprintln(stderr, "portfolio quote does not accept symbols")
		return 2
	}
	book, ok := loadBasisPathOption(options.basisPath, stderr, getenv)
	if !ok {
		return 2
	}
	rows, code := fetchPortfolioQuotes(ctx, stderr, getenv, book)
	if code != 0 {
		return code
	}
	if err := writeQuotesWithBasis(stdout, rows, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

type portfolioWatchOptions struct {
	basisPath string
	interval  time.Duration
	count     int
	jsonl     bool
}

func runPortfolioWatch(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	options, ok := parsePortfolioWatchOptions(args, stderr)
	if !ok {
		return 2
	}
	book, ok := loadBasisPathOption(options.basisPath, stderr, getenv)
	if !ok {
		return 2
	}

	iteration := 0
	for {
		if options.count > 0 && iteration >= options.count {
			return 0
		}
		iteration++

		now := time.Now().Format(time.RFC3339)
		rows, code := fetchPortfolioQuotes(ctx, stderr, getenv, book)
		if options.jsonl {
			if err := writePortfolioQuoteJSONL(stdout, now, rows, code); err != nil {
				fmt.Fprintln(stderr, "failed to write output")
				return 1
			}
		} else {
			if iteration > 1 {
				fmt.Fprint(stdout, "\033[H\033[2J")
			}
			fmt.Fprintf(stdout, "Updated: %s\n\n", now)
			if code != 0 {
				fmt.Fprintln(stdout, "ERROR\tportfolio quote failed")
			} else if err := writeQuotesWithBasisTable(stdout, rows); err != nil {
				fmt.Fprintf(stderr, "failed to write output: %v\n", err)
				return 1
			}
		}
		if code != 0 {
			return code
		}
		if options.count > 0 && iteration >= options.count {
			return 0
		}

		select {
		case <-ctx.Done():
			return 130
		case <-time.After(options.interval):
		}
	}
}

func parsePortfolioWatchOptions(args []string, stderr io.Writer) (portfolioWatchOptions, bool) {
	options := portfolioWatchOptions{interval: 5 * time.Second}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--basis":
			value, ok := nextFlagValue(args, &i, "--basis", stderr)
			if !ok {
				return portfolioWatchOptions{}, false
			}
			options.basisPath = value
		case "--jsonl":
			options.jsonl = true
		case "--interval":
			value, ok := nextFlagValue(args, &i, "--interval", stderr)
			if !ok {
				return portfolioWatchOptions{}, false
			}
			interval, err := time.ParseDuration(value)
			if err != nil || interval <= 0 {
				fmt.Fprintf(stderr, "invalid --interval value %q\n", value)
				return portfolioWatchOptions{}, false
			}
			options.interval = interval
		case "--count":
			value, ok := nextFlagValue(args, &i, "--count", stderr)
			if !ok {
				return portfolioWatchOptions{}, false
			}
			count, err := parseNonNegativeInt(value)
			if err != nil {
				fmt.Fprintf(stderr, "invalid --count value %q\n", value)
				return portfolioWatchOptions{}, false
			}
			options.count = count
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return portfolioWatchOptions{}, false
			}
			fmt.Fprintln(stderr, "portfolio watch does not accept symbols")
			return portfolioWatchOptions{}, false
		}
	}
	return options, true
}

func parseNonNegativeInt(value string) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("invalid non-negative integer")
	}
	return parsed, nil
}

func loadBasisPathOption(flagPath string, stderr io.Writer, getenv getenvFunc) (*basisBook, bool) {
	path := resolveBasisPath(flagPath, getenv)
	if path == "" {
		fmt.Fprintln(stderr, "--basis is required or STONK_PORTFOLIO_FILE must be configured")
		return nil, false
	}
	book, err := loadBasisBook(path)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load basis file: %v\n", err)
		return nil, false
	}
	return book, true
}

func fetchPortfolioQuotes(ctx context.Context, stderr io.Writer, getenv getenvFunc, book *basisBook) ([]quoteWithBasis, int) {
	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return nil, 1
	}
	client := fmp.NewClient(apiKey, http.DefaultClient)

	var rows []quoteWithBasis
	for _, domain := range []string{"stocks", "funds", "crypto", "forex", "commodities", "indexes"} {
		symbols := symbolsForDomain(book, domain)
		if len(symbols) == 0 {
			continue
		}
		quotes, err := fetchDomainQuotes(ctx, client, domain, symbols)
		if err != nil {
			fmt.Fprintf(stderr, "portfolio quote failed for %s: %v\n", domain, err)
			return nil, 1
		}
		rows = append(rows, attachBasis(domain, quotes, book)...)
	}
	return rows, 0
}

func symbolsForDomain(book *basisBook, domain string) []string {
	var symbols []string
	for _, entry := range book.Entries() {
		if entry.Domain == domain {
			symbols = append(symbols, entry.Symbol)
		}
	}
	return symbols
}

func fetchDomainQuotes(ctx context.Context, client *fmp.Client, domain string, symbols []string) ([]fmp.Quote, error) {
	switch domain {
	case "stocks":
		return client.StockQuotes(ctx, symbols)
	case "indexes":
		return client.IndexQuotes(ctx, symbols)
	default:
		return client.BatchQuotes(ctx, symbols)
	}
}

func writePortfolioQuoteJSONL(w io.Writer, timestamp string, rows []quoteWithBasis, code int) error {
	update := portfolioQuoteUpdate{
		Timestamp: timestamp,
		Quotes:    rows,
	}
	if code != 0 {
		update.Quotes = nil
		update.Error = "portfolio quote failed"
	}
	return json.NewEncoder(w).Encode(update)
}
