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

func runFundsSectorWeightings(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeFundsSectorWeightingsHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeFundsSectorWeightingsHelp(stdout)
		return 0
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) != 1 {
		fmt.Fprintln(stderr, "funds sector-weightings requires exactly one symbol")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	weightings, err := client.ETFSectorWeightings(ctx, remaining[0])
	if err != nil {
		fmt.Fprintf(stderr, "funds sector-weightings failed: %v\n", err)
		return 1
	}

	if err := writeETFSectorWeightings(stdout, weightings, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeFundsSectorWeightingsHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch ETF or fund sector allocation weights.

Usage:
  stonk funds sector-weightings <symbol> [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func writeETFSectorWeightings(w io.Writer, weightings []fmp.ETFSectorWeighting, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(weightings)
	case outputCSV:
		return writeETFSectorWeightingsCSV(w, weightings)
	default:
		return writeETFSectorWeightingsTable(w, weightings)
	}
}

func writeETFSectorWeightingsTable(w io.Writer, weightings []fmp.ETFSectorWeighting) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SYMBOL\tSECTOR\tWEIGHT%")
	for _, weighting := range weightings {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", weighting.Symbol, weighting.Sector, formatFloat(weighting.WeightPercentage))
	}
	return tw.Flush()
}

func writeETFSectorWeightingsCSV(w io.Writer, weightings []fmp.ETFSectorWeighting) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"symbol", "sector", "weight_percent"}); err != nil {
		return err
	}
	for _, weighting := range weightings {
		if err := cw.Write([]string{
			weighting.Symbol,
			weighting.Sector,
			formatFloat(weighting.WeightPercentage),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
