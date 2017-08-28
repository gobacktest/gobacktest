package internal

import (
	"time"
)

// Event declares the basic event interface
type Event interface {
	Timestamp() time.Time
	Symbol() string
}

// Timestamp returns the time property of an event.
func (e event) Timestamp() time.Time {
	return e.timestamp
}

// Symbol returns the symbol property of an event.
func (e event) Symbol() string {
	return e.symbol
}

// DataEvent declares a data event interface
type DataEvent interface {
	Metrics() map[string]float64
	Metric(string) float64
}

// dataEvent is the implementation of the basic DataEvent.
type dataEvent struct {
	metrics map[string]float64
}

// Metrics returns the a map of metrics
func (d dataEvent) Metrics() map[string]float64 {
	return d.metrics
}

// Metric() returns the specific metrics by key
func (d dataEvent) Metrics(s string) (metric float64, ok bool) {
	return d.metric[s]
}

// barEvent declares an event for an OHLCV bar (Open, High, Low, Close, Volume).
type barEvent struct {
	event
	dataEvent
	OpenPrice     float64
	HighPrice     float64
	LowPrice      float64
	ClosePrice    float64
	AdjClosePrice float64
	Volume        int64
}

// signalEvent declares a basic signal event
type signalEvent struct {
	event
	Direction    string // long or short
}

// orderEvent declares a basic order event
type orderEvent struct {
	event
	Direction string  // buy or sell
	Qty       int64   // quantity of the order
	OrderType string  // market or limit
	Limit     float64 // limit for the order
}

// fillEvent declares a basic fill event
type fillEvent struct {
	event
	Exchange    string // exchange symbol
	Direction   string // BOT for buy or SLD for sell
	Qty         int64
	Price       float64
	Commission  float64
	ExchangeFee float64
	Cost        float64 // the total cost of the filled order incl commision and fees
	Net         float64 // the net value of the filled order e.g. spend/taken incl expenses
}
