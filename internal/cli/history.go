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

type historyOptions struct {
	format outputFormat
	symbol string
	from   string
	to     string
	limit  int
}

func runStocksHistory(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksHistoryHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksHistoryHelp(stdout)
		return 0
	}

	options, ok := parseHistoryOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	prices, err := client.StockHistory(ctx, fmp.StockHistoryRequest{
		Symbol: options.symbol,
		From:   options.from,
		To:     options.to,
	})
	if err != nil {
		fmt.Fprintf(stderr, "stocks history failed: %v\n", err)
		return 1
	}
	if options.limit > 0 && len(prices) > options.limit {
		prices = prices[:options.limit]
	}

	if err := writeStockHistory(stdout, prices, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func parseHistoryOptions(args []string, stderr io.Writer) (historyOptions, bool) {
	options := historyOptions{format: outputTable}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--json":
			if options.format == outputCSV {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return historyOptions{}, false
			}
			options.format = outputJSON
		case "--csv":
			if options.format == outputJSON {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return historyOptions{}, false
			}
			options.format = outputCSV
		case "--from":
			value, ok := nextFlagValue(args, &i, "--from", stderr)
			if !ok || !validDateFlag(value, "--from", stderr) {
				return historyOptions{}, false
			}
			options.from = value
		case "--to":
			value, ok := nextFlagValue(args, &i, "--to", stderr)
			if !ok || !validDateFlag(value, "--to", stderr) {
				return historyOptions{}, false
			}
			options.to = value
		case "--limit":
			value, ok := nextFlagValue(args, &i, "--limit", stderr)
			if !ok {
				return historyOptions{}, false
			}
			limit, err := strconv.Atoi(value)
			if err != nil || limit < 0 {
				fmt.Fprintf(stderr, "invalid --limit value %q\n", value)
				return historyOptions{}, false
			}
			options.limit = limit
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return historyOptions{}, false
			}
			if options.symbol != "" {
				fmt.Fprintln(stderr, "stocks history requires exactly one symbol")
				return historyOptions{}, false
			}
			options.symbol = arg
		}
	}

	if options.symbol == "" {
		fmt.Fprintln(stderr, "stocks history requires exactly one symbol")
		return historyOptions{}, false
	}
	return options, true
}

func nextFlagValue(args []string, index *int, name string, stderr io.Writer) (string, bool) {
	if *index+1 >= len(args) {
		fmt.Fprintf(stderr, "%s requires a value\n", name)
		return "", false
	}
	*index++
	return args[*index], true
}

func validDateFlag(value string, name string, stderr io.Writer) bool {
	if _, err := time.Parse("2006-01-02", value); err != nil {
		fmt.Fprintf(stderr, "invalid %s date %q; use YYYY-MM-DD\n", name, value)
		return false
	}
	return true
}

func writeStockHistory(w io.Writer, prices []fmp.StockPrice, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(prices)
	case outputCSV:
		return writeStockHistoryCSV(w, prices)
	default:
		return writeStockHistoryTable(w, prices)
	}
}

func writeStockHistoryTable(w io.Writer, prices []fmp.StockPrice) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "DATE\tSYMBOL\tOPEN\tHIGH\tLOW\tCLOSE\tCHANGE%\tVOLUME\tVWAP")
	for _, price := range prices {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			price.Date,
			price.Symbol,
			formatFloat(price.Open),
			formatFloat(price.High),
			formatFloat(price.Low),
			formatFloat(price.Close),
			formatFloat(price.ChangePercent),
			formatFloat(price.Volume),
			formatFloat(price.VWAP),
		)
	}
	return tw.Flush()
}

func writeStockHistoryCSV(w io.Writer, prices []fmp.StockPrice) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"date", "symbol", "open", "high", "low", "close", "change", "change_percent", "volume", "vwap"}); err != nil {
		return err
	}
	for _, price := range prices {
		if err := cw.Write([]string{
			price.Date,
			price.Symbol,
			formatFloat(price.Open),
			formatFloat(price.High),
			formatFloat(price.Low),
			formatFloat(price.Close),
			formatFloat(price.Change),
			formatFloat(price.ChangePercent),
			formatFloat(price.Volume),
			formatFloat(price.VWAP),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
