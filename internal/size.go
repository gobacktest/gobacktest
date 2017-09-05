package internal

import (
	"errors"
)

// SizeHandler is the basic interface for setting the size of an order
type SizeHandler interface {
	SizeOrder(OrderEvent, DataEvent, map[string]position) (OrderEvent, error)
}

// Size is a basic size handler implementation
type Size struct {
	DefaultSize  int64
	DefaultValue float64
}

// SizeOrder adjusts the size of an order
func (s *Size) SizeOrder(order OrderEvent, latest DataEvent, positions map[string]position) (OrderEvent, error) {
	// no default set, no sizing possible, order rejected
	if (s.DefaultSize == 0) || (s.DefaultValue == 0) {
		return order, errors.New("cannot size order: no defaulSize or defaultValue set,")
	}

	// simple implementation, just gives the received order back
	// with default size
	order.SetQty(s.DefaultSize)

	// decide on order direction
	switch order.Direction() {
	case "long":
		order.SetDirection("buy")
	case "short":
		order.SetDirection("sell")
	case "exit": // all shares should be sold or bought, depending on position
		// poll postions
		position, ok := positions[order.Symbol()]
		if !ok {
			return order, errors.New("cannot exit order: no position to symbol in portfolio,")
		}
		if position.qty >= 0 {
			order.SetQty(position.qty)
			order.SetDirection("sell")
		}

		if position.qty < 0 {
			order.SetQty(position.qty * -1)
			order.SetDirection("buy")
		}

	}

	return order, nil
}
