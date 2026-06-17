package fmp

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// StockQuote is a normalized stock quote from FMP.
type StockQuote struct {
	Symbol           string  `json:"symbol"`
	Name             string  `json:"name,omitempty"`
	Price            float64 `json:"price"`
	Change           float64 `json:"change"`
	ChangePercentage float64 `json:"changePercentage"`
	Volume           float64 `json:"volume"`
	MarketCap        float64 `json:"marketCap"`
	Exchange         string  `json:"exchange,omitempty"`
	Timestamp        int64   `json:"timestamp"`
}

// StockQuotes returns current quote data for one or more stock symbols.
func (c *Client) StockQuotes(ctx context.Context, symbols []string) ([]StockQuote, error) {
	symbols = normalizeSymbols(symbols)
	if len(symbols) == 0 {
		return nil, fmt.Errorf("at least one symbol is required")
	}

	var quotes []StockQuote
	query := url.Values{}
	query.Set("symbols", strings.Join(symbols, ","))
	if err := c.get(ctx, "/batch-quote", query, &quotes); err != nil {
		return nil, err
	}
	return quotes, nil
}

func normalizeSymbols(symbols []string) []string {
	normalized := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		symbol = strings.TrimSpace(symbol)
		if symbol == "" {
			continue
		}
		normalized = append(normalized, strings.ToUpper(symbol))
	}
	return normalized
}
