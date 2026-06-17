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

func runStocksPeers(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeStocksPeersHelp(stdout)
		return 0
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		writeStocksPeersHelp(stdout)
		return 0
	}

	format, remaining, ok := parseOutputFlags(args, stderr)
	if !ok {
		return 2
	}
	if len(remaining) != 1 {
		fmt.Fprintln(stderr, "stocks peers requires exactly one symbol")
		return 2
	}

	apiKey := getenv("FMP_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, "FMP_API_KEY is not configured")
		return 1
	}

	client := fmp.NewClient(apiKey, http.DefaultClient)
	peers, err := client.StockPeers(ctx, remaining[0])
	if err != nil {
		fmt.Fprintf(stderr, "stocks peers failed: %v\n", err)
		return 1
	}

	if err := writeStockPeers(stdout, peers, format); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func writeStocksPeersHelp(w io.Writer) {
	fmt.Fprint(w, `Fetch peer companies.

Usage:
  stonk stocks peers <symbol> [flags]

Flags:
  --json  Write JSON output
  --csv   Write CSV output
`)
}

func writeStockPeers(w io.Writer, peers []fmp.StockPeer, format outputFormat) error {
	switch format {
	case outputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(peers)
	case outputCSV:
		return writeStockPeersCSV(w, peers)
	default:
		return writeStockPeersTable(w, peers)
	}
}

func writeStockPeersTable(w io.Writer, peers []fmp.StockPeer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SYMBOL\tNAME\tPRICE\tMARKET CAP")
	for _, peer := range peers {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", peer.Symbol, peer.CompanyName, formatFloat(peer.Price), formatFloat(peer.MarketCap))
	}
	return tw.Flush()
}

func writeStockPeersCSV(w io.Writer, peers []fmp.StockPeer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"symbol", "name", "price", "market_cap"}); err != nil {
		return err
	}
	for _, peer := range peers {
		if err := cw.Write([]string{
			peer.Symbol,
			peer.CompanyName,
			formatFloat(peer.Price),
			formatFloat(peer.MarketCap),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
