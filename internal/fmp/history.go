package fmp

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// StockHistoryRequest describes a historical end-of-day price request.
type StockHistoryRequest struct {
	Symbol string
	From   string
	To     string
}

// StockPrice is one historical end-of-day price row from FMP.
type StockPrice struct {
	Symbol        string  `json:"symbol"`
	Date          string  `json:"date"`
	Open          float64 `json:"open"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Close         float64 `json:"close"`
	Volume        float64 `json:"volume"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
	VWAP          float64 `json:"vwap"`
}

// StockHistory returns historical end-of-day prices for a stock symbol.
func (c *Client) StockHistory(ctx context.Context, request StockHistoryRequest) ([]StockPrice, error) {
	return c.PriceHistory(ctx, request)
}

// PriceHistory returns historical end-of-day prices for any supported symbol.
func (c *Client) PriceHistory(ctx context.Context, request StockHistoryRequest) ([]StockPrice, error) {
	symbol := strings.ToUpper(strings.TrimSpace(request.Symbol))
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	query := url.Values{}
	query.Set("symbol", symbol)
	if request.From != "" {
		query.Set("from", request.From)
	}
	if request.To != "" {
		query.Set("to", request.To)
	}

	var prices []StockPrice
	if err := c.get(ctx, "/historical-price-eod/full", query, &prices); err != nil {
		return nil, err
	}
	return prices, nil
}
