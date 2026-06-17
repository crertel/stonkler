package cli

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/crertel/stonkler/internal/fmp"
)

type searchScreenerOptions struct {
	format outputFormat
	fmp    fmp.ScreenerOptions
}

func runSearchScreener(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		writeSearchScreenerHelp(stdout)
		return 0
	}

	options, ok := parseSearchScreenerOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	results, err := client.CompanyScreener(ctx, options.fmp)
	if err != nil {
		fmt.Fprintf(stderr, "search screener failed: %v\n", err)
		return 1
	}

	if err := writeScreenerResults(stdout, results, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func parseSearchScreenerOptions(args []string, stderr io.Writer) (searchScreenerOptions, bool) {
	options := searchScreenerOptions{
		format: outputTable,
		fmp:    fmp.ScreenerOptions{Limit: 25},
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--json":
			if options.format == outputCSV {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return searchScreenerOptions{}, false
			}
			options.format = outputJSON
		case "--csv":
			if options.format == outputJSON {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return searchScreenerOptions{}, false
			}
			options.format = outputCSV
		case "--sector":
			value, ok := nextFlagValue(args, &i, "--sector", stderr)
			if !ok {
				return searchScreenerOptions{}, false
			}
			options.fmp.Sector = value
		case "--country":
			value, ok := nextFlagValue(args, &i, "--country", stderr)
			if !ok {
				return searchScreenerOptions{}, false
			}
			options.fmp.Country = strings.ToUpper(value)
		case "--market-cap-min":
			value, ok := nextFlagValue(args, &i, "--market-cap-min", stderr)
			if !ok {
				return searchScreenerOptions{}, false
			}
			amount, err := parseMarketCapAmount(value)
			if err != nil {
				fmt.Fprintf(stderr, "invalid --market-cap-min value %q\n", value)
				return searchScreenerOptions{}, false
			}
			options.fmp.MarketCapMin = amount
		case "--limit":
			value, ok := nextFlagValue(args, &i, "--limit", stderr)
			if !ok {
				return searchScreenerOptions{}, false
			}
			limit, err := strconv.Atoi(value)
			if err != nil || limit < 0 {
				fmt.Fprintf(stderr, "invalid --limit value %q\n", value)
				return searchScreenerOptions{}, false
			}
			options.fmp.Limit = limit
		default:
			fmt.Fprintf(stderr, "unknown flag %q\n", arg)
			return searchScreenerOptions{}, false
		}
	}
	return options, true
}

func parseMarketCapAmount(value string) (float64, error) {
	value = strings.TrimSpace(strings.ToUpper(value))
	if value == "" {
		return 0, fmt.Errorf("empty amount")
	}

	multiplier := 1.0
	switch value[len(value)-1] {
	case 'K':
		multiplier = 1_000
		value = value[:len(value)-1]
	case 'M':
		multiplier = 1_000_000
		value = value[:len(value)-1]
	case 'B':
		multiplier = 1_000_000_000
		value = value[:len(value)-1]
	case 'T':
		multiplier = 1_000_000_000_000
		value = value[:len(value)-1]
	}

	amount, err := strconv.ParseFloat(value, 64)
	if err != nil || amount < 0 {
		return 0, fmt.Errorf("invalid amount")
	}
	return amount * multiplier, nil
}

func writeSearchScreenerHelp(w io.Writer) {
	fmt.Fprint(w, `Screen companies.

Usage:
  stonk search screener [flags]

Flags:
  --sector <sector>           Filter by sector
  --country <country>         Filter by country code
  --market-cap-min <amount>   Minimum market cap, such as 100B
  --limit <n>                 Maximum rows to request
  --json                      Write JSON output
  --csv                       Write CSV output
`)
}

func writeScreenerResults(w io.Writer, results []fmp.ScreenerResult, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	case outputCSV:
		return writeScreenerResultsCSV(w, results)
	default:
		return writeScreenerResultsTable(w, results)
	}
}

func writeScreenerResultsTable(w io.Writer, results []fmp.ScreenerResult) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SYMBOL\tNAME\tPRICE\tMARKET CAP\tSECTOR\tINDUSTRY\tEXCHANGE\tCOUNTRY")
	for _, result := range results {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", result.Symbol, result.CompanyName, formatFloat(result.Price), formatFloat(result.MarketCap), result.Sector, result.Industry, result.ExchangeShortName, result.Country)
	}
	return tw.Flush()
}

func writeScreenerResultsCSV(w io.Writer, results []fmp.ScreenerResult) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"symbol", "name", "price", "market_cap", "sector", "industry", "exchange", "country", "is_etf", "is_fund", "is_active"}); err != nil {
		return err
	}
	for _, result := range results {
		if err := cw.Write([]string{
			result.Symbol,
			result.CompanyName,
			formatFloat(result.Price),
			formatFloat(result.MarketCap),
			result.Sector,
			result.Industry,
			result.ExchangeShortName,
			result.Country,
			fmt.Sprint(result.IsETF),
			fmt.Sprint(result.IsFund),
			fmt.Sprint(result.IsActivelyTrading),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
