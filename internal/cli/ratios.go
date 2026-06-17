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

	"github.com/crertel/stonkler/internal/fmp"
)

type ratiosOptions struct {
	format outputFormat
	symbol string
}

type ratioColumn struct {
	header string
	field  string
}

var stockRatioColumns = []ratioColumn{
	{header: "GROSS MARGIN", field: "grossProfitMarginTTM"},
	{header: "OPERATING MARGIN", field: "operatingProfitMarginTTM"},
	{header: "NET MARGIN", field: "netProfitMarginTTM"},
	{header: "CURRENT", field: "currentRatioTTM"},
	{header: "QUICK", field: "quickRatioTTM"},
	{header: "P/E", field: "priceToEarningsRatioTTM"},
	{header: "P/B", field: "priceToBookRatioTTM"},
	{header: "P/S", field: "priceToSalesRatioTTM"},
	{header: "DEBT/EQUITY", field: "debtToEquityRatioTTM"},
	{header: "DIV YIELD", field: "dividendYieldTTM"},
}

func runStocksRatios(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksRatiosHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksRatiosHelp(stdout)
		return 0
	}

	options, ok := parseRatiosOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	ratios, err := client.StockRatiosTTM(ctx, options.symbol)
	if err != nil {
		fmt.Fprintf(stderr, "stocks ratios failed: %v\n", err)
		return 1
	}

	if err := writeStockRatios(stdout, ratios, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func parseRatiosOptions(args []string, stderr io.Writer) (ratiosOptions, bool) {
	options := ratiosOptions{format: outputTable}

	for _, arg := range args {
		switch arg {
		case "--ttm":
			continue
		case "--json":
			if options.format == outputCSV {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return ratiosOptions{}, false
			}
			options.format = outputJSON
		case "--csv":
			if options.format == outputJSON {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return ratiosOptions{}, false
			}
			options.format = outputCSV
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return ratiosOptions{}, false
			}
			if options.symbol != "" {
				fmt.Fprintln(stderr, "stocks ratios requires exactly one symbol")
				return ratiosOptions{}, false
			}
			options.symbol = arg
		}
	}

	if options.symbol == "" {
		fmt.Fprintln(stderr, "stocks ratios requires exactly one symbol")
		return ratiosOptions{}, false
	}
	return options, true
}

func writeStocksRatiosHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch trailing-twelve-month stock valuation and operating ratios.

Usage:
  stonk stocks ratios <symbol> [flags]

Flags:
  --ttm   Fetch trailing-twelve-month ratios
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func writeStockRatios(w io.Writer, ratios []fmp.StockRatioRow, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(ratios)
	case outputCSV:
		return writeStockRatiosCSV(w, ratios)
	default:
		return writeStockRatiosTable(w, ratios)
	}
}

func writeStockRatiosTable(w io.Writer, ratios []fmp.StockRatioRow) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprint(tw, "SYMBOL")
	for _, column := range stockRatioColumns {
		fmt.Fprintf(tw, "\t%s", column.header)
	}
	fmt.Fprintln(tw)

	for _, row := range ratios {
		fmt.Fprintf(tw, "%s", ratioValue(row, "symbol"))
		for _, column := range stockRatioColumns {
			fmt.Fprintf(tw, "\t%s", ratioValue(row, column.field))
		}
		fmt.Fprintln(tw)
	}
	return tw.Flush()
}

func writeStockRatiosCSV(w io.Writer, ratios []fmp.StockRatioRow) error {
	cw := csv.NewWriter(w)
	header := []string{"symbol"}
	for _, column := range stockRatioColumns {
		header = append(header, column.field)
	}
	if err := cw.Write(header); err != nil {
		return err
	}

	for _, ratio := range ratios {
		record := []string{ratioValue(ratio, "symbol")}
		for _, column := range stockRatioColumns {
			record = append(record, ratioValue(ratio, column.field))
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func ratioValue(row fmp.StockRatioRow, field string) string {
	value, ok := row[field]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return formatFloat(typed)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return fmt.Sprint(typed)
	}
}
