package backtest

import (
// "fmt"
)

// ExecutionHandler is the basic interface for executing orders
type ExecutionHandler interface {
	ExecuteOrder(OrderEvent, DataHandler) (*Fill, error)
}

// Exchange is a basic execution handler implementation
type Exchange struct {
	Symbol      string
	Commission  CommissionHandler
	ExchangeFee float64
}

// ExecuteOrder executes an order event
func (e *Exchange) ExecuteOrder(order OrderEvent, data DataHandler) (*Fill, error) {
	// fetch latest known data event for the symbol
	latest := data.Latest(order.GetSymbol())

	// simple implementation, creates a direct fill from the order
	// based on the last known data price
	f := &Fill{
		Event:    Event{Timestamp: order.GetTime(), Symbol: order.GetSymbol()},
		Exchange: e.Symbol,
		Qty:      order.GetQty(),
		Price:    latest.LatestPrice(), // last price from data event
	}

	switch order.GetDirection() {
	case "buy":
		f.Direction = "BOT"
	case "sell":
		f.Direction = "SLD"
	}

	commission, err := e.Commission.Calculate(float64(f.Qty), f.Price)
	if err != nil {
		return f, err
	}
	f.Commission = commission

	f.ExchangeFee = e.calculateExchangeFee()
	f.Cost = e.calculateCost(commission, f.ExchangeFee)

	return f, nil
}

// calculateExchangeFee() calculates the exchange fee for a stock trade
func (e *Exchange) calculateExchangeFee() float64 {
	return e.ExchangeFee
}

// calculateCost() calculates the total cost for a stock trade
func (e *Exchange) calculateCost(commission, fee float64) float64 {
	return commission + fee
}
