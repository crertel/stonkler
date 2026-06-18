package fmp

import "testing"

func TestParseStreamTradeShortFields(t *testing.T) {
	trade, ok, err := ParseStreamTrade([]byte(`{"s":"aapl","p":195.23,"q":12,"t":1710000000123}`))
	if err != nil {
		t.Fatalf("ParseStreamTrade() error = %v", err)
	}
	if !ok {
		t.Fatalf("ParseStreamTrade() ok = false, want true")
	}
	if trade.Symbol != "AAPL" {
		t.Fatalf("symbol = %q, want AAPL", trade.Symbol)
	}
	if trade.Price != 195.23 {
		t.Fatalf("price = %v, want 195.23", trade.Price)
	}
	if trade.Size != 12 {
		t.Fatalf("size = %v, want 12", trade.Size)
	}
	if trade.Timestamp != 1710000000123 {
		t.Fatalf("timestamp = %d, want 1710000000123", trade.Timestamp)
	}
}

func TestParseStreamTradeLongFields(t *testing.T) {
	trade, ok, err := ParseStreamTrade([]byte(`{"ticker":"MSFT","price":"412.5","size":"3","timestamp":"1710000000"}`))
	if err != nil {
		t.Fatalf("ParseStreamTrade() error = %v", err)
	}
	if !ok {
		t.Fatalf("ParseStreamTrade() ok = false, want true")
	}
	if trade.Symbol != "MSFT" {
		t.Fatalf("symbol = %q, want MSFT", trade.Symbol)
	}
	if trade.Price != 412.5 {
		t.Fatalf("price = %v, want 412.5", trade.Price)
	}
	if trade.Size != 3 {
		t.Fatalf("size = %v, want 3", trade.Size)
	}
}

func TestParseStreamTradeIgnoresControlMessage(t *testing.T) {
	_, ok, err := ParseStreamTrade([]byte(`{"event":"login","status":200}`))
	if err != nil {
		t.Fatalf("ParseStreamTrade() error = %v", err)
	}
	if ok {
		t.Fatalf("ParseStreamTrade() ok = true, want false")
	}
}

func TestParseStreamTradeRequiresPrice(t *testing.T) {
	_, ok, err := ParseStreamTrade([]byte(`{"s":"AAPL"}`))
	if err == nil {
		t.Fatalf("ParseStreamTrade() error = nil, want error")
	}
	if ok {
		t.Fatalf("ParseStreamTrade() ok = true, want false")
	}
}
