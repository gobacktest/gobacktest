package backtest

import (
	"time"

	"github.com/shopspring/decimal"
)

// EventHandler declares the basic event interface
type EventHandler interface {
	Timer
	Symboler
}

// Timer declares the timer interface
type Timer interface {
	GetTime() time.Time
}

// Symboler declares the symboler interface
type Symboler interface {
	GetSymbol() string
}

// Event is the implementation of the basic event interface.
type Event struct {
	Timestamp time.Time
	Symbol    string
}

// GetTime returns the timestamp of an event
func (e Event) GetTime() time.Time {
	return e.Timestamp
}

// GetSymbol returns the symbol string of the event
func (e Event) GetSymbol() string {
	return e.Symbol
}

// DataEventHandler declares a data event interface
type DataEventHandler interface {
	EventHandler
	LatestPrice() float64
}

// DataEvent is the basic implementation of a data event handler.
type DataEvent struct {
	Metrics map[string]float64
}

// BarEvent declares a bar event interface.
type BarEvent interface {
	DataEventHandler
}

// Bar declares an event for an OHLCV bar (Open, High, Low, Close, Volume).
type Bar struct {
	Event
	DataEvent
	Open     float64
	High     float64
	Low      float64
	Close    float64
	AdjClose float64
	Volume   int64
}

// LatestPrice returns the close proce of the bar event.
func (b Bar) LatestPrice() float64 {
	return b.Close
}

// TickEvent declares a tick event interface.
type TickEvent interface {
	DataEventHandler
}

// Tick declares an tick event
type Tick struct {
	Event
	DataEvent
	Bid float64
	Ask float64
}

// LatestPrice returns the middle of Bid and Ask.
func (t Tick) LatestPrice() float64 {
	bid := decimal.NewFromFloat(t.Bid)
	ask := decimal.NewFromFloat(t.Ask)
	diff := decimal.New(2, 0)
	latest, _ := bid.Add(ask).Div(diff).Round(DP).Float64()
	return latest
}

// SignalEvent declares the signal event interface.
type SignalEvent interface {
	EventHandler
	Directioner
}

// Signal declares a basic signal event
type Signal struct {
	Event
	Direction string // long or short
}

// SetDirection sets the Directions field of a Signal
func (s *Signal) SetDirection(st string) {
	s.Direction = st
}

// GetDirection returns the Direction of a Signal
func (s Signal) GetDirection() string {
	return s.Direction
}

// OrderEvent declares the order event interface.
type OrderEvent interface {
	EventHandler
	Directioner
	Quantifier
}

// Directioner defines a direction interface
type Directioner interface {
	SetDirection(string)
	GetDirection() string
}

// Quantifier defines a qty interface
type Quantifier interface {
	SetQty(int64)
	GetQty() int64
}

// Order declares a basic order event
type Order struct {
	Event
	Direction string  // buy or sell
	Qty       int64   // quantity of the order
	OrderType string  // market or limit
	Limit     float64 // limit for the order
}

// SetDirection sets the Directions field of an Order
func (o *Order) SetDirection(s string) {
	o.Direction = s
}

// GetDirection returns the Direction of an Order
func (o Order) GetDirection() string {
	return o.Direction
}

// SetQty sets the Qty field of an Order
func (o *Order) SetQty(i int64) {
	o.Qty = i
}

// GetQty returns the Qty field of an Order
func (o Order) GetQty() int64 {
	return o.Qty
}

// FillEvent declares the fill event interface.
type FillEvent interface {
	EventHandler
	Directioner
	Quantifier
	GetPrice() float64
	GetCommission() float64
	GetExchangeFee() float64
	GetCost() float64
	Value() float64
	NetValue() float64
}

// Fill declares a basic fill event
type Fill struct {
	Event
	Exchange    string // exchange symbol
	Direction   string // BOT for buy or SLD for sell
	Qty         int64
	Price       float64
	Commission  float64
	ExchangeFee float64
	Cost        float64 // the total cost of the filled order incl commission and fees
}

// SetDirection sets the Directions field of a Fill
func (f *Fill) SetDirection(s string) {
	f.Direction = s
}

// GetDirection returns the direction of a Fill
func (f Fill) GetDirection() string {
	return f.Direction
}

// SetQty sets the Qty field of a Fill
func (f *Fill) SetQty(i int64) {
	f.Qty = i
}

// GetQty returns the qty field of a fill
func (f Fill) GetQty() int64 {
	return f.Qty
}

// GetPrice returns the Price field of a fill
func (f Fill) GetPrice() float64 {
	return f.Price
}

// GetCommission returns the Commission field of a fill.
func (f Fill) GetCommission() float64 {
	return f.Commission
}

// GetExchangeFee returns the ExchangeFee Field of a fill
func (f Fill) GetExchangeFee() float64 {
	return f.ExchangeFee
}

// GetCost returns the Cost field of a Fill
func (f Fill) GetCost() float64 {
	return f.Cost
}

// Value returns the value without cost.
func (f Fill) Value() float64 {
	qty := decimal.New(f.Qty, 0)
	price := decimal.NewFromFloat(f.Price)
	value, _ := qty.Mul(price).Round(DP).Float64()
	return value
}

// NetValue returns the net value including cost.
func (f Fill) NetValue() float64 {
	qty := decimal.New(f.Qty, 0)
	price := decimal.NewFromFloat(f.Price)
	cost := decimal.NewFromFloat(f.Cost)

	if f.Direction == "BOT" {
		// qty * price + cost
		netValue, _ := qty.Mul(price).Add(cost).Round(DP).Float64()
		return netValue
	}
	// SLD
	//qty * price - cost
	netValue, _ := qty.Mul(price).Sub(cost).Round(DP).Float64()
	return netValue
}
