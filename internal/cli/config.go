package cli

import (
	"fmt"
	"io"
)

type getenvFunc func(string) string

func runConfig(args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeConfigHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeConfigHelp(stdout)
		return 0
	case "show":
		writeConfigShow(stdout, getenv)
		return 0
	case "doctor":
		return runConfigDoctor(stdout, getenv)
	case "providers":
		writeConfigProviders(stdout, getenv)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown config command %q\n\n", args[0])
		writeConfigHelp(stderr)
		return 2
	}
}

func writeConfigHelp(w io.Writer) {
	fmt.Fprint(w, `Inspect stonk configuration and provider readiness.

Usage:
  stonk config <command>

Commands:
  show       Print non-secret configuration
  doctor     Check provider configuration
  providers  List configured data providers
`)
}

func writeConfigShow(w io.Writer, getenv getenvFunc) {
	fmt.Fprintln(w, "default_provider=fmp")
	fmt.Fprintf(w, "fmp.api_key_configured=%t\n", hasFMPAPIKey(getenv))
}

func runConfigDoctor(w io.Writer, getenv getenvFunc) int {
	if hasFMPAPIKey(getenv) {
		fmt.Fprintln(w, "ok: FMP_API_KEY is configured")
		return 0
	}

	fmt.Fprintln(w, "error: FMP_API_KEY is not configured")
	return 1
}

func writeConfigProviders(w io.Writer, getenv getenvFunc) {
	state := "missing_key"
	if hasFMPAPIKey(getenv) {
		state = "ready"
	}

	fmt.Fprintf(w, "fmp\t%s\n", state)
}

func hasFMPAPIKey(getenv getenvFunc) bool {
	return getenv("FMP_API_KEY") != ""
}
