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

type statementsOptions struct {
	format       outputFormat
	symbol       string
	statement    fmp.StatementType
	period       string
	limit        int
	displaySpecs []statementDisplaySpec
}

type statementDisplaySpec struct {
	Header string
	Field  string
}

func runStocksStatements(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksStatementsHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksStatementsHelp(stdout)
		return 0
	}

	options, ok := parseStatementsOptions(args, stderr)
	if !ok {
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	statements, err := client.StockStatements(ctx, fmp.StockStatementsRequest{
		Symbol:    options.symbol,
		Statement: options.statement,
		Period:    options.period,
		Limit:     options.limit,
	})
	if err != nil {
		fmt.Fprintf(stderr, "stocks statements failed: %v\n", err)
		return 1
	}

	if err := writeStockStatements(stdout, statements, options); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func parseStatementsOptions(args []string, stderr io.Writer) (statementsOptions, bool) {
	options := statementsOptions{
		format: outputTable,
		period: "annual",
		limit:  5,
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

	if len(positionals) != 2 {
		fmt.Fprintln(stderr, "stocks statements requires a symbol and statement type")
		return statementsOptions{}, false
	}
	statement, specs, ok := parseStatementType(positionals[1], stderr)
	if !ok {
		return statementsOptions{}, false
	}
	options.symbol = positionals[0]
	options.statement = statement
	options.displaySpecs = specs
	return options, true
}

func parseStatementType(value string, stderr io.Writer) (fmp.StatementType, []statementDisplaySpec, bool) {
	switch value {
	case "income", "income-statement":
		return fmp.StatementIncome, []statementDisplaySpec{
			{"REVENUE", "revenue"},
			{"GROSS PROFIT", "grossProfit"},
			{"OPERATING INCOME", "operatingIncome"},
			{"NET INCOME", "netIncome"},
			{"EPS", "eps"},
		}, true
	case "balance", "balance-sheet":
		return fmp.StatementBalanceSheet, []statementDisplaySpec{
			{"ASSETS", "totalAssets"},
			{"LIABILITIES", "totalLiabilities"},
			{"EQUITY", "totalStockholdersEquity"},
			{"CASH", "cashAndCashEquivalents"},
			{"DEBT", "totalDebt"},
		}, true
	case "cash-flow", "cashflow", "cash":
		return fmp.StatementCashFlow, []statementDisplaySpec{
			{"OPERATING CF", "operatingCashFlow"},
			{"CAPEX", "capitalExpenditure"},
			{"FREE CF", "freeCashFlow"},
			{"ENDING CASH", "cashAtEndOfPeriod"},
			{"NET CHANGE CASH", "netChangeInCash"},
		}, true
	default:
		fmt.Fprintf(stderr, "unknown statement type %q; use income, balance, or cash-flow\n", value)
		return "", nil, false
	}
}

func writeStockStatements(w io.Writer, statements []fmp.FinancialStatement, options statementsOptions) error {
	switch options.format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(statements)
	case outputCSV:
		return writeStockStatementsCSV(w, statements, options.displaySpecs)
	default:
		return writeStockStatementsTable(w, statements, options.displaySpecs)
	}
}

func writeStockStatementsTable(w io.Writer, statements []fmp.FinancialStatement, specs []statementDisplaySpec) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprint(tw, "DATE\tSYMBOL\tFY\tPERIOD")
	for _, spec := range specs {
		fmt.Fprintf(tw, "\t%s", spec.Header)
	}
	fmt.Fprintln(tw)

	for _, statement := range statements {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s",
			statementString(statement, "date"),
			statementString(statement, "symbol"),
			statementString(statement, "fiscalYear"),
			statementString(statement, "period"),
		)
		for _, spec := range specs {
			fmt.Fprintf(tw, "\t%s", statementValue(statement, spec.Field))
		}
		fmt.Fprintln(tw)
	}
	return tw.Flush()
}

func writeStockStatementsCSV(w io.Writer, statements []fmp.FinancialStatement, specs []statementDisplaySpec) error {
	cw := csv.NewWriter(w)
	header := []string{"date", "symbol", "fiscal_year", "period"}
	for _, spec := range specs {
		header = append(header, spec.Field)
	}
	if err := cw.Write(header); err != nil {
		return err
	}

	for _, statement := range statements {
		row := []string{
			statementString(statement, "date"),
			statementString(statement, "symbol"),
			statementString(statement, "fiscalYear"),
			statementString(statement, "period"),
		}
		for _, spec := range specs {
			row = append(row, statementValue(statement, spec.Field))
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func statementString(statement fmp.FinancialStatement, field string) string {
	value, ok := statement[field]
	if !ok || value == nil {
		return ""
	}
	if text, ok := value.(string); ok {
		return text
	}
	return fmt.Sprint(value)
}

func statementValue(statement fmp.FinancialStatement, field string) string {
	value, ok := statement[field]
	if !ok || value == nil {
		return ""
	}
	switch value := value.(type) {
	case float64:
		return formatFloat(value)
	case string:
		return value
	default:
		return fmt.Sprint(value)
	}
}
