package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/crertel/stonkler/internal/fmp"
)

type getStatementsClient interface {
	SearchName(context.Context, string) ([]fmp.SearchResult, error)
	StockStatements(context.Context, fmp.StockStatementsRequest) ([]fmp.FinancialStatement, error)
}

func runGetStatements(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetStatementsHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeGetStatementsHelp(stdout)
		return 0
	}

	options, ok := parseGetStatementsOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	symbol, err := getStatementsSymbol(ctx, client, options.symbol)
	if err != nil {
		fmt.Fprintf(stderr, "get statements failed: %v\n", err)
		return 1
	}
	options.symbol = symbol

	statements, err := client.StockStatements(ctx, fmp.StockStatementsRequest{
		Symbol:    options.symbol,
		Statement: options.statement,
		Period:    options.period,
		Limit:     options.limit,
	})
	if err != nil {
		fmt.Fprintf(stderr, "get statements failed: %v\n", err)
		return 1
	}

	if err := writeStockStatements(stdout, statements, options); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeGetStatementsHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch financial statements, resolving a name query and defaulting to income statements.

Usage:
  stonk get statements <symbol|name> [statement] [flags]
  stonk get statement <symbol|name> [statement] [flags]

Flags:
  --period <annual|quarter>  Reporting period
  --limit <n>                Maximum statements to request
  --json                     Write JSON output
  --csv                      Write CSV output
`)
}

func parseGetStatementsOptions(args []string, stderr io.Writer) (statementsOptions, bool) {
	options := statementsOptions{
		format:    outputTable,
		period:    "annual",
		limit:     5,
		statement: fmp.StatementIncome,
		displaySpecs: []statementDisplaySpec{
			{"REVENUE", "revenue"},
			{"GROSS PROFIT", "grossProfit"},
			{"OPERATING INCOME", "operatingIncome"},
			{"NET INCOME", "netIncome"},
			{"EPS", "eps"},
		},
	}

	positionals := make([]string, 0, 2)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--json":
			if options.format == outputCSV {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return statementsOptions{}, false
			}
			options.format = outputJSON
		case "--csv":
			if options.format == outputJSON {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return statementsOptions{}, false
			}
			options.format = outputCSV
		case "--period":
			value, ok := nextFlagValue(args, &i, "--period", stderr)
			if !ok {
				return statementsOptions{}, false
			}
			switch value {
			case "annual", "quarter":
				options.period = value
			default:
				fmt.Fprintf(stderr, "invalid --period value %q; use annual or quarter\n", value)
				return statementsOptions{}, false
			}
		case "--limit":
			value, ok := nextFlagValue(args, &i, "--limit", stderr)
			if !ok {
				return statementsOptions{}, false
			}
			limit, err := strconv.Atoi(value)
			if err != nil || limit < 0 {
				fmt.Fprintf(stderr, "invalid --limit value %q\n", value)
				return statementsOptions{}, false
			}
			options.limit = limit
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return statementsOptions{}, false
			}
			positionals = append(positionals, arg)
		}
	}

	if len(positionals) < 1 || len(positionals) > 2 {
		fmt.Fprintln(stderr, "get statements requires a symbol or name and optional statement type")
		return statementsOptions{}, false
	}
	options.symbol = positionals[0]
	if len(positionals) == 2 {
		statement, specs, ok := parseStatementType(positionals[1], stderr)
		if !ok {
			return statementsOptions{}, false
		}
		options.statement = statement
		options.displaySpecs = specs
	}
	return options, true
}

func getStatementsSymbol(ctx context.Context, client getStatementsClient, query string) (string, error) {
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
