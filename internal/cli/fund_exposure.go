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

func runFundsExposure(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeFundsExposureHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeFundsExposureHelp(stdout)
		return 0
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) != 1 {
		fmt.Fprintln(stderr, "funds exposure requires exactly one asset symbol")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	exposures, err := client.ETFAssetExposure(ctx, remaining[0])
	if err != nil {
		fmt.Fprintf(stderr, "funds exposure failed: %v\n", err)
		return 1
	}

	if err := writeETFAssetExposures(stdout, exposures, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeFundsExposureHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch ETF or fund exposure to an asset.

Usage:
  stonk funds exposure <asset-symbol> [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func writeETFAssetExposures(w io.Writer, exposures []fmp.ETFAssetExposure, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(exposures)
	case outputCSV:
		return writeETFAssetExposuresCSV(w, exposures)
	default:
		return writeETFAssetExposuresTable(w, exposures)
	}
}

func writeETFAssetExposuresTable(w io.Writer, exposures []fmp.ETFAssetExposure) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SYMBOL\tASSET\tNAME\tSHARES\tWEIGHT%\tMARKET VALUE")
	for _, exposure := range exposures {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s\t%s\t%s\n",
			rawValue(exposure, "symbol"),
			rawValue(exposure, "asset"),
			rawValue(exposure, "name"),
			rawValue(exposure, "sharesNumber"),
			rawValue(exposure, "weightPercentage"),
			rawValue(exposure, "marketValue"),
		)
	}
	return tw.Flush()
}

func writeETFAssetExposuresCSV(w io.Writer, exposures []fmp.ETFAssetExposure) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"symbol", "asset", "name", "isin", "cusip", "shares", "weight_percent", "market_value"}); err != nil {
		return err
	}
	for _, exposure := range exposures {
		if err := cw.Write([]string{
			rawValue(exposure, "symbol"),
			rawValue(exposure, "asset"),
			rawValue(exposure, "name"),
			rawValue(exposure, "isin"),
			rawValue(exposure, "cusip"),
			rawValue(exposure, "sharesNumber"),
			rawValue(exposure, "weightPercentage"),
			rawValue(exposure, "marketValue"),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
