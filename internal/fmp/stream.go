package fmp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const defaultWebSocketURL = "wss://websockets.financialmodelingprep.com"

// StreamTrade is a real-time trade update from FMP's stock websocket feed.
type StreamTrade struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Size      float64 `json:"size,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

// StreamStockTrades subscribes to FMP's real-time stock trade websocket feed.
func (c *Client) StreamStockTrades(ctx context.Context, symbols []string) (<-chan StreamTrade, <-chan error) {
	trades := make(chan StreamTrade)
	errs := make(chan error, 1)

	go func() {
		defer close(trades)
		defer close(errs)

		endpoint, err := url.Parse(defaultWebSocketURL)
		if err != nil {
			errs <- err
			return
		}

		conn, _, err := websocket.DefaultDialer.DialContext(ctx, endpoint.String(), nil)
		if err != nil {
			errs <- err
			return
		}
		defer conn.Close()

		done := make(chan struct{})
		defer close(done)
		go func() {
			select {
			case <-ctx.Done():
				_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
				_ = conn.Close()
			case <-done:
			}
		}()

		if err := writeFMPStreamMessage(conn, "login", map[string]string{"apiKey": c.apiKey}); err != nil {
			errs <- err
			return
		}
		for _, symbol := range symbols {
			ticker := strings.ToUpper(strings.TrimSpace(symbol))
			if ticker == "" {
				continue
			}
			if err := writeFMPStreamMessage(conn, "subscribe", map[string]string{"ticker": ticker}); err != nil {
				errs <- err
				return
			}
		}

		for {
			_, payload, err := conn.ReadMessage()
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				errs <- err
				return
			}

			trade, ok, err := ParseStreamTrade(payload)
			if err != nil {
				errs <- err
				return
			}
			if !ok {
				continue
			}

			select {
			case trades <- trade:
			case <-ctx.Done():
				return
			}
		}
	}()

	return trades, errs
}

func writeFMPStreamMessage(conn *websocket.Conn, event string, data any) error {
	return conn.WriteJSON(map[string]any{
		"event": event,
		"data":  data,
	})
}

// ParseStreamTrade decodes a stock websocket trade payload. FMP sends control
// acknowledgements on the same socket; those return ok=false.
func ParseStreamTrade(payload []byte) (StreamTrade, bool, error) {
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()

	var message map[string]any
	if err := dec.Decode(&message); err != nil {
		return StreamTrade{}, false, err
	}

	symbol := firstString(message, "s", "symbol", "ticker")
	if symbol == "" {
		return StreamTrade{}, false, nil
	}

	price, ok := firstFloat(message, "p", "price")
	if !ok {
		return StreamTrade{}, false, fmt.Errorf("stream trade for %s missing price", symbol)
	}

	size, _ := firstFloat(message, "q", "size", "volume")
	timestamp, _ := firstInt(message, "t", "timestamp")

	return StreamTrade{
		Symbol:    strings.ToUpper(symbol),
		Price:     price,
		Size:      size,
		Timestamp: timestamp,
	}, true, nil
}

func firstString(values map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := values[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case string:
			return strings.TrimSpace(typed)
		}
	}
	return ""
}

func firstFloat(values map[string]any, keys ...string) (float64, bool) {
	for _, key := range keys {
		value, ok := values[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case json.Number:
			parsed, err := typed.Float64()
			return parsed, err == nil
		case float64:
			return typed, true
		case string:
			parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
			return parsed, err == nil
		}
	}
	return 0, false
}

func firstInt(values map[string]any, keys ...string) (int64, bool) {
	for _, key := range keys {
		value, ok := values[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case json.Number:
			parsed, err := typed.Int64()
			return parsed, err == nil
		case float64:
			return int64(typed), true
		case string:
			parsed, err := strconv.ParseInt(strings.TrimSpace(typed), 10, 64)
			return parsed, err == nil
		}
	}
	return 0, false
}
