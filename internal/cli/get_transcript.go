package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/crertel/stonkler/internal/fmp"
)

type getTranscriptClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	EarningsCallTranscript(context.Context, string, int, int) ([]fmp.EarningsCallTranscript, error)
	EarningsCallTranscriptDates(context.Context, string) ([]fmp.EarningsCallTranscriptDate, error)
}

func runGetTranscript(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetTranscriptHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetTranscriptHelp(stdout)
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
	symbol, err := getTranscriptSymbol(ctx, client, options.symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get transcript failed: %v\n", err)
		return 1
	}
	options.symbol = symbol

	if options.latest {
		dates, err := client.EarningsCallTranscriptDates(ctx, options.symbol)
		if err != nil {
			fmt.Fprintf(stderr, "get transcript failed: %v\n", err)
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
		fmt.Fprintf(stderr, "get transcript failed: %v\n", err)
		return 1
	}

	if err := writeTranscripts(stdout, transcripts, options.format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetTranscriptHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch an earnings call transcript, resolving a name query when needed.

Usage:
  stonk get transcript <symbol|name> --year <year> --quarter <1-4> [flags]
  stonk get transcript <symbol|name> --latest [flags]

Flags:
  --year <year>     Fiscal year
  --quarter <1-4>   Fiscal quarter
  --latest          Fetch the most recent transcript period
  --json            Write JSON output
  --csv             Write CSV output
`)
}

func getTranscriptSymbol(ctx context.Context, client getTranscriptClient, query string) (string, error) {
	if !shouldResolveGetNameQuery(query) {
		return query, nil
	}

	results, err := client.SearchName(ctx, query)
	if err != nil {
		return "", err
	}
	symbol, ok := bestSearchSymbol(results)
	if !ok {
		return "", fmt.Errorf("no symbol found for %q", query)
	}
	return symbol, nil
}
