package cli

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/crertel/stonkler/internal/fmp"
)

func runStocks(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeStocksHelp(stdout)
		return 0
	case "quote", "quotes":
		return runStocksQuote(ctx, args[1:], stdout, stderr, getenv)
	case "watch":
		return runStocksWatch(ctx, args[1:], stdout, stderr, getenv)
	default:
		fmt.Fprintf(stderr, "unknown stocks command %q\n\n", args[0])
		writeStocksHelp(stderr)
		return 2
	}
}

func runStocksQuote(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksQuoteHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksQuoteHelp(stdout)
		return 0
	}

	format, symbols, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(symbols) == 0 {
		fmt.Fprintln(stderr, "stocks quote requires at least one symbol")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	quotes, err := client.StockQuotes(ctx, symbols)
	if err != nil {
		fmt.Fprintf(stderr, "stocks quote failed: %v\n", err)
		return 1
	}

	if err := writeStockQuotes(stdout, quotes, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeStocksHelp(w io.Writer) {
	fmt.Fprint(w, `Stock quotes, history, fundamentals, and watch views.

Usage:
  stonk stocks <command> [flags]

Commands:
  quote   Fetch one or more stock quotes
  quotes  Alias for quote
  watch   Refresh stock quotes in a terminal view
`)
}

func writeStocksQuoteHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch one or more stock quotes.

Usage:
  stonk stocks quote <symbol> [symbol...] [flags]
  stonk stocks quotes <symbol> [symbol...] [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func writeStocksWatchHelp(w io.Writer) {
	fmt.Fprint(w, `Refresh stock quotes in a terminal view.

Usage:
  stonk stocks watch <symbol> [symbol...] [flags]

Flags:
  --interval <duration>  Refresh interval, such as 5s or 1m
  --count <n>            Number of refreshes before exiting
  --jsonl                Write newline-delimited JSON updates
`)
}

func writeStockQuotes(w io.Writer, quotes []fmp.StockQuote, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(quotes)
	case outputCSV:
		return writeStockQuotesCSV(w, quotes)
	default:
		return writeStockQuotesTable(w, quotes)
	}
}

func writeStockQuotesTable(w io.Writer, quotes []fmp.StockQuote) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SYMBOL\tNAME\tPRICE\tCHANGE\tCHANGE%\tVOLUME\tMARKET CAP\tUPDATED")
	for _, quote := range quotes {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			quote.Symbol,
			quote.Name,
			formatFloat(quote.Price),
			formatFloat(quote.Change),
			formatFloat(quote.ChangePercentage),
			formatFloat(quote.Volume),
			formatFloat(quote.MarketCap),
			formatUnixTimestamp(quote.Timestamp),
		)
	}
	return tw.Flush()
}

func writeStockQuotesCSV(w io.Writer, quotes []fmp.StockQuote) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"symbol", "name", "price", "change", "change_percent", "volume", "market_cap", "timestamp"}); err != nil {
		return err
	}
	for _, quote := range quotes {
		if err := cw.Write([]string{
			quote.Symbol,
			quote.Name,
			formatFloat(quote.Price),
			formatFloat(quote.Change),
			formatFloat(quote.ChangePercentage),
			formatFloat(quote.Volume),
			formatFloat(quote.MarketCap),
			strconv.FormatInt(quote.Timestamp, 10),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func formatUnixTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(timestamp, 0).Format(time.RFC3339)
}
