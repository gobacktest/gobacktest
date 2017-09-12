package internal

import (
	"errors"
	"math"
)

// SizeHandler is the basic interface for setting the size of an order
type SizeHandler interface {
	SizeOrder(OrderEvent, DataEvent, PortfolioHandler) (OrderEvent, error)
}

// Size is a basic size handler implementation
type Size struct {
	DefaultSize  int64
	DefaultValue float64
}

// SizeOrder adjusts the size of an order
func (s *Size) SizeOrder(order OrderEvent, data DataEvent, pf PortfolioHandler) (OrderEvent, error) {
	// no default set, no sizing possible, order rejected
	if (s.DefaultSize == 0) || (s.DefaultValue == 0) {
		return order, errors.New("cannot size order: no defaulSize or defaultValue set,")
	}

	// decide on order direction
	switch order.Direction() {
	case "long":
		order.SetDirection("buy")
		order.SetQty(s.setDefaultSize(data.LatestPrice()))
	case "short":
		order.SetDirection("sell")
		order.SetQty(s.setDefaultSize(data.LatestPrice()))
	case "exit": // all shares should be sold or bought, depending on position
		// poll postions
		if _, ok := pf.IsInvested(order.Symbol()); !ok {
			return order, errors.New("cannot exit order: no position to symbol in portfolio,")
		}
		if pos, ok := pf.IsLong(order.Symbol()); ok {
			order.SetDirection("sell")
			order.SetQty(pos.qty)
		}
		if pos, ok := pf.IsShort(order.Symbol()); ok {
			order.SetDirection("buy")
			order.SetQty(pos.qty * -1)
		}
	}

	return order, nil
}

func (s *Size) setDefaultSize(price float64) (qty int64) {
	if (float64(s.DefaultSize) * price) > s.DefaultValue {
		correctedQty := int64(math.Floor(s.DefaultValue / price))
		return correctedQty
	}
	return s.DefaultSize
}
