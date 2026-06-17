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

type transcriptOptions struct {
	format  outputFormat
	symbol  string
	year    int
	quarter int
	latest  bool
}

func runStocksTranscript(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksTranscriptHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksTranscriptHelp(stdout)
		return 0
	}

	options, ok := parseTranscriptOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	if options.latest {
		dates, err := client.EarningsCallTranscriptDates(ctx, options.symbol)
		if err != nil {
			fmt.Fprintf(stderr, "stocks transcript failed: %v\n", err)
			return 1
		}

		year, quarter, ok := latestTranscriptPeriod(dates)
		if !ok {
			fmt.Fprintf(stderr, "no transcript dates found for %s\n", options.symbol)
			return 1
		}
		options.year = year
		options.quarter = quarter
	}

	transcripts, err := client.EarningsCallTranscript(ctx, options.symbol, options.year, options.quarter)
	if err != nil {
		fmt.Fprintf(stderr, "stocks transcript failed: %v\n", err)
		return 1
	}

	if err := writeTranscripts(stdout, transcripts, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func parseTranscriptOptions(args []string, stderr io.Writer) (transcriptOptions, bool) {
	options := transcriptOptions{format: outputTable}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--json":
			if options.format == outputCSV {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return transcriptOptions{}, false
			}
			options.format = outputJSON
		case "--csv":
			if options.format == outputJSON {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return transcriptOptions{}, false
			}
			options.format = outputCSV
		case "--latest":
			options.latest = true
		case "--year":
			value, ok := nextFlagValue(args, &i, "--year", stderr)
			if !ok {
				return transcriptOptions{}, false
			}
			year, err := strconv.Atoi(value)
			if err != nil || year <= 0 {
				fmt.Fprintf(stderr, "invalid --year value %q\n", value)
				return transcriptOptions{}, false
			}
			options.year = year
		case "--quarter":
			value, ok := nextFlagValue(args, &i, "--quarter", stderr)
			if !ok {
				return transcriptOptions{}, false
			}
			quarter, err := strconv.Atoi(value)
			if err != nil || quarter < 1 || quarter > 4 {
				fmt.Fprintf(stderr, "invalid --quarter value %q\n", value)
				return transcriptOptions{}, false
			}
			options.quarter = quarter
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return transcriptOptions{}, false
			}
			if options.symbol != "" {
				fmt.Fprintln(stderr, "stocks transcript requires exactly one symbol")
				return transcriptOptions{}, false
			}
			options.symbol = arg
		}
	}

	if options.symbol == "" {
		fmt.Fprintln(stderr, "stocks transcript requires exactly one symbol")
		return transcriptOptions{}, false
	}
	if options.latest && (options.year != 0 || options.quarter != 0) {
		fmt.Fprintln(stderr, "--latest cannot be combined with --year or --quarter")
		return transcriptOptions{}, false
	}
	if options.latest {
		return options, true
	}
	if options.year == 0 {
		fmt.Fprintln(stderr, "stocks transcript requires --year")
		return transcriptOptions{}, false
	}
	if options.quarter == 0 {
		fmt.Fprintln(stderr, "stocks transcript requires --quarter")
		return transcriptOptions{}, false
	}
	return options, true
}

func writeStocksTranscriptHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch an earnings call transcript.

Usage:
  stonk stocks transcript <symbol> --year <year> --quarter <1-4> [flags]
  stonk stocks transcript <symbol> --latest [flags]

Flags:
  --year <year>     Fiscal year
  --quarter <1-4>   Fiscal quarter
  --latest          Fetch the most recent transcript period
  --json            Write JSON output
  --csv             Write CSV output
`)
}

func latestTranscriptPeriod(dates []fmp.EarningsCallTranscriptDate) (int, int, bool) {
	var best fmp.EarningsCallTranscriptDate
	found := false
	for _, date := range dates {
		if date.Year <= 0 || date.Quarter < 1 || date.Quarter > 4 {
			continue
		}
		if !found || newerTranscriptDate(date, best) {
			best = date
			found = true
		}
	}
	if !found {
		return 0, 0, false
	}
	return best.Year, best.Quarter, true
}

func newerTranscriptDate(candidate, current fmp.EarningsCallTranscriptDate) bool {
	if candidate.Date != "" || current.Date != "" {
		return candidate.Date > current.Date
	}
	if candidate.Year != current.Year {
		return candidate.Year > current.Year
	}
	return candidate.Quarter > current.Quarter
}

func writeTranscripts(w io.Writer, transcripts []fmp.EarningsCallTranscript, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(transcripts)
	case outputCSV:
		return writeTranscriptsCSV(w, transcripts)
	default:
		return writeTranscriptsTable(w, transcripts)
	}
}

func writeTranscriptsTable(w io.Writer, transcripts []fmp.EarningsCallTranscript) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SYMBOL\tYEAR\tQUARTER\tDATE\tTITLE\tCONTENT CHARS")
	for _, transcript := range transcripts {
		content := rawValue(transcript, "content")
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s\t%s\t%d\n",
			rawValue(transcript, "symbol"),
			rawValue(transcript, "year"),
			rawValue(transcript, "quarter"),
			rawValue(transcript, "date"),
			rawValue(transcript, "title"),
			len(content),
		)
	}
	return tw.Flush()
}

func writeTranscriptsCSV(w io.Writer, transcripts []fmp.EarningsCallTranscript) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"symbol", "year", "quarter", "date", "title", "content"}); err != nil {
		return err
	}
	for _, transcript := range transcripts {
		if err := cw.Write([]string{
			rawValue(transcript, "symbol"),
			rawValue(transcript, "year"),
			rawValue(transcript, "quarter"),
			rawValue(transcript, "date"),
			rawValue(transcript, "title"),
			rawValue(transcript, "content"),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func rawValue(row map[string]any, field string) string {
	value, ok := row[field]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return formatFloat(typed)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return fmt.Sprint(typed)
	}
}
