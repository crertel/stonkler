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

// ETFSectorWeighting is one sector allocation row returned by FMP for an ETF or fund.
type ETFSectorWeighting struct {
	Symbol           string  `json:"symbol"`
	Sector           string  `json:"sector"`
	WeightPercentage float64 `json:"weightPercentage"`
}

// ETFCountryWeighting is one country allocation row returned by FMP for an ETF or fund.
type ETFCountryWeighting struct {
	Country          string `json:"country"`
	WeightPercentage string `json:"weightPercentage"`
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

// ETFSectorWeightings returns sector allocation weights for an ETF or fund symbol.
func (c *Client) ETFSectorWeightings(ctx context.Context, symbol string) ([]ETFSectorWeighting, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	var weightings []ETFSectorWeighting
	if err := c.get(ctx, "/etf/sector-weightings", url.Values{"symbol": []string{symbol}}, &weightings); err != nil {
		return nil, err
	}
	return weightings, nil
}

// ETFCountryWeightings returns country allocation weights for an ETF or fund symbol.
func (c *Client) ETFCountryWeightings(ctx context.Context, symbol string) ([]ETFCountryWeighting, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	var weightings []ETFCountryWeighting
	if err := c.get(ctx, "/etf/country-weightings", url.Values{"symbol": []string{symbol}}, &weightings); err != nil {
		return nil, err
	}
	return weightings, nil
}
