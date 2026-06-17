package cli

import (
	"context"
	"fmt"
	"io"
)

func runGet(ctx context.Context, args []string, stdout, stderr io.Writer, getenv getenvFunc) int {
	if len(args) == 0 {
		writeGetHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeGetHelp(stdout)
		return 0
	case "analyst":
		return runStocksAnalyst(ctx, args[1:], stdout, stderr, getenv)
	case "company", "profile":
		return runStocksProfile(ctx, args[1:], stdout, stderr, getenv)
	case "country-weightings":
		return runFundsCountryWeightings(ctx, args[1:], stdout, stderr, getenv)
	case "commodity", "commodities":
		return runDomainQuote(ctx, args[1:], stdout, stderr, getenv, "commodities", writeCommoditiesQuoteHelp, commodityQuotes)
	case "commodity-history", "commodities-history":
		return runDomainHistory(ctx, args[1:], stdout, stderr, getenv, "commodities", writeCommoditiesHistoryHelp, nil)
	case "crypto":
		return runDomainQuote(ctx, args[1:], stdout, stderr, getenv, "crypto", writeCryptoQuoteHelp, cryptoQuotes)
	case "crypto-history":
		return runDomainHistory(ctx, args[1:], stdout, stderr, getenv, "crypto", writeCryptoHistoryHelp, nil)
	case "exposure":
		return runFundsExposure(ctx, args[1:], stdout, stderr, getenv)
	case "etf", "fund", "fund-info":
		return runFundsInfo(ctx, args[1:], stdout, stderr, getenv)
	case "forex":
		return runDomainQuote(ctx, args[1:], stdout, stderr, getenv, "forex", writeForexQuoteHelp, forexQuotes)
	case "forex-history":
		return runDomainHistory(ctx, args[1:], stdout, stderr, getenv, "forex", writeForexHistoryHelp, nil)
	case "history":
		return runStocksHistory(ctx, args[1:], stdout, stderr, getenv)
	case "holdings":
		return runFundsHoldings(ctx, args[1:], stdout, stderr, getenv)
	case "index", "indexes":
		return runDomainQuote(ctx, args[1:], stdout, stderr, getenv, "indexes", writeIndexesQuoteHelp, indexQuotes)
	case "index-history", "indexes-history":
		return runDomainHistory(ctx, args[1:], stdout, stderr, getenv, "indexes", writeIndexesHistoryHelp, normalizeIndexSymbol)
	case "insiders":
		return runStocksInsiders(ctx, args[1:], stdout, stderr, getenv)
	case "metrics":
		return runStocksMetrics(ctx, args[1:], stdout, stderr, getenv)
	case "peers":
		return runStocksPeers(ctx, args[1:], stdout, stderr, getenv)
	case "quote", "quotes":
		return runStocksQuote(ctx, args[1:], stdout, stderr, getenv)
	case "ratios":
		return runStocksRatios(ctx, args[1:], stdout, stderr, getenv)
	case "sec":
		return runStocksSEC(ctx, args[1:], stdout, stderr, getenv)
	case "sector-weightings":
		return runFundsSectorWeightings(ctx, args[1:], stdout, stderr, getenv)
	case "statement", "statements":
		return runStocksStatements(ctx, args[1:], stdout, stderr, getenv)
	case "transcript":
		return runStocksTranscript(ctx, args[1:], stdout, stderr, getenv)
	default:
		fmt.Fprintf(stderr, "unknown get command %q\n\n", args[0])
		writeGetHelp(stderr)
		return 2
	}
}

func writeGetHelp(w io.Writer) {
	fmt.Fprint(w, `Workflow-oriented shortcuts.

Usage:
  stonk get <command> [flags]

Commands:
  analyst Fetch stock analyst rating snapshot, inferring the stock domain for now
  company Fetch company profile data, inferring the stock domain for now
  country-weightings Fetch ETF or fund country allocation weights
  commodity Fetch commodity quotes
  commodity-history Fetch commodity history
  crypto  Fetch cryptocurrency quotes
  crypto-history Fetch cryptocurrency history
  exposure Fetch ETF or fund exposure to an asset
  fund    Fetch ETF or fund profile information
  forex   Fetch foreign exchange quotes
  forex-history Fetch foreign exchange history
  history Fetch historical prices, inferring the stock domain for now
  holdings Fetch ETF holdings, inferring the funds domain for now
  index   Fetch index quotes
  index-history Fetch index history
  insiders Fetch stock insider transactions, inferring the stock domain for now
  metrics Fetch stock key metrics, inferring the stock domain for now
  peers   Fetch stock peers, inferring the stock domain for now
  profile Alias for company
  quote   Fetch one or more quotes, inferring the stock domain for now
  quotes  Alias for quote
  ratios  Fetch stock ratios, inferring the stock domain for now
  sec     Fetch stock SEC filings, inferring the stock domain for now
  sector-weightings Fetch ETF or fund sector allocation weights
  statements Fetch financial statements, inferring the stock domain for now
  transcript Fetch an earnings call transcript, inferring the stock domain for now
`)
}
