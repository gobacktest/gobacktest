package internal

import (
	"time"
)

// EventHandler declares the basic event interface
type EventHandler interface {
}

// BarEvent declares an event for an OHLCV bar (Open, High, Low, Close, Volume).
type BarEvent struct {
	date       time.Time
	openPrice  int
	highPrice  int
	lowPrice   int
	closePrice int
	volume     int
}

// SignalEvent declares a basic signal event
type SignalEvent struct {
	timestamp    time.Time
	symbol       string
	direction    string // long or short
	suggestedQty int    // suggested quantitity
}

// OrderEvent declrares a basic order event
type OrderEvent struct {
	symbol    string
	orderType string // market or limit
	limit     int64  // limit for the order
	qty       int    // quantity of the order
	direction string // buy or sell
}

// FillEvent declrares a basic fill event
type FillEvent struct {
	timestamp   time.Time
	symbol      string
	exchange    string
	direction   string // buy or sell
	qty         int
	price       float64
	commission  float64
	exchangeFee float64
	cost        float64 // the total cost of the filled order incl commision and fees
}

// calculateComission() calculates the commission for a stock trade
//
// based on the conditions of IngDiba
// see https://www.ing-diba.de/wertpapiere/direkt-depot/konditionen
func (f FillEvent) calculateComission() float64 {
	var comMin = 9.90
	var comMax = 59.90
	var comRate = 0.0025 // in percent

	switch {
	case (float64(f.qty) * f.price * comRate) < comMin:
		return comMin
	case (float64(f.qty) * f.price * comRate) > comMax:
		return comMax
	default:
		return float64(f.qty) * f.price * comRate
	}
}

// calculateExchangeFee() calculates the exchange fee for a stock trade
//
// based on the conditions of IngDiba
// see https://www.ing-diba.de/wertpapiere/direkt-depot/konditionen
func (f FillEvent) calculateExchangeFee() float64 {
	return 1.75 // Xetra trade
}

// calculateCost() calculates the total cost for a stock trade
func (f FillEvent) calculateCost() float64 {
	return float64(f.qty)*f.price + f.commission + f.exchangeFee
}
