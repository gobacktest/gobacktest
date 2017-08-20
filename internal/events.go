package internal

import (
	"time"
)

// EventHandler declares the basic event interface
type EventHandler interface {
	Timestamp() time.Time
	Symbol() string
}

// Event is a basic impementation of an event.
type Event struct {
	timestamp time.Time
	symbol    string
}

// Timestamp returns the time property of an event.
func (e Event) Timestamp() time.Time {
	return e.timestamp
}

// Symbol returns the symbol property of an event.
func (e Event) Symbol() string {
	return e.symbol
}

// BarEvent declares an event for an OHLCV bar (Open, High, Low, Close, Volume).
type BarEvent struct {
	Event
	OpenPrice     float64
	HighPrice     float64
	LowPrice      float64
	ClosePrice    float64
	AdjClosePrice float64
	Volume        int64
	metrics       map[string]float64
}

// SignalEvent declares a basic signal event
type SignalEvent struct {
	Event
	Direction    string // long or short
}

// OrderEvent declares a basic order event
type OrderEvent struct {
	Event
	Direction string  // buy or sell
	Qty       int64   // quantity of the order
	OrderType string  // market or limit
	Limit     float64 // limit for the order
}

// FillEvent declares a basic fill event
type FillEvent struct {
	Event
	Exchange    string // exchange symbol
	Direction   string // BOT for buy or SLD for sell
	Qty         int64
	Price       float64
	Commission  float64
	ExchangeFee float64
	Cost        float64 // the total cost of the filled order incl commision and fees
	Net         float64 // the net value of the filled order e.g. spend/taken incl expenses
}
