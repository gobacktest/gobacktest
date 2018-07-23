package gobacktest

// OrderBook represents a order book.
type OrderBook struct {
	orders  []OrderEvent
	history []OrderEvent
}

// Add an order to the order book.
func (ob *OrderBook) Add(order OrderEvent) error {
	ob.orders = append(ob.orders, order)
	return nil
}

// Orders returns all Orders from the order book
func (ob OrderBook) Orders() ([]OrderEvent, bool) {
	if len(ob.orders) == 0 {
		return ob.orders, false
	}

	return ob.orders, true
}

// OrderBy returns the order by a select function from the order book.
func (ob OrderBook) OrderBy(fn func(order OrderEvent) bool) ([]OrderEvent, bool) {
	var orders = []OrderEvent{}

	for _, order := range ob.orders {
		if fn(order) {
			orders = append(orders, order)
		}
	}

	if len(orders) == 0 {
		return orders, false
	}

	return orders, true
}

// OrderBySymbol returns the order of a specific symbol from the order book.
func (ob OrderBook) OrderBySymbol(symbol string) ([]OrderEvent, bool) {
	var orders = []OrderEvent{}

	var fn = func(order OrderEvent) bool {
		if order.Symbol() != symbol {
			return false
		}
		return true
	}

	orders, ok := ob.OrderBy(fn)
	return orders, ok
}

// OpenOrders returns all orders which are open from the order book.
func (ob OrderBook) OpenOrders() ([]OrderEvent, bool) {
	var orders = []OrderEvent{}

	var fn = func(order OrderEvent) bool {
		if (order.Status() != OrderNew) || (order.Status() != OrderSubmitted) || (order.Status() != OrderPartiallyFilled) {
			return false
		}
		return true
	}

	orders, ok := ob.OrderBy(fn)
	return orders, ok
}

// CanceledOrders returns all orders which are canceled from the order book.
func (ob OrderBook) CanceledOrders() ([]OrderEvent, bool) {
	var orders = []OrderEvent{}

	var fn = func(order OrderEvent) bool {
		if (order.Status() == OrderCanceled) || (order.Status() == OrderCancelPending) {
			return true
		}
		return false
	}

	orders, ok := ob.OrderBy(fn)
	return orders, ok
}
