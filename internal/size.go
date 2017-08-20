package internal

// SizeHandler is the basic interface for setting the size of an order
type SizeHandler interface {
	SizeOrder(OrderEvent, EventHandler, map[string]Position) (OrderEvent, bool)
}

// Size is a basic size handler implementation
type Size struct {
	DefaultSize  int64
	DefaultValue float64
}

// SizeOrder adjusts the size of an order
func (s *Size) SizeOrder(order OrderEvent, current EventHandler, positions map[string]Position) (OrderEvent, bool) {
	// no default set, no sizing possible, order rejected
	if (s.DefaultSize == nil) && (s.DefaultValue == nil) {
		return order, false
	}

	// simple implementation, just gives the received order back
	// with default size
	order.Qty = s.DefaultSize
	
	return order, true
}
