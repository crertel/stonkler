package cli

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/tabwriter"

	"github.com/crertel/stonkler/internal/fmp"
)

type metricsOptions struct {
	format outputFormat
	symbol string
}

var stockMetricColumns = []ratioColumn{
	{header: "MARKET CAP", field: "marketCap"},
	{header: "ENTERPRISE VALUE", field: "enterpriseValueTTM"},
	{header: "EV/SALES", field: "evToSalesTTM"},
	{header: "EV/EBITDA", field: "evToEBITDATTM"},
	{header: "NET DEBT/EBITDA", field: "netDebtToEBITDATTM"},
	{header: "ROA", field: "returnOnAssetsTTM"},
	{header: "ROE", field: "returnOnEquityTTM"},
	{header: "ROIC", field: "returnOnInvestedCapitalTTM"},
	{header: "FCF YIELD", field: "freeCashFlowYieldTTM"},
	{header: "CASH CONV CYCLE", field: "cashConversionCycleTTM"},
}

func runStocksMetrics(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksMetricsHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksMetricsHelp(stdout)
		return 0
	}

	options, ok := parseMetricsOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	metrics, err := client.StockKeyMetricsTTM(ctx, options.symbol)
	if err != nil {
		fmt.Fprintf(stderr, "stocks metrics failed: %v\n", err)
		return 1
	}

	if err := writeStockMetrics(stdout, metrics, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func parseMetricsOptions(args []string, stderr io.Writer) (metricsOptions, bool) {
	options := metricsOptions{format: outputTable}

	for _, arg := range args {
		switch arg {
		case "--ttm":
			continue
		case "--json":
			if options.format == outputCSV {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return metricsOptions{}, false
			}
			options.format = outputJSON
		case "--csv":
			if options.format == outputJSON {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return metricsOptions{}, false
			}
			options.format = outputCSV
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return metricsOptions{}, false
			}
			if options.symbol != "" {
				fmt.Fprintln(stderr, "stocks metrics requires exactly one symbol")
				return metricsOptions{}, false
			}
			options.symbol = arg
		}
	}

	if options.symbol == "" {
		fmt.Fprintln(stderr, "stocks metrics requires exactly one symbol")
		return metricsOptions{}, false
	}
	return options, true
}

func writeStocksMetricsHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch trailing-twelve-month stock key metrics.

Usage:
  stonk stocks metrics <symbol> [flags]

Flags:
  --ttm   Fetch trailing-twelve-month metrics
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func writeStockMetrics(w io.Writer, metrics []fmp.StockMetricRow, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(metrics)
	case outputCSV:
		return writeStockMetricsCSV(w, metrics)
	default:
		return writeStockMetricsTable(w, metrics)
	}
}

func writeStockMetricsTable(w io.Writer, metrics []fmp.StockMetricRow) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprint(tw, "SYMBOL")
	for _, column := range stockMetricColumns {
		fmt.Fprintf(tw, "\t%s", column.header)
	}
	fmt.Fprintln(tw)

	for _, row := range metrics {
		fmt.Fprintf(tw, "%s", ratioValue(fmp.StockRatioRow(row), "symbol"))
		for _, column := range stockMetricColumns {
			fmt.Fprintf(tw, "\t%s", ratioValue(fmp.StockRatioRow(row), column.field))
		}
		fmt.Fprintln(tw)
	}
	return tw.Flush()
}

func writeStockMetricsCSV(w io.Writer, metrics []fmp.StockMetricRow) error {
	cw := csv.NewWriter(w)
	header := []string{"symbol"}
	for _, column := range stockMetricColumns {
		header = append(header, column.field)
	}
	if err := cw.Write(header); err != nil {
		return err
	}

	for _, metric := range metrics {
		row := fmp.StockRatioRow(metric)
		record := []string{ratioValue(row, "symbol")}
		for _, column := range stockMetricColumns {
			record = append(record, ratioValue(row, column.field))
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
