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

type insidersOptions struct {
	format outputFormat
	symbol string
	limit  int
}

func runStocksInsiders(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksInsidersHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksInsidersHelp(stdout)
		return 0
	}

	options, ok := parseInsidersOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	trades, err := client.InsiderTrades(ctx, options.symbol, options.limit)
	if err != nil {
		fmt.Fprintf(stderr, "stocks insiders failed: %v\n", err)
		return 1
	}

	if err := writeInsiderTrades(stdout, trades, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func parseInsidersOptions(args []string, stderr io.Writer) (insidersOptions, bool) {
	options := insidersOptions{
		format: outputTable,
		limit:  25,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--json":
			if options.format == outputCSV {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return insidersOptions{}, false
			}
			options.format = outputJSON
		case "--csv":
			if options.format == outputJSON {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return insidersOptions{}, false
			}
			options.format = outputCSV
		case "--limit":
			value, ok := nextFlagValue(args, &i, "--limit", stderr)
			if !ok {
				return insidersOptions{}, false
			}
			limit, err := strconv.Atoi(value)
			if err != nil || limit < 0 {
				fmt.Fprintf(stderr, "invalid --limit value %q\n", value)
				return insidersOptions{}, false
			}
			options.limit = limit
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return insidersOptions{}, false
			}
			if options.symbol != "" {
				fmt.Fprintln(stderr, "stocks insiders requires exactly one symbol")
				return insidersOptions{}, false
			}
			options.symbol = arg
		}
	}

	if options.symbol == "" {
		fmt.Fprintln(stderr, "stocks insiders requires exactly one symbol")
		return insidersOptions{}, false
	}
	return options, true
}

func writeStocksInsidersHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch insider transactions.

Usage:
  stonk stocks insiders <symbol> [flags]

Flags:
  --limit <n>  Maximum transactions to request
  --json       Write JSON output
  --csv        Write CSV output
`)
}

func writeInsiderTrades(w io.Writer, trades []fmp.InsiderTrade, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(trades)
	case outputCSV:
		return writeInsiderTradesCSV(w, trades)
	default:
		return writeInsiderTradesTable(w, trades)
	}
}

func writeInsiderTradesTable(w io.Writer, trades []fmp.InsiderTrade) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SYMBOL\tFILING\tTRANSACTION\tREPORTING NAME\tTYPE\tACTION\tSHARES\tPRICE\tOWNED\tFORM")
	for _, trade := range trades {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			trade.Symbol,
			trade.FilingDate,
			trade.TransactionDate,
			trade.ReportingName,
			trade.TransactionType,
			trade.AcquisitionOrDisposition,
			formatFloat(trade.SecuritiesTransacted),
			formatFloat(trade.Price),
			formatFloat(trade.SecuritiesOwned),
			trade.FormType,
		)
	}
	return tw.Flush()
}

func writeInsiderTradesCSV(w io.Writer, trades []fmp.InsiderTrade) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"symbol", "filing_date", "transaction_date", "reporting_name", "reporting_cik", "company_cik", "owner_type", "transaction_type", "action", "direct_or_indirect", "form_type", "shares", "price", "securities_owned", "security_name", "url"}); err != nil {
		return err
	}
	for _, trade := range trades {
		if err := cw.Write([]string{
			trade.Symbol,
			trade.FilingDate,
			trade.TransactionDate,
			trade.ReportingName,
			trade.ReportingCIK,
			trade.CompanyCIK,
			trade.TypeOfOwner,
			trade.TransactionType,
			trade.AcquisitionOrDisposition,
			trade.DirectOrIndirect,
			trade.FormType,
			formatFloat(trade.SecuritiesTransacted),
			formatFloat(trade.Price),
			formatFloat(trade.SecuritiesOwned),
			trade.SecurityName,
			trade.URL,
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
