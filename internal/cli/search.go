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

func runSearch(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeSearchHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeSearchHelp(stdout)
		return 0
	}
	if args[0] == "screener" {
		return runSearchScreener(ctx, args[1:], stdout, stderr, getenv)
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) == 0 {
		fmt.Fprintln(stderr, "search requires a query")
		return 2
	}

	mode := "name"
	queryArgs := remaining
	switch remaining[0] {
	case "stocks", "funds", "name", "symbol", "cik", "isin":
		mode = remaining[0]
		queryArgs = remaining[1:]
	}
	if len(queryArgs) == 0 {
		fmt.Fprintf(stderr, "search %s requires a query\n", mode)
		return 2
	}
	if len(queryArgs) > 1 {
		fmt.Fprintln(stderr, "search query must be a single argument; quote multi-word names")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	results, err := runSearchQuery(ctx, client, mode, queryArgs[0])
	if err != nil {
		fmt.Fprintf(stderr, "search failed: %v\n", err)
		return 1
	}

	if err := writeSearchResults(stdout, results, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func runSearchQuery(ctx context.Context, client *fmp.Client, mode string, query string) ([]fmp.SearchResult, error) {
	switch mode {
	case "stocks", "funds", "name":
		return client.SearchName(ctx, query)
	case "symbol":
		return client.SearchSymbol(ctx, query)
	case "cik":
		return client.SearchCIK(ctx, query)
	case "isin":
		return client.SearchISIN(ctx, query)
	default:
		return nil, fmt.Errorf("unsupported search mode %q", mode)
	}
}

func writeSearchHelp(w io.Writer) {
	fmt.Fprint(w, `Discover symbols and securities.

Usage:
  stonk search [flags] <query>
  stonk search [flags] stocks <query>
  stonk search [flags] funds <query>
  stonk search [flags] name <query>
  stonk search [flags] symbol <query>
  stonk search [flags] cik <cik>
  stonk search [flags] isin <isin>
  stonk search screener [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output

Screener Flags:
  --sector <sector>           Filter by sector
  --country <country>         Filter by country code
  --market-cap-min <amount>   Minimum market cap, such as 100B
  --limit <n>                 Maximum rows to request
`)
}

func writeSearchResults(w io.Writer, results []fmp.SearchResult, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	case outputCSV:
		return writeSearchResultsCSV(w, results)
	default:
		return writeSearchResultsTable(w, results)
	}
}

func writeSearchResultsTable(w io.Writer, results []fmp.SearchResult) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SYMBOL\tNAME\tEXCHANGE\tCURRENCY")
	for _, result := range results {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", result.Symbol, result.Name, result.ExchangeShortName, result.Currency)
	}
	return tw.Flush()
}

func writeSearchResultsCSV(w io.Writer, results []fmp.SearchResult) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"symbol", "name", "exchange", "currency"}); err != nil {
		return err
	}
	for _, result := range results {
		if err := cw.Write([]string{result.Symbol, result.Name, result.ExchangeShortName, result.Currency}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
