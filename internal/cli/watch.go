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

type watchOptions struct {
	interval time.Duration
	count    int
	jsonl    bool
	symbols  []string
}

type stockWatchUpdate struct {
	Timestamp string           `json:"timestamp"`
	Quotes    []fmp.StockQuote `json:"quotes,omitempty"`
	Error     string           `json:"error,omitempty"`
}

func runStocksWatch(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksWatchHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksWatchHelp(stdout)
		return 0
	}

	options, ok := parseWatchOptions(args, stderr)
	if !ok {
		return 2
	}
	if len(options.symbols) == 0 {
		fmt.Fprintln(stderr, "stocks watch requires at least one symbol")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	return runQuoteWatchLoop(ctx, stdout, stderr, client, options, func(ctx context.Context, client *fmp.Client, symbols []string) ([]fmp.Quote, error) {
		return client.StockQuotes(ctx, symbols)
	})
}

func parseWatchOptions(args []string, stderr io.Writer) (watchOptions, bool) {
	options := watchOptions{
		interval: 5 * time.Second,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--jsonl":
			options.jsonl = true
		case "--interval":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "--interval requires a duration")
				return watchOptions{}, false
			}
			i++
			interval, err := time.ParseDuration(args[i])
			if err != nil || interval <= 0 {
				fmt.Fprintf(stderr, "invalid --interval value %q\n", args[i])
				return watchOptions{}, false
			}
			options.interval = interval
		case "--count":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "--count requires a number")
				return watchOptions{}, false
			}
			i++
			count, err := strconv.Atoi(args[i])
			if err != nil || count < 0 {
				fmt.Fprintf(stderr, "invalid --count value %q\n", args[i])
				return watchOptions{}, false
			}
			options.count = count
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return watchOptions{}, false
			}
			options.symbols = append(options.symbols, arg)
		}
	}

	return options, true
}

func runQuoteWatchLoop(ctx context.Context, stdout, stderr io.Writer, client *fmp.Client, options watchOptions, fetch quoteFetcher) int {
	iteration := 0
	for {
		if options.count > 0 && iteration >= options.count {
			return 0
		}
		iteration++

		now := time.Now().Format(time.RFC3339)
		quotes, err := fetch(ctx, client, options.symbols)
		if options.jsonl {
			if writeStockWatchJSONL(stdout, now, quotes, err) != nil {
				fmt.Fprintln(stderr, "failed to write output")
				return 1
			}
		} else {
			if iteration > 1 {
				fmt.Fprint(stdout, "\033[H\033[2J")
			}
			fmt.Fprintf(stdout, "Updated: %s\n\n", now)
			if err != nil {
				fmt.Fprintf(stdout, "ERROR\t%s\n", err)
			} else if err := writeStockQuotesTable(stdout, quotes); err != nil {
				fmt.Fprintf(stderr, "failed to write output: %v\n", err)
				return 1
			}
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

func writeStockWatchJSONL(w io.Writer, timestamp string, quotes []fmp.StockQuote, err error) error {
	update := stockWatchUpdate{
		Timestamp: timestamp,
		Quotes:    quotes,
	}
	if err != nil {
		update.Error = err.Error()
		update.Quotes = nil
	}

	return json.NewEncoder(w).Encode(update)
}
