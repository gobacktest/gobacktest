package internal

import (
	"time"

	"github.com/dirkolbrich/gobacktest/internal/utils"
)

// EventHandler declares the basic event interface
type EventHandler interface {
}

// BarEvent declares an event for an OHLCV bar (Open, High, Low, Close, Volume).
type BarEvent struct {
	Timestamp     time.Time
	Symbol        string
	OpenPrice     float64
	HighPrice     float64
	LowPrice      float64
	ClosePrice    float64
	AdjClosePrice float64
	Volume        int64
}

// SignalEvent declares a basic signal event
type SignalEvent struct {
	Timestamp    time.Time
	Symbol       string
	Direction    string // long or short
	SuggestedQty int64  // suggested quantitity
}

// OrderEvent declares a basic order event
type OrderEvent struct {
	Timestamp time.Time
	Symbol    string
	Direction string  // buy or sell
	Qty       int64   // quantity of the order
	OrderType string  // market or limit
	Limit     float64 // limit for the order
}

// FillEvent declares a basic fill event
type FillEvent struct {
	Timestamp   time.Time
	Symbol      string
	Exchange    string // exchange symbol
	Direction   string // buy or sell
	Qty         int64  // positive for buy, negativ for sell
	Price       float64
	Commission  float64
	ExchangeFee float64
	Cost        float64 // the total cost of the filled order incl commision and fees
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
	case (float64(f.Qty) * f.Price * comRate) < comMin:
		return comMin
	case (float64(f.Qty) * f.Price * comRate) > comMax:
		return comMax
	default:
		return utils.Round(float64(f.Qty)*f.Price*comRate, 3)
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
	return f.Commission + f.ExchangeFee
}
