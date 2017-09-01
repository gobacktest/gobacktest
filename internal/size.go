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
func (s *Size) SizeOrder(order OrderEvent, current DataEvent, positions map[string]position) (OrderEvent, error) {
	// no default set, no sizing possible, order rejected
	if (s.DefaultSize == 0) || (s.DefaultValue == 0) {
		return order, errors.New("cant size order: no defaulsize or defaultValue set,")
	}

	// simple implementation, just gives the received order back
	// with default size
	order.SetQty(s.DefaultSize)

	return order, nil
}
