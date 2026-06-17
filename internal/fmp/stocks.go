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

// StockRatioRow is one raw ratio row returned by FMP.
type StockRatioRow map[string]any

// StockMetricRow is one raw key metrics row returned by FMP.
type StockMetricRow map[string]any

// StockPeer is one peer company row returned by FMP.
type StockPeer struct {
	Symbol      string  `json:"symbol"`
	CompanyName string  `json:"companyName"`
	Price       float64 `json:"price"`
	MarketCap   float64 `json:"mktCap"`
}

// StockRatingSnapshot is FMP's compact analyst rating score snapshot.
type StockRatingSnapshot struct {
	Symbol                  string  `json:"symbol"`
	Rating                  string  `json:"rating"`
	OverallScore            float64 `json:"overallScore"`
	DiscountedCashFlowScore float64 `json:"discountedCashFlowScore"`
	ReturnOnEquityScore     float64 `json:"returnOnEquityScore"`
	ReturnOnAssetsScore     float64 `json:"returnOnAssetsScore"`
	DebtToEquityScore       float64 `json:"debtToEquityScore"`
	PriceToEarningsScore    float64 `json:"priceToEarningsScore"`
	PriceToBookScore        float64 `json:"priceToBookScore"`
}

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

// StockRatiosTTM returns trailing-twelve-month stock ratios for a symbol.
func (c *Client) StockRatiosTTM(ctx context.Context, symbol string) ([]StockRatioRow, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	var ratios []StockRatioRow
	if err := c.get(ctx, "/ratios-ttm", url.Values{"symbol": []string{symbol}}, &ratios); err != nil {
		return nil, err
	}
	return ratios, nil
}

// StockKeyMetricsTTM returns trailing-twelve-month stock key metrics for a symbol.
func (c *Client) StockKeyMetricsTTM(ctx context.Context, symbol string) ([]StockMetricRow, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	var metrics []StockMetricRow
	if err := c.get(ctx, "/key-metrics-ttm", url.Values{"symbol": []string{symbol}}, &metrics); err != nil {
		return nil, err
	}
	return metrics, nil
}

// StockPeers returns peer companies for a stock symbol.
func (c *Client) StockPeers(ctx context.Context, symbol string) ([]StockPeer, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	var peers []StockPeer
	if err := c.get(ctx, "/stock-peers", url.Values{"symbol": []string{symbol}}, &peers); err != nil {
		return nil, err
	}
	return peers, nil
}

// StockRatingSnapshot returns FMP's compact analyst rating score snapshot.
func (c *Client) StockRatingSnapshot(ctx context.Context, symbol string) ([]StockRatingSnapshot, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	var ratings []StockRatingSnapshot
	if err := c.get(ctx, "/ratings-snapshot", url.Values{"symbol": []string{symbol}}, &ratings); err != nil {
		return nil, err
	}
	return ratings, nil
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
