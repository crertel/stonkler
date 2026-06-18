package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/crertel/stonkler/internal/fmp"
)

const portfolioFileVersion = 1

type portfolioFile struct {
	Version     int                          `json:"version"`
	Stocks      map[string]portfolioPosition `json:"stocks,omitempty"`
	Funds       map[string]portfolioPosition `json:"funds,omitempty"`
	Crypto      map[string]portfolioPosition `json:"crypto,omitempty"`
	Forex       map[string]portfolioPosition `json:"forex,omitempty"`
	Commodities map[string]portfolioPosition `json:"commodities,omitempty"`
	Indexes     map[string]portfolioPosition `json:"indexes,omitempty"`
}

type portfolioPosition struct {
	Lots []portfolioLot `json:"lots"`
}

type portfolioLot struct {
	Basis      float64  `json:"basis"`
	Quantity   *float64 `json:"quantity,omitempty"`
	AcquiredOn string   `json:"acquired_on,omitempty"`
}

type basisBook struct {
	entries map[string]basisEntry
}

type basisEntry struct {
	Domain             string         `json:"domain"`
	Symbol             string         `json:"symbol"`
	Lots               []portfolioLot `json:"lots"`
	AverageBasis       float64        `json:"averageBasis,omitempty"`
	HasAverageBasis    bool           `json:"hasAverageBasis"`
	TotalQuantity      float64        `json:"totalQuantity,omitempty"`
	HasTotalQuantity   bool           `json:"hasTotalQuantity"`
	TotalCost          float64        `json:"totalCost,omitempty"`
	HasTotalCost       bool           `json:"hasTotalCost"`
	AcquiredOn         string         `json:"acquiredOn,omitempty"`
	MissingQuantityLot bool           `json:"missingQuantityLot,omitempty"`
}

type quoteWithBasis struct {
	Symbol             string  `json:"symbol"`
	Name               string  `json:"name,omitempty"`
	Price              float64 `json:"price"`
	Change             float64 `json:"change"`
	ChangePercentage   float64 `json:"changePercentage"`
	Volume             float64 `json:"volume"`
	MarketCap          float64 `json:"marketCap"`
	Timestamp          int64   `json:"timestamp"`
	Basis              float64 `json:"basis,omitempty"`
	HasBasis           bool    `json:"hasBasis"`
	Quantity           float64 `json:"quantity,omitempty"`
	HasQuantity        bool    `json:"hasQuantity"`
	Cost               float64 `json:"cost,omitempty"`
	HasCost            bool    `json:"hasCost"`
	MarketValue        float64 `json:"marketValue,omitempty"`
	HasMarketValue     bool    `json:"hasMarketValue"`
	UnrealizedGain     float64 `json:"unrealizedGain,omitempty"`
	HasUnrealizedGain  bool    `json:"hasUnrealizedGain"`
	UnrealizedGainPct  float64 `json:"unrealizedGainPercent,omitempty"`
	HasUnrealizedPct   bool    `json:"hasUnrealizedGainPercent"`
	AcquiredOn         string  `json:"acquiredOn,omitempty"`
	MissingQuantityLot bool    `json:"missingQuantityLot,omitempty"`
}

func loadBasisBook(path string) (*basisBook, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return decodeBasisBook(file)
}

func decodeBasisBook(r io.Reader) (*basisBook, error) {
	var portfolio portfolioFile
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&portfolio); err != nil {
		return nil, err
	}
	if portfolio.Version != portfolioFileVersion {
		return nil, fmt.Errorf("unsupported portfolio version %d", portfolio.Version)
	}

	book := &basisBook{entries: make(map[string]basisEntry)}
	if err := addBasisDomain(book, "stocks", portfolio.Stocks); err != nil {
		return nil, err
	}
	if err := addBasisDomain(book, "funds", portfolio.Funds); err != nil {
		return nil, err
	}
	if err := addBasisDomain(book, "crypto", portfolio.Crypto); err != nil {
		return nil, err
	}
	if err := addBasisDomain(book, "forex", portfolio.Forex); err != nil {
		return nil, err
	}
	if err := addBasisDomain(book, "commodities", portfolio.Commodities); err != nil {
		return nil, err
	}
	if err := addBasisDomain(book, "indexes", portfolio.Indexes); err != nil {
		return nil, err
	}
	return book, nil
}

func addBasisDomain(book *basisBook, domain string, positions map[string]portfolioPosition) error {
	for symbol, position := range positions {
		normalized := strings.ToUpper(strings.TrimSpace(symbol))
		if normalized == "" {
			return fmt.Errorf("%s contains an empty symbol", domain)
		}
		entry, err := summarizeBasisEntry(domain, normalized, position)
		if err != nil {
			return err
		}
		book.entries[basisKey(domain, normalized)] = entry
	}
	return nil
}

func summarizeBasisEntry(domain, symbol string, position portfolioPosition) (basisEntry, error) {
	if len(position.Lots) == 0 {
		return basisEntry{}, fmt.Errorf("%s %s requires at least one lot", domain, symbol)
	}

	entry := basisEntry{
		Domain: symbolDomain(domain),
		Symbol: symbol,
		Lots:   position.Lots,
	}
	var totalQuantity float64
	var totalCost float64
	var basisSum float64
	earliest := ""
	for index, lot := range position.Lots {
		if lot.Basis < 0 {
			return basisEntry{}, fmt.Errorf("%s %s lot %d has negative basis", domain, symbol, index+1)
		}
		basisSum += lot.Basis
		if lot.Quantity == nil {
			entry.MissingQuantityLot = true
		} else {
			if *lot.Quantity < 0 {
				return basisEntry{}, fmt.Errorf("%s %s lot %d has negative quantity", domain, symbol, index+1)
			}
			totalQuantity += *lot.Quantity
			totalCost += *lot.Quantity * lot.Basis
		}
		if lot.AcquiredOn != "" {
			if _, err := time.Parse("2006-01-02", lot.AcquiredOn); err != nil {
				return basisEntry{}, fmt.Errorf("%s %s lot %d has invalid acquired_on %q", domain, symbol, index+1, lot.AcquiredOn)
			}
			if earliest == "" || lot.AcquiredOn < earliest {
				earliest = lot.AcquiredOn
			}
		}
	}

	entry.AcquiredOn = earliest
	if !entry.MissingQuantityLot {
		entry.HasTotalQuantity = true
		entry.HasTotalCost = true
		entry.TotalQuantity = totalQuantity
		entry.TotalCost = totalCost
		if totalQuantity > 0 {
			entry.HasAverageBasis = true
			entry.AverageBasis = totalCost / totalQuantity
		}
	} else {
		entry.HasAverageBasis = true
		entry.AverageBasis = basisSum / float64(len(position.Lots))
	}
	return entry, nil
}

func (b *basisBook) Entry(domain, symbol string) (basisEntry, bool) {
	if b == nil {
		return basisEntry{}, false
	}
	entry, ok := b.entries[basisKey(domain, strings.ToUpper(strings.TrimSpace(symbol)))]
	return entry, ok
}

func (b *basisBook) Entries() []basisEntry {
	if b == nil {
		return nil
	}
	entries := make([]basisEntry, 0, len(b.entries))
	for _, entry := range b.entries {
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Domain == entries[j].Domain {
			return entries[i].Symbol < entries[j].Symbol
		}
		return entries[i].Domain < entries[j].Domain
	})
	return entries
}

func basisKey(domain, symbol string) string {
	return symbolDomain(domain) + ":" + strings.ToUpper(strings.TrimSpace(symbol))
}

func symbolDomain(domain string) string {
	switch domain {
	case "stock":
		return "stocks"
	case "fund":
		return "funds"
	case "commodity":
		return "commodities"
	case "index":
		return "indexes"
	default:
		return domain
	}
}

func attachBasis(domain string, quotes []fmp.Quote, book *basisBook) []quoteWithBasis {
	rows := make([]quoteWithBasis, 0, len(quotes))
	for _, quote := range quotes {
		row := quoteWithBasis{
			Symbol:           quote.Symbol,
			Name:             quote.Name,
			Price:            quote.Price,
			Change:           quote.Change,
			ChangePercentage: quote.ChangePercentage,
			Volume:           quote.Volume,
			MarketCap:        quote.MarketCap,
			Timestamp:        quote.Timestamp,
		}
		if entry, ok := book.Entry(domain, row.Symbol); ok {
			row.HasBasis = entry.HasAverageBasis
			row.Basis = entry.AverageBasis
			row.HasQuantity = entry.HasTotalQuantity
			row.Quantity = entry.TotalQuantity
			row.HasCost = entry.HasTotalCost
			row.Cost = entry.TotalCost
			row.AcquiredOn = entry.AcquiredOn
			row.MissingQuantityLot = entry.MissingQuantityLot
			if entry.HasTotalQuantity {
				row.HasMarketValue = true
				row.MarketValue = row.Price * entry.TotalQuantity
			}
			if entry.HasTotalCost && row.HasMarketValue {
				row.HasUnrealizedGain = true
				row.UnrealizedGain = row.MarketValue - entry.TotalCost
				if entry.TotalCost != 0 {
					row.HasUnrealizedPct = true
					row.UnrealizedGainPct = row.UnrealizedGain / entry.TotalCost * 100
				}
			}
		}
		rows = append(rows, row)
	}
	return rows
}
