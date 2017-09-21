package backtest

// RiskHandler is the basic interface for accessing risks of a portfolio
type RiskHandler interface {
	EvaluateOrder(OrderEvent, DataEventHandler, map[string]position) (OrderEvent, error)
}

// Risk is a basic risk handler implementation
type Risk struct {
}

// EvaluateOrder handles the risk of an order, refines or cancel it
func (r *Risk) EvaluateOrder(order OrderEvent, data DataEventHandler, positions map[string]position) (OrderEvent, error) {
	// simple implementation, just gives the received order back
	// no risk management
	return order, nil
}
