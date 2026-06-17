package cli

import (
	"fmt"
	"io"
)

type outputFormat string

const (
	outputTable outputFormat = "table"
	outputJSON  outputFormat = "json"
	outputCSV   outputFormat = "csv"
)

func parseOutputFlags(args []string, stderr io.Writer) (outputFormat, []string, bool) {
	format := outputTable
	remaining := make([]string, 0, len(args))

	for _, arg := range args {
		switch arg {
		case "--json":
			if format == outputCSV {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return outputTable, nil, false
			}
			format = outputJSON
		case "--csv":
			if format == outputJSON {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return outputTable, nil, false
			}
			format = outputCSV
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return outputTable, nil, false
			}
			remaining = append(remaining, arg)
		}
	}

	return format, remaining, true
}
