package internal

// RiskHandler is the basic interface for accessing risks of a portfolio
type RiskHandler interface {
	EvaluateOrder(OrderEvent, map[string]Position) (OrderEvent, bool)
}

// Risk is a basic risk handler implementation
type Risk struct {
}

// EvaluateOrder handles the risk of an order, refines or cancel it
func (r *Risk) EvaluateOrder(order OrderEvent, positions map[string]Position) (OrderEvent, bool) {
	// simple implementation, just gives the received order back
	// no risk management
	return order, true
}
