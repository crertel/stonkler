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

func runFundsInfo(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeFundsInfoHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeFundsInfoHelp(stdout)
		return 0
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) != 1 {
		fmt.Fprintln(stderr, "funds info requires exactly one symbol")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	profile, err := client.FundProfile(ctx, remaining[0])
	if err != nil {
		fmt.Fprintf(stderr, "funds info failed: %v\n", err)
		return 1
	}

	if err := writeFundInfo(stdout, profile, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeFundInfo(w io.Writer, profile fmp.StockProfile, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(profile)
	case outputCSV:
		return writeFundInfoCSV(w, profile)
	default:
		return writeFundInfoTable(w, profile)
	}
}

func writeFundInfoTable(w io.Writer, profile fmp.StockProfile) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	rows := [][2]string{
		{"Symbol", profile.Symbol},
		{"Name", profile.CompanyName},
		{"Exchange", profile.Exchange},
		{"Currency", profile.Currency},
		{"Category", profile.Industry},
		{"Issuer/Sector", profile.Sector},
		{"Market Cap", formatFloat(profile.MarketCap)},
		{"Price", formatFloat(profile.Price)},
		{"Dividend", formatFloat(profile.LastDividend)},
		{"Website", profile.Website},
		{"IPO Date", profile.IPODate},
		{"Active", fmt.Sprint(profile.IsActivelyTrading)},
	}
	for _, row := range rows {
		fmt.Fprintf(tw, "%s:\t%s\n", row[0], row[1])
	}
	return tw.Flush()
}

func writeFundInfoCSV(w io.Writer, profile fmp.StockProfile) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"symbol",
		"name",
		"exchange",
		"currency",
		"category",
		"issuer_sector",
		"market_cap",
		"price",
		"last_dividend",
		"website",
		"ipo_date",
		"is_etf",
		"is_fund",
		"is_active",
	}); err != nil {
		return err
	}

	if err := cw.Write([]string{
		profile.Symbol,
		profile.CompanyName,
		profile.Exchange,
		profile.Currency,
		profile.Industry,
		profile.Sector,
		formatFloat(profile.MarketCap),
		formatFloat(profile.Price),
		formatFloat(profile.LastDividend),
		profile.Website,
		profile.IPODate,
		fmt.Sprint(profile.IsETF),
		fmt.Sprint(profile.IsFund),
		fmt.Sprint(profile.IsActivelyTrading),
	}); err != nil {
		return err
	}
	cw.Flush()
	return cw.Error()
}
