package internal

import "log"

// ExecutionHandler is the basic interface for executing orders
type ExecutionHandler interface {
	ExecuteOrder(OrderEvent) (FillEvent, bool)
}

// Exchange is a basic execution handler implementation
type Exchange struct {
}

// ExecuteOrder executes an order event
func (e *Exchange) ExecuteOrder(o OrderEvent) (fill FillEvent, ok bool) {
	log.Printf("Reveiving order: Type: %s Symbol: %s Qty: %d", o.Direction, o.Symbol, o.Qty)

	return fill, true
}
