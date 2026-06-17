package fmp

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// StockProfile is normalized company profile data from FMP.
type StockProfile struct {
	Symbol            string  `json:"symbol"`
	CompanyName       string  `json:"companyName"`
	Price             float64 `json:"price"`
	MarketCap         float64 `json:"marketCap"`
	Beta              float64 `json:"beta,omitempty"`
	LastDividend      float64 `json:"lastDividend,omitempty"`
	Range             string  `json:"range,omitempty"`
	Change            float64 `json:"change,omitempty"`
	ChangePercentage  float64 `json:"changePercentage,omitempty"`
	Volume            float64 `json:"volume,omitempty"`
	AverageVolume     float64 `json:"averageVolume,omitempty"`
	Currency          string  `json:"currency,omitempty"`
	CIK               string  `json:"cik,omitempty"`
	ISIN              string  `json:"isin,omitempty"`
	CUSIP             string  `json:"cusip,omitempty"`
	ExchangeFullName  string  `json:"exchangeFullName,omitempty"`
	Exchange          string  `json:"exchange,omitempty"`
	Industry          string  `json:"industry,omitempty"`
	Website           string  `json:"website,omitempty"`
	Description       string  `json:"description,omitempty"`
	CEO               string  `json:"ceo,omitempty"`
	Sector            string  `json:"sector,omitempty"`
	Country           string  `json:"country,omitempty"`
	FullTimeEmployees string  `json:"fullTimeEmployees,omitempty"`
	Phone             string  `json:"phone,omitempty"`
	Address           string  `json:"address,omitempty"`
	City              string  `json:"city,omitempty"`
	State             string  `json:"state,omitempty"`
	ZIP               string  `json:"zip,omitempty"`
	Image             string  `json:"image,omitempty"`
	IPODate           string  `json:"ipoDate,omitempty"`
	IsETF             bool    `json:"isEtf,omitempty"`
	IsActivelyTrading bool    `json:"isActivelyTrading,omitempty"`
	IsADR             bool    `json:"isAdr,omitempty"`
	IsFund            bool    `json:"isFund,omitempty"`
}

// StockProfile returns profile data for one stock symbol.
func (c *Client) StockProfile(ctx context.Context, symbol string) (StockProfile, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return StockProfile{}, fmt.Errorf("symbol is required")
	}

	var profiles []StockProfile
	query := url.Values{}
	query.Set("symbol", symbol)
	if err := c.get(ctx, "/profile", query, &profiles); err != nil {
		return StockProfile{}, err
	}
	if len(profiles) == 0 {
		return StockProfile{}, fmt.Errorf("no profile returned for %s", symbol)
	}
	return profiles[0], nil
}
