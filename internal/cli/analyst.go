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

func runStocksAnalyst(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksAnalystHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksAnalystHelp(stdout)
		return 0
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) != 1 {
		fmt.Fprintln(stderr, "stocks analyst requires exactly one symbol")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	ratings, err := client.StockRatingSnapshot(ctx, remaining[0])
	if err != nil {
		fmt.Fprintf(stderr, "stocks analyst failed: %v\n", err)
		return 1
	}

	if err := writeStockAnalystRatings(stdout, ratings, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeStocksAnalystHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch analyst rating snapshot scores.

Usage:
  stonk stocks analyst <symbol> [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func writeStockAnalystRatings(w io.Writer, ratings []fmp.StockRatingSnapshot, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(ratings)
	case outputCSV:
		return writeStockAnalystRatingsCSV(w, ratings)
	default:
		return writeStockAnalystRatingsTable(w, ratings)
	}
}

func writeStockAnalystRatingsTable(w io.Writer, ratings []fmp.StockRatingSnapshot) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SYMBOL\tRATING\tOVERALL\tDCF\tROE\tROA\tDEBT/EQUITY\tP/E\tP/B")
	for _, rating := range ratings {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			rating.Symbol,
			rating.Rating,
			formatFloat(rating.OverallScore),
			formatFloat(rating.DiscountedCashFlowScore),
			formatFloat(rating.ReturnOnEquityScore),
			formatFloat(rating.ReturnOnAssetsScore),
			formatFloat(rating.DebtToEquityScore),
			formatFloat(rating.PriceToEarningsScore),
			formatFloat(rating.PriceToBookScore),
		)
	}
	return tw.Flush()
}

func writeStockAnalystRatingsCSV(w io.Writer, ratings []fmp.StockRatingSnapshot) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"symbol", "rating", "overall_score", "dcf_score", "roe_score", "roa_score", "debt_to_equity_score", "price_to_earnings_score", "price_to_book_score"}); err != nil {
		return err
	}
	for _, rating := range ratings {
		if err := cw.Write([]string{
			rating.Symbol,
			rating.Rating,
			formatFloat(rating.OverallScore),
			formatFloat(rating.DiscountedCashFlowScore),
			formatFloat(rating.ReturnOnEquityScore),
			formatFloat(rating.ReturnOnAssetsScore),
			formatFloat(rating.DebtToEquityScore),
			formatFloat(rating.PriceToEarningsScore),
			formatFloat(rating.PriceToBookScore),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
