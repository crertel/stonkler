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

func runStocksProfile(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksProfileHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksProfileHelp(stdout)
		return 0
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) != 1 {
		fmt.Fprintln(stderr, "stocks profile requires exactly one symbol")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	profile, err := client.StockProfile(ctx, remaining[0])
	if err != nil {
		fmt.Fprintf(stderr, "stocks profile failed: %v\n", err)
		return 1
	}

	if err := writeStockProfile(stdout, profile, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeStockProfile(w io.Writer, profile fmp.StockProfile, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(profile)
	case outputCSV:
		return writeStockProfileCSV(w, profile)
	default:
		return writeStockProfileTable(w, profile)
	}
}

func writeStockProfileTable(w io.Writer, profile fmp.StockProfile) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	rows := [][2]string{
		{"Symbol", profile.Symbol},
		{"Company", profile.CompanyName},
		{"Exchange", profile.Exchange},
		{"Currency", profile.Currency},
		{"Sector", profile.Sector},
		{"Industry", profile.Industry},
		{"CEO", profile.CEO},
		{"Country", profile.Country},
		{"Market Cap", formatFloat(profile.MarketCap)},
		{"Price", formatFloat(profile.Price)},
		{"Website", profile.Website},
		{"IPO Date", profile.IPODate},
	}
	for _, row := range rows {
		fmt.Fprintf(tw, "%s:\t%s\n", row[0], row[1])
	}
	return tw.Flush()
}

func writeStockProfileCSV(w io.Writer, profile fmp.StockProfile) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"symbol",
		"company_name",
		"exchange",
		"currency",
		"sector",
		"industry",
		"ceo",
		"country",
		"market_cap",
		"price",
		"website",
		"ipo_date",
	}); err != nil {
		return err
	}

	if err := cw.Write([]string{
		profile.Symbol,
		profile.CompanyName,
		profile.Exchange,
		profile.Currency,
		profile.Sector,
		profile.Industry,
		profile.CEO,
		profile.Country,
		formatFloat(profile.MarketCap),
		formatFloat(profile.Price),
		profile.Website,
		profile.IPODate,
	}); err != nil {
		return err
	}
	cw.Flush()
	return cw.Error()
}
