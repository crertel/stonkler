package fmp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Quote is a normalized quote from FMP.
type Quote struct {
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

// StockQuote is kept as the stock-domain name for shared quote data.
type StockQuote = Quote

// UnmarshalJSON accepts both stable and v3 quote field variants.
func (q *Quote) UnmarshalJSON(data []byte) error {
	var raw struct {
		Symbol            string   `json:"symbol"`
		Name              string   `json:"name"`
		Price             *float64 `json:"price"`
		Change            *float64 `json:"change"`
		ChangePercentage  *float64 `json:"changePercentage"`
		ChangesPercentage *float64 `json:"changesPercentage"`
		Volume            *float64 `json:"volume"`
		MarketCap         *float64 `json:"marketCap"`
		Exchange          string   `json:"exchange"`
		Timestamp         *int64   `json:"timestamp"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	q.Symbol = raw.Symbol
	q.Name = raw.Name
	q.Price = valueOrZero(raw.Price)
	q.Change = valueOrZero(raw.Change)
	q.ChangePercentage = valueOrZero(raw.ChangePercentage)
	if raw.ChangePercentage == nil {
		q.ChangePercentage = valueOrZero(raw.ChangesPercentage)
	}
	q.Volume = valueOrZero(raw.Volume)
	q.MarketCap = valueOrZero(raw.MarketCap)
	q.Exchange = raw.Exchange
	q.Timestamp = intValueOrZero(raw.Timestamp)
	return nil
}

// StockQuotes returns current quote data for one or more stock symbols.
func (c *Client) StockQuotes(ctx context.Context, symbols []string) ([]StockQuote, error) {
	return c.BatchQuotes(ctx, symbols)
}

// BatchQuotes returns current quote data for symbols supported by the stable batch quote endpoint.
func (c *Client) BatchQuotes(ctx context.Context, symbols []string) ([]Quote, error) {
	symbols = normalizeSymbols(symbols)
	if len(symbols) == 0 {
		return nil, fmt.Errorf("at least one symbol is required")
	}

	var quotes []Quote
	query := url.Values{}
	query.Set("symbols", strings.Join(symbols, ","))
	if err := c.get(ctx, "/batch-quote", query, &quotes); err != nil {
		return nil, err
	}
	return quotes, nil
}

// IndexQuotes returns current quote data for index symbols.
func (c *Client) IndexQuotes(ctx context.Context, symbols []string) ([]Quote, error) {
	symbols = normalizeIndexSymbols(symbols)
	if len(symbols) == 0 {
		return nil, fmt.Errorf("at least one symbol is required")
	}

	var quotes []Quote
	if err := c.getV3(ctx, "/quote/"+url.PathEscape(strings.Join(symbols, ",")), url.Values{}, &quotes); err != nil {
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

func normalizeIndexSymbols(symbols []string) []string {
	normalized := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		symbol = strings.TrimSpace(symbol)
		if symbol == "" {
			continue
		}
		symbol = strings.ToUpper(symbol)
		if !strings.HasPrefix(symbol, "^") {
			symbol = "^" + symbol
		}
		normalized = append(normalized, symbol)
	}
	return normalized
}

func valueOrZero(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func intValueOrZero(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}
