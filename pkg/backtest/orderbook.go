package backtest

// OrderBook represents a order book.
type OrderBook struct {
	orders  []OrderEvent
	history []OrderEvent
}

// AddOrder adds an order to the order book.
func (ob *OrderBook) AddOrder(order OrderEvent) error {
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

// OrdersBy returns the order by a select function from the order book.
func (ob OrderBook) OrdersBy(fn func(order OrderEvent) bool) ([]OrderEvent, bool) {
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

// OrdersBySymbol returns the order of a specific symbol from the order book.
func (ob OrderBook) OrdersBySymbol(symbol string) ([]OrderEvent, bool) {
	var orders = []OrderEvent{}

	var symbolFunc = func(order OrderEvent) bool {
		if order.Symbol() != symbol {
			return false
		}
		return true
	}

	orders, ok := ob.OrdersBy(symbolFunc)
	return orders, ok
}
