package gobacktest

import (
	"errors"
	"time"
)

// EventHandler declares the basic event interface
type EventHandler interface {
	Timer
	Symboler
}

// Timer declares the timer interface
type Timer interface {
	Time() time.Time
	SetTime(time.Time)
}

// Symboler declares the symboler interface
type Symboler interface {
	Symbol() string
	SetSymbol(string)
}

// Event is the implementation of the basic event interface.
type Event struct {
	timestamp time.Time
	symbol    string
}

// Time returns the timestamp of an event
func (e Event) Time() time.Time {
	return e.timestamp
}

// SetTime returns the timestamp of an event
func (e *Event) SetTime(t time.Time) {
	e.timestamp = t
}

// Symbol returns the symbol string of the event
func (e Event) Symbol() string {
	return e.symbol
}

// SetSymbol returns the symbol string of the event
func (e *Event) SetSymbol(s string) {
	e.symbol = s
}

// MetricHandler defines the handling of metrics to a data event
type MetricHandler interface {
	Add(string, float64) error
	Get(string) (float64, bool)
}

// Metric defines a metric property to a data point.
type Metric struct {
	metrics map[string]float64
}

// Add ads a value to the metrics map
func (m *Metric) Add(key string, value float64) error {
	// initialize map if nil
	if m.metrics == nil {
		m.metrics = map[string]float64{}
	}

	if key == "" {
		return errors.New("invalid key given")
	}

	m.metrics[key] = value
	return nil
}

// Get return a metric by name, if not found it returns false.
func (m Metric) Get(key string) (float64, bool) {
	value, ok := m.metrics[key]
	return value, ok
}

// Pricer defines the handling otf the latest Price Information
type Pricer interface {
	Price() float64
}

// DataEvent declares a data event interface
type DataEvent interface {
	EventHandler
	MetricHandler
	Pricer
}

// BarEvent declares a bar event interface.
type BarEvent interface {
	DataEvent
}

// Bar declares an event for an OHLCV bar (Open, High, Low, Close, Volume).
type Bar struct {
	Event
	Metric
	Open     float64
	High     float64
	Low      float64
	Close    float64
	AdjClose float64
	Volume   int64
}

// Price returns the close price of the bar event.
func (b Bar) Price() float64 {
	return b.Close
}

// TickEvent declares a tick event interface.
// type TickEvent interface {
// 	DataEvent
// }

// Tick declares an tick event
type Tick struct {
	Event
	Metric
	Bid       float64
	Ask       float64
	BidVolume int64
	AskVolume int64
}

// Price returns the middle of Bid and Ask.
func (t Tick) Price() float64 {
	latest := (t.Bid + t.Ask) / float64(2)
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
	direction OrderDirection // long or short
}

// Direction returns the Direction of a Signal
func (s Signal) Direction() OrderDirection {
	return s.direction
}

// SetDirection sets the Directions field of a Signal
func (s *Signal) SetDirection(dir OrderDirection) {
	s.direction = dir
}

// Directioner defines a direction interface
type Directioner interface {
	Direction() OrderDirection
	SetDirection(OrderDirection)
}

// Quantifier defines a qty interface
type Quantifier interface {
	Qty() int64
	SetQty(int64)
}

// OrderEvent declares the order event interface.
type OrderEvent interface {
	EventHandler
	Directioner
	Quantifier
	Status() OrderStatus
}

// FillEvent declares the fill event interface.
type FillEvent interface {
	EventHandler
	Directioner
	Quantifier
	Price() float64
	Commission() float64
	ExchangeFee() float64
	Cost() float64
	Value() float64
	NetValue() float64
}

// Fill declares a basic fill event
type Fill struct {
	Event
	direction   OrderDirection // BOT for buy, SLD for sell, HLD for hold
	Exchange    string         // exchange symbol
	qty         int64
	price       float64
	commission  float64
	exchangeFee float64
	cost        float64 // the total cost of the filled order incl commission and fees
}

// Direction returns the direction of a Fill
func (f Fill) Direction() OrderDirection {
	return f.direction
}

// SetDirection sets the Directions field of a Fill
func (f *Fill) SetDirection(dir OrderDirection) {
	f.direction = dir
}

// Qty returns the qty field of a fill
func (f Fill) Qty() int64 {
	return f.qty
}

// SetQty sets the Qty field of a Fill
func (f *Fill) SetQty(i int64) {
	f.qty = i
}

// Price returns the Price field of a fill
func (f Fill) Price() float64 {
	return f.price
}

// Commission returns the Commission field of a fill.
func (f Fill) Commission() float64 {
	return f.commission
}

// ExchangeFee returns the ExchangeFee Field of a fill
func (f Fill) ExchangeFee() float64 {
	return f.exchangeFee
}

// Cost returns the Cost field of a Fill
func (f Fill) Cost() float64 {
	return f.cost
}

// Value returns the value without cost.
func (f Fill) Value() float64 {
	value := float64(f.qty) * f.price
	return value
}

// NetValue returns the net value including cost.
func (f Fill) NetValue() float64 {
	if f.direction == BOT {
		// qty * price + cost
		netValue := float64(f.qty)*f.price + f.cost
		return netValue
	}
	// SLD
	//qty * price - cost
	netValue := float64(f.qty)*f.price - f.cost
	return netValue
}
