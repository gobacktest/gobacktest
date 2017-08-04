package internal

import "time"

// ExecutionHandler is the basic interface for executing orders
type ExecutionHandler interface {
	ExecuteOrder(OrderEvent) (FillEvent, bool)
}

// Exchange is a basic execution handler implementation
type Exchange struct {
}

// ExecuteOrder executes an order event
func (e *Exchange) ExecuteOrder(o OrderEvent) (fill FillEvent, ok bool) {
	// log.Printf("Exchange receives Order: %#v \n", o)

	// simple implementation, creates a direct fill from the order
	// based on the last known closing price
	fill = FillEvent{
		Event:     Event{timestamp: time.Now(), symbol: o.Symbol()},
		Exchange:  "XETRA", // default Xetra exchange
		Direction: o.Direction,
		Qty:       o.Qty,
		Price:     10, // implement fetching last price from data handler
	}

	fill.Commission = fill.calculateComission()
	fill.ExchangeFee = fill.calculateExchangeFee()
	fill.Cost = fill.calculateCost()
	fill.Net = fill.calculateNet()

	return fill, true
}
