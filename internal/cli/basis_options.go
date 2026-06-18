package cli

import (
	"fmt"
	"io"
)

type basisOutputOptions struct {
	format    outputFormat
	basisPath string
	remaining []string
}

func parseBasisOutputOptions(args []string, stderr io.Writer) (basisOutputOptions, bool) {
	options := basisOutputOptions{format: outputTable}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--json":
			if options.format == outputCSV {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return basisOutputOptions{}, false
			}
			options.format = outputJSON
		case "--csv":
			if options.format == outputJSON {
				fmt.Fprintln(stderr, "--json and --csv are mutually exclusive")
				return basisOutputOptions{}, false
			}
			options.format = outputCSV
		case "--basis":
			value, ok := nextFlagValue(args, &i, "--basis", stderr)
			if !ok {
				return basisOutputOptions{}, false
			}
			options.basisPath = value
		default:
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(stderr, "unknown flag %q\n", arg)
				return basisOutputOptions{}, false
			}
			options.remaining = append(options.remaining, arg)
		}
	}
	return options, true
}

func resolveBasisPath(flagPath string, getenv getenvFunc) string {
	if flagPath != "" {
		return flagPath
	}
	return getenv("STONK_PORTFOLIO_FILE")
}
