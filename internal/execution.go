package internal

import (
	"time"

	"github.com/dirkolbrich/gobacktest/internal/utils"
)

// ExecutionHandler is the basic interface for executing orders
type ExecutionHandler interface {
	ExecuteOrder(OrderEvent, DataEvent) (FillEvent, bool)
}

// Exchange is a basic execution handler implementation
type Exchange struct {
	Symbol      string
	ExchangeFee float64
}

// ExecuteOrder executes an order event
func (e *Exchange) ExecuteOrder(order OrderEvent, data DataEvent) (FillEvent, bool) {
	// log.Printf("Exchange receives Order: %#v \n", order)

	var latestPrice float64

	// parse data event
	switch data := data.(type) {
	case BarEvent:
		latestPrice = data.Close()
	case TickEvent:
		latestPrice = (data.Bid() + data.Ask()) / 2
	}

	// simple implementation, creates a direct fill from the order
	// based on the last known closing price
	f := &fill{
		event:     event{timestamp: time.Now(), symbol: order.Symbol()},
		exchange:  e.Symbol,
		direction: order.Direction(),
		qty:       order.Qty(),
		price:     latestPrice, // implement fetching last price from data handler
	}

	f.commission = e.calculateComission(float64(f.qty), f.price)
	f.exchangeFee = e.calculateExchangeFee()
	f.cost = e.calculateCost(f.commission, f.exchangeFee)
	f.net = e.calculateNet(f.direction, float64(f.qty), f.price, f.cost)

	return f, true
}

// calculateComission() calculates the commission for a stock trade
//
// based on the conditions of IngDiba
// see https://www.ing-diba.de/wertpapiere/direkt-depot/konditionen
func (e *Exchange) calculateComission(qty, price float64) float64 {
	var comMin = 9.90
	var comMax = 59.90
	var comRate = 0.0025 // in percent

	switch {
	case (qty * price * comRate) < comMin:
		return comMin
	case (qty * price * comRate) > comMax:
		return comMax
	default:
		return utils.Round(qty*price*comRate, 4)
	}
}

// calculateExchangeFee() calculates the exchange fee for a stock trade
func (e *Exchange) calculateExchangeFee() float64 {
	return e.ExchangeFee
}

// calculateCost() calculates the total cost for a stock trade
func (e *Exchange) calculateCost(commission, fee float64) float64 {
	return commission + fee
}

// calculateCost() calculates the total cost for a stock trade
func (e *Exchange) calculateNet(dir string, qty, price, cost float64) float64 {
	if dir == "BOT" {
		return utils.Round(qty*price+cost, 4)
	}
	// if "SLD"
	return utils.Round(qty*price-cost, 4)
}
