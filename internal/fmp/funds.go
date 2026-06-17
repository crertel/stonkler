package fmp

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// FundProfile returns profile data for an ETF or fund symbol.
func (c *Client) FundProfile(ctx context.Context, symbol string) (StockProfile, error) {
	return c.StockProfile(ctx, symbol)
}

// ETFHolding is one holding row returned by FMP for an ETF.
type ETFHolding struct {
	Asset            string  `json:"asset"`
	Name             string  `json:"name"`
	ISIN             string  `json:"isin,omitempty"`
	CUSIP            string  `json:"cusip,omitempty"`
	SharesNumber     float64 `json:"sharesNumber"`
	WeightPercentage float64 `json:"weightPercentage"`
	MarketValue      float64 `json:"marketValue"`
	Updated          string  `json:"updated,omitempty"`
}

// ETFHoldings returns holdings for an ETF symbol.
func (c *Client) ETFHoldings(ctx context.Context, symbol string) ([]ETFHolding, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	var holdings []ETFHolding
	if err := c.getV3(ctx, "/etf-holder/"+url.PathEscape(symbol), url.Values{}, &holdings); err != nil {
		return nil, err
	}
	return holdings, nil
}
