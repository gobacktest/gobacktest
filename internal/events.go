package internal

import (
	"time"
)

// Event declares the basic event interface
type Event interface {
	Timestamp() time.Time
	Symbol() string
}

// event is the implementation of the basic event interface.
type event struct {
	timestamp time.Time
	symbol    string
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
	Event
	LatestPrice() float64
	Metrics() map[string]float64
}

// dataEvent is the implementation of the basic DataEvent.
type dataEvent struct {
	metrics map[string]float64
}

// Metrics returns the a map of metrics
func (d dataEvent) Metrics() map[string]float64 {
	return d.metrics
}

// BarEvent declares a bar event interface.
type BarEvent interface {
	DataEvent
	Open() float64
	High() float64
	Low() float64
	Close() float64
	AdjClose() float64
	Volume() int64
}

// barEvent declares an event for an OHLCV bar (Open, High, Low, Close, Volume).
type bar struct {
	event
	dataEvent
	openPrice     float64
	highPrice     float64
	lowPrice      float64
	closePrice    float64
	adjClosePrice float64
	volume        int64
}

func (b bar) Open() float64 {
	return b.openPrice
}

func (b bar) High() float64 {
	return b.highPrice
}

func (b bar) Low() float64 {
	return b.lowPrice
}

func (b bar) Close() float64 {
	return b.closePrice
}

func (b bar) AdjClose() float64 {
	return b.adjClosePrice
}

func (b bar) Volume() int64 {
	return b.volume
}

func (b bar) LatestPrice() float64 {
	return b.closePrice
}

// TickEvent declares a tick event interface.
type TickEvent interface {
	DataEvent
	Bid() float64
	Ask() float64
}

// tick declares an tick event
type tick struct {
	event
	dataEvent
	bidPrice float64
	askPrice float64
}

func (t tick) LatestPrice() float64 {
	return (t.bidPrice + t.askPrice) / 2
}

func (t tick) Bid() float64 {
	return t.bidPrice
}

func (t tick) Ask() float64 {
	return t.askPrice
}

/***** BackTest Event Interfaces *****/

// Directioner defines the direction interface
type Directioner interface {
	SetDirection(string)
	Direction() string
}

// Qtyer defines the Quantity interface
type Qtyer interface {
	SetQty(int64)
	Qty() int64
}

// Limiter defines the limit interface
type Limiter interface {
	SetLimit(float64)
	Limit() float64
}

/***** Backtest event interfaces and concrete struct implementations *****/

// SignalEvent declares the signal event interface.
type SignalEvent interface {
	Event
	Directioner
	IsSignal() bool
}

// signal declares a basic signal event
type signal struct {
	event
	direction string // long or short
}

func (s signal) IsSignal() bool {
	return true
}

func (s *signal) SetDirection(dir string) {
	s.direction = dir
}

func (s signal) Direction() string {
	return s.direction
}

// OrderEvent declares the order event interface.
type OrderEvent interface {
	Event
	Directioner
	Qtyer
	IsOrder() bool
}

// orderEvent declares a basic order event
type order struct {
	event
	direction string  // buy or sell
	qty       int64   // quantity of the order
	orderType string  // market or limit
	limit     float64 // limit for the order
}

func (o order) IsOrder() bool {
	return true
}

func (o *order) SetDirection(dir string) {
	o.direction = dir
}

func (o order) Direction() string {
	return o.direction
}

func (o *order) SetQty(qty int64) {
	o.qty = qty
}

func (o order) Qty() int64 {
	return o.qty
}

func (o *order) SetOrderType(ot string) {
	o.orderType = ot
}

func (o order) OrderType() string {
	return o.orderType
}

func (o *order) SetLimit(l float64) {
	o.limit = l
}

func (o order) Limit() float64 {
	return o.limit
}

// FillEvent declares the fill event interface.
type FillEvent interface {
	Event
	Directioner
	Qtyer
	IsFill() bool
	Price() float64
	Commission() float64
	ExchangeFee() float64
	Cost() float64
	Net() float64
}

// fillEvent declares a basic fill event
type fill struct {
	event
	exchange    string // exchange symbol
	direction   string // BOT for buy or SLD for sell
	qty         int64
	price       float64
	commission  float64
	exchangeFee float64
	cost        float64 // the total cost of the filled order incl commission and fees
	net         float64 // the net value of the filled order e.g. spend/taken incl expenses
}

func (f fill) IsFill() bool {
	return true
}

func (f *fill) SetDirection(dir string) {
	f.direction = dir
}

func (f fill) Direction() string {
	return f.direction
}

func (f *fill) SetQty(qty int64) {
	f.qty = qty
}

func (f fill) Qty() int64 {
	return f.qty
}

func (f fill) Price() float64 {
	return f.price
}

func (f fill) Commission() float64 {
	return f.commission
}

func (f fill) ExchangeFee() float64 {
	return f.exchangeFee
}

func (f fill) Cost() float64 {
	return f.cost
}

func (f fill) Net() float64 {
	return f.net
}
