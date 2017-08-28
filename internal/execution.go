package internal

import "time"

// ExecutionHandler is the basic interface for executing orders
type ExecutionHandler interface {
	ExecuteOrder(orderEvent, DataEvent) (fillEvent, bool)
}

// Exchange is a basic execution handler implementation
type Exchange struct {
	Symbol      string
	ExchangeFee float64
}

// ExecuteOrder executes an order event
func (e *Exchange) ExecuteOrder(order orderEvent, data DataEvent) (fillEvent, bool) {
	// log.Printf("Exchange receives Order: %#v \n", order)

	// parse data event

	// simple implementation, creates a direct fill from the order
	// based on the last known closing price
	f := fillEvent{
		event:     event{timestamp: time.Now(), symbol: order.Symbol()},
		Exchange:  e.Symbol,
		Direction: order.Direction,
		Qty:       order.Qty,
		Price:     data.Close // implement fetching last price from data handler
	}

	f.Commission = e.calculateComission(float64(f.Qty), f.Price)
	f.ExchangeFee = e.calculateExchangeFee()
	f.Cost = e.calculateCost(f.Commission, f.ExchangeFee)
	f.Net = e.calculateNet(f.Direction, float64(f.Qty), f.Price, f.Cost)

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
		return utils.Round(qty*price + cost, 4)
	}
	// if "SLD"
	return utils.Round(qty*price - cost, 4)
}
