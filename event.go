package gobacktest

import (
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

// SignalEvent declares the signal event interface.
type SignalEvent interface {
	EventHandler
	Directioner
}

// Directioner defines a direction interface
type Directioner interface {
	Direction() Direction
	SetDirection(Direction)
}

// OrderEvent declares the order event interface.
type OrderEvent interface {
	EventHandler
	Directioner
	Quantifier
	IDer
	Status() OrderStatus
	Limit() float64
	Stop() float64
}

// Quantifier defines a qty interface.
type Quantifier interface {
	Qty() int64
	SetQty(int64)
}

// IDer declares setting and retrieving of an Id.
type IDer interface {
	ID() int
	SetID(int)
}

// FillEvent declares fill event functionality.
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
