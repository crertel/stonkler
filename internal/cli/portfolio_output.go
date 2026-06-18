package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"
)

func writeBasisEntries(w io.Writer, entries []basisEntry, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	case outputCSV:
		return writeBasisEntriesCSV(w, entries)
	default:
		return writeBasisEntriesTable(w, entries)
	}
}

func writeBasisEntriesTable(w io.Writer, entries []basisEntry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "DOMAIN\tSYMBOL\tLOTS\tBASIS\tQUANTITY\tCOST\tACQUIRED")
	for _, entry := range entries {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
			entry.Domain,
			entry.Symbol,
			len(entry.Lots),
			optionalFloat(entry.AverageBasis, entry.HasAverageBasis),
			optionalFloat(entry.TotalQuantity, entry.HasTotalQuantity),
			optionalFloat(entry.TotalCost, entry.HasTotalCost),
			optionalString(entry.AcquiredOn),
		)
	}
	return tw.Flush()
}

func writeBasisEntriesCSV(w io.Writer, entries []basisEntry) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"domain", "symbol", "lots", "basis", "quantity", "cost", "acquired_on"}); err != nil {
		return err
	}
	for _, entry := range entries {
		if err := cw.Write([]string{
			entry.Domain,
			entry.Symbol,
			strconv.Itoa(len(entry.Lots)),
			optionalFloat(entry.AverageBasis, entry.HasAverageBasis),
			optionalFloat(entry.TotalQuantity, entry.HasTotalQuantity),
			optionalFloat(entry.TotalCost, entry.HasTotalCost),
			entry.AcquiredOn,
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeQuotesWithBasis(w io.Writer, rows []quoteWithBasis, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(rows)
	case outputCSV:
		return writeQuotesWithBasisCSV(w, rows)
	default:
		return writeQuotesWithBasisTable(w, rows)
	}
}

func writeQuotesWithBasisTable(w io.Writer, rows []quoteWithBasis) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SYMBOL\tNAME\tPRICE\tCHANGE%\tBASIS\tQTY\tVALUE\tGAIN\tGAIN%\tACQUIRED")
	for _, row := range rows {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			row.Symbol,
			row.Name,
			formatFloat(row.Price),
			formatFloat(row.ChangePercentage),
			optionalFloat(row.Basis, row.HasBasis),
			optionalFloat(row.Quantity, row.HasQuantity),
			optionalFloat(row.MarketValue, row.HasMarketValue),
			optionalFloat(row.UnrealizedGain, row.HasUnrealizedGain),
			optionalFloat(row.UnrealizedGainPct, row.HasUnrealizedPct),
			optionalString(row.AcquiredOn),
		)
	}
	return tw.Flush()
}

func writeQuotesWithBasisCSV(w io.Writer, rows []quoteWithBasis) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"symbol", "name", "price", "change_percent", "basis", "quantity", "market_value", "unrealized_gain", "unrealized_gain_percent", "acquired_on"}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := cw.Write([]string{
			row.Symbol,
			row.Name,
			formatFloat(row.Price),
			formatFloat(row.ChangePercentage),
			optionalFloat(row.Basis, row.HasBasis),
			optionalFloat(row.Quantity, row.HasQuantity),
			optionalFloat(row.MarketValue, row.HasMarketValue),
			optionalFloat(row.UnrealizedGain, row.HasUnrealizedGain),
			optionalFloat(row.UnrealizedGainPct, row.HasUnrealizedPct),
			row.AcquiredOn,
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func optionalFloat(value float64, ok bool) string {
	if !ok {
		return "-"
	}
	return formatFloat(value)
}

func optionalString(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
