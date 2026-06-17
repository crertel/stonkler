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

type secOptions struct {
	format outputFormat
	symbol string
	limit  int
}

func runStocksSEC(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksSECHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksSECHelp(stdout)
		return 0
	}

	options, ok := parseSECOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	filings, err := client.SECFilings(ctx, options.symbol, options.limit)
	if err != nil {
		fmt.Fprintf(stderr, "stocks sec failed: %v\n", err)
		return 1
	}

	if err := writeSECFilings(stdout, filings, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func parseSECOptions(args []string, stderr io.Writer) (secOptions, bool) {
	options := secOptions{
		format: outputTable,
		limit:  25,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--json":
			if options.format == outputCSV {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return secOptions{}, false
			}
			options.format = outputJSON
		case "--csv":
			if options.format == outputJSON {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return secOptions{}, false
			}
			options.format = outputCSV
		case "--limit":
			value, ok := nextFlagValue(args, &i, "--limit", stderr)
			if !ok {
				return secOptions{}, false
			}
			limit, err := strconv.Atoi(value)
			if err != nil || limit < 0 {
				fmt.Fprintf(stderr, "invalid --limit value %q\n", value)
				return secOptions{}, false
			}
			options.limit = limit
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return secOptions{}, false
			}
			if options.symbol != "" {
				fmt.Fprintln(stderr, "stocks sec requires exactly one symbol")
				return secOptions{}, false
			}
			options.symbol = arg
		}
	}

	if options.symbol == "" {
		fmt.Fprintln(stderr, "stocks sec requires exactly one symbol")
		return secOptions{}, false
	}
	return options, true
}

func writeStocksSECHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch SEC filings.

Usage:
  stonk stocks sec <symbol> [flags]

Flags:
  --limit <n>  Maximum filings to request
  --json       Write JSON output
  --csv        Write CSV output
`)
}

func writeSECFilings(w io.Writer, filings []fmp.SECFiling, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(filings)
	case outputCSV:
		return writeSECFilingsCSV(w, filings)
	default:
		return writeSECFilingsTable(w, filings)
	}
}

func writeSECFilingsTable(w io.Writer, filings []fmp.SECFiling) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SYMBOL\tTYPE\tFILING DATE\tACCEPTED\tCIK\tLINK")
	for _, filing := range filings {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", filing.Symbol, filing.Type, filing.FilingDate, filing.AcceptedDate, filing.CIK, filing.Link)
	}
	return tw.Flush()
}

func writeSECFilingsCSV(w io.Writer, filings []fmp.SECFiling) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"symbol", "type", "filing_date", "accepted_date", "cik", "link", "final_link"}); err != nil {
		return err
	}
	for _, filing := range filings {
		if err := cw.Write([]string{
			filing.Symbol,
			filing.Type,
			filing.FilingDate,
			filing.AcceptedDate,
			filing.CIK,
			filing.Link,
			filing.FinalLink,
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
