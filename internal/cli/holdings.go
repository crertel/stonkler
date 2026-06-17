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

type holdingsOptions struct {
	format outputFormat
	symbol string
	limit  int
}

func runFundsHoldings(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeFundsHoldingsHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeFundsHoldingsHelp(stdout)
		return 0
	}

	options, ok := parseHoldingsOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	holdings, err := client.ETFHoldings(ctx, options.symbol)
	if err != nil {
		fmt.Fprintf(stderr, "funds holdings failed: %v\n", err)
		return 1
	}
	if options.limit > 0 && len(holdings) > options.limit {
		holdings = holdings[:options.limit]
	}

	if err := writeETFHoldings(stdout, holdings, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func parseHoldingsOptions(args []string, stderr io.Writer) (holdingsOptions, bool) {
	options := holdingsOptions{
		format: outputTable,
		limit:  25,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--json":
			if options.format == outputCSV {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return holdingsOptions{}, false
			}
			options.format = outputJSON
		case "--csv":
			if options.format == outputJSON {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return holdingsOptions{}, false
			}
			options.format = outputCSV
		case "--limit":
			value, ok := nextFlagValue(args, &i, "--limit", stderr)
			if !ok {
				return holdingsOptions{}, false
			}
			limit, err := strconv.Atoi(value)
			if err != nil || limit < 0 {
				fmt.Fprintf(stderr, "invalid --limit value %q\n", value)
				return holdingsOptions{}, false
			}
			options.limit = limit
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return holdingsOptions{}, false
			}
			if options.symbol != "" {
				fmt.Fprintln(stderr, "funds holdings requires exactly one symbol")
				return holdingsOptions{}, false
			}
			options.symbol = arg
		}
	}

	if options.symbol == "" {
		fmt.Fprintln(stderr, "funds holdings requires exactly one symbol")
		return holdingsOptions{}, false
	}
	return options, true
}

func writeETFHoldings(w io.Writer, holdings []fmp.ETFHolding, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(holdings)
	case outputCSV:
		return writeETFHoldingsCSV(w, holdings)
	default:
		return writeETFHoldingsTable(w, holdings)
	}
}

func writeETFHoldingsTable(w io.Writer, holdings []fmp.ETFHolding) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ASSET\tNAME\tWEIGHT%\tSHARES\tMARKET VALUE\tUPDATED")
	for _, holding := range holdings {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s\t%s\t%s\n",
			holding.Asset,
			holding.Name,
			formatFloat(holding.WeightPercentage),
			formatFloat(holding.SharesNumber),
			formatFloat(holding.MarketValue),
			holding.Updated,
		)
	}
	return tw.Flush()
}

func writeETFHoldingsCSV(w io.Writer, holdings []fmp.ETFHolding) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"asset", "name", "isin", "cusip", "shares", "weight_percent", "market_value", "updated"}); err != nil {
		return err
	}
	for _, holding := range holdings {
		if err := cw.Write([]string{
			holding.Asset,
			holding.Name,
			holding.ISIN,
			holding.CUSIP,
			formatFloat(holding.SharesNumber),
			formatFloat(holding.WeightPercentage),
			formatFloat(holding.MarketValue),
			holding.Updated,
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
