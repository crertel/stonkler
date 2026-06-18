package cli

import (
	"strings"
	"testing"
)

func TestDecodeBasisBookSummarizesWeightedLots(t *testing.T) {
	book, err := decodeBasisBook(strings.NewReader(`{
		"version": 1,
		"stocks": {
			"AAPL": {
				"lots": [
					{"basis": 100, "quantity": 2, "acquired_on": "2024-02-01"},
					{"basis": 130, "quantity": 1, "acquired_on": "2024-01-01"}
				]
			}
		}
	}`))
	if err != nil {
		t.Fatalf("decodeBasisBook() error = %v", err)
	}

	entry, ok := book.Entry("stocks", "aapl")
	if !ok {
		t.Fatalf("Entry() ok = false, want true")
	}
	if !entry.HasAverageBasis || entry.AverageBasis != 110 {
		t.Fatalf("average basis = %v/%t, want 110/true", entry.AverageBasis, entry.HasAverageBasis)
	}
	if !entry.HasTotalQuantity || entry.TotalQuantity != 3 {
		t.Fatalf("quantity = %v/%t, want 3/true", entry.TotalQuantity, entry.HasTotalQuantity)
	}
	if !entry.HasTotalCost || entry.TotalCost != 330 {
		t.Fatalf("cost = %v/%t, want 330/true", entry.TotalCost, entry.HasTotalCost)
	}
	if entry.AcquiredOn != "2024-01-01" {
		t.Fatalf("acquired = %q, want 2024-01-01", entry.AcquiredOn)
	}
}

func TestDecodeBasisBookAllowsMissingQuantity(t *testing.T) {
	book, err := decodeBasisBook(strings.NewReader(`{
		"version": 1,
		"crypto": {
			"BTCUSD": {
				"lots": [
					{"basis": 60000},
					{"basis": 70000}
				]
			}
		}
	}`))
	if err != nil {
		t.Fatalf("decodeBasisBook() error = %v", err)
	}

	entry, ok := book.Entry("crypto", "BTCUSD")
	if !ok {
		t.Fatalf("Entry() ok = false, want true")
	}
	if !entry.HasAverageBasis || entry.AverageBasis != 65000 {
		t.Fatalf("average basis = %v/%t, want 65000/true", entry.AverageBasis, entry.HasAverageBasis)
	}
	if entry.HasTotalQuantity {
		t.Fatalf("HasTotalQuantity = true, want false")
	}
	if !entry.MissingQuantityLot {
		t.Fatalf("MissingQuantityLot = false, want true")
	}
}

func TestDecodeBasisBookRejectsInvalidAcquiredDate(t *testing.T) {
	_, err := decodeBasisBook(strings.NewReader(`{
		"version": 1,
		"stocks": {
			"AAPL": {"lots": [{"basis": 100, "acquired_on": "2024/01/01"}]}
		}
	}`))
	if err == nil {
		t.Fatal("decodeBasisBook() error = nil, want error")
	}
}
