package backtest

import "github.com/dirkolbrich/gobacktest/pkg/utils"

// ExecutionHandler is the basic interface for executing orders
type ExecutionHandler interface {
	ExecuteOrder(OrderEvent, DataHandler) (FillEvent, error)
}

// Exchange is a basic execution handler implementation
type Exchange struct {
	Symbol      string
	ExchangeFee float64
}

// ExecuteOrder executes an order event
func (e *Exchange) ExecuteOrder(order OrderEvent, data DataHandler) (FillEvent, error) {
	// log.Printf("Exchange receives Order: \n%#v \n%#v\n", order, data)

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

	f.Commission = e.calculateCommission(float64(f.Qty), f.Price)
	f.ExchangeFee = e.calculateExchangeFee()
	f.Cost = e.calculateCost(f.Commission, f.ExchangeFee)

	return f, nil
}

// calculateComission() calculates the commission for a stock trade
//
// based on the conditions of IngDiba
// see https://www.ing-diba.de/wertpapiere/direkt-depot/konditionen
func (e *Exchange) calculateCommission(qty, price float64) float64 {
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
