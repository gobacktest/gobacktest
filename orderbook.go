package gobacktest

import (
	"fmt"
	"sort"
)

// OrderBook represents an order book.
type OrderBook struct {
	counter int
	orders  []OrderEvent
	history []OrderEvent
}

// Add an order to the order book.
func (ob *OrderBook) Add(order OrderEvent) error {
	// increment counter
	ob.counter++
	// assign an ID to the Order
	order.SetID(ob.counter)

	ob.orders = append(ob.orders, order)
	return nil
}

// Remove an order from the order book, append it to history.
func (ob *OrderBook) Remove(id int) error {
	for i, order := range ob.orders {
		// order found
		if order.ID() == id {
			ob.history = append(ob.history, ob.orders[i])

			ob.orders = append(ob.orders[:i], ob.orders[i+1:]...)

			return nil
		}
	}

	// order not found in orderbook
	return fmt.Errorf("order with id %v not found", id)
}

// Orders returns all Orders from the order book.
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

// OrdersBySymbol returns the order of a specific symbol from the order book.
func (ob OrderBook) OrdersBySymbol(symbol string) ([]OrderEvent, bool) {
	var fn = func(order OrderEvent) bool {
		if order.Symbol() != symbol {
			return false
		}
		return true
	}

	orders, ok := ob.OrderBy(fn)
	return orders, ok
}

// OrdersBidBySymbol returns all bid orders of a specific symbol from the order book.
func (ob OrderBook) OrdersBidBySymbol(symbol string) ([]OrderEvent, bool) {
	var fn = func(order OrderEvent) bool {
		if (order.Symbol() != symbol) || (order.Direction() != BOT) {
			return false
		}
		return true
	}
	orders, ok := ob.OrderBy(fn)

	// sort bid orders ascending, lowest price first
	sort.Slice(orders, func(i, j int) bool {
		o1 := orders[i]
		o2 := orders[j]

		return o1.Limit() < o2.Limit()

	})

	return orders, ok
}

// OrdersAskBySymbol returns all bid orders of a specific symbol from the order book.
func (ob OrderBook) OrdersAskBySymbol(symbol string) ([]OrderEvent, bool) {
	var fn = func(order OrderEvent) bool {
		if (order.Symbol() != symbol) || (order.Direction() != SLD) {
			return false
		}
		return true
	}
	orders, ok := ob.OrderBy(fn)

	// sort bid orders descending, highest price first
	sort.Slice(orders, func(i, j int) bool {
		o1 := orders[i]
		o2 := orders[j]

		return o1.Limit() > o2.Limit()

	})

	return orders, ok
}

// OrdersOpen returns all open orders from the order book.
func (ob OrderBook) OrdersOpen() ([]OrderEvent, bool) {
	var fn = func(order OrderEvent) bool {
		if (order.Status() != OrderNew) || (order.Status() != OrderSubmitted) || (order.Status() != OrderPartiallyFilled) {
			return false
		}
		return true
	}

	orders, ok := ob.OrderBy(fn)
	return orders, ok
}

// OrdersCanceled returns all canceled orders from the order book.
func (ob OrderBook) OrdersCanceled() ([]OrderEvent, bool) {
	var fn = func(order OrderEvent) bool {
		if (order.Status() == OrderCanceled) || (order.Status() == OrderCancelPending) {
			return true
		}
		return false
	}

	orders, ok := ob.OrderBy(fn)
	return orders, ok
}
