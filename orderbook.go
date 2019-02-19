package gobacktest

import (
	"fmt"
	"sort"
)

// OrderBookHandler defines the orderbook interface.
type OrderBookHandler interface {
	Add(Order) error
	Remove(Order) error
	Orders() ([]Order, bool)
	OrderBy(func(Order) bool) ([]Order, bool)
	OrdersBySymbol(string) ([]Order, bool)
	BidOrders() ([]Order, bool)
	AskOrders() ([]Order, bool)
	BidOrdersBySymbol(string) ([]Order, bool)
	AskOrdersBySymbol(string) ([]Order, bool)
	OpenOrders() ([]Order, bool)
	CanceledOrders() ([]Order, bool)
	History() ([]Order, bool)
}

// OrderBook represents a basic orderbook implementation.
type OrderBook struct {
	counter int
	orders  []Order
	history []Order
}

// Add an order to the orderbook.
func (ob *OrderBook) Add(order Order) error {
	// increment counter
	ob.counter++
	// assign an ID to the Order
	order.SetID(ob.counter)

	ob.orders = append(ob.orders, order)
	return nil
}

// Remove an order from the orderbook, append it to history.
func (ob *OrderBook) Remove(order Order) error {
	for i, o := range ob.orders {
		// order found
		if o.ID() == order.id {
			ob.history = append(ob.history, ob.orders[i])

			ob.orders = append(ob.orders[:i], ob.orders[i+1:]...)

			return nil
		}
	}

	// order not found in orderbook
	return fmt.Errorf("order with id %v not found", id)
}

// Orders returns all Orders from the order book.
func (ob OrderBook) Orders() ([]Order, bool) {
	if len(ob.orders) == 0 {
		return ob.orders, false
	}

	return ob.orders, true
}

// OrderBy returns the order defined by a select function from the orderbook.
func (ob OrderBook) OrderBy(fn func(order Order) bool) ([]Order, bool) {
	var orders = []Order{}

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

// OrdersBySymbol returns all orders of a specific symbol from the orderbook.
func (ob OrderBook) OrdersBySymbol(symbol string) ([]Order, bool) {
	var fn = func(order Order) bool {
		if order.Symbol() != symbol {
			return false
		}
		return true
	}

	orders, ok := ob.OrderBy(fn)
	return orders, ok
}

// BidOrders returns all bid orders from the orderbook.
func (ob OrderBook) BidOrders() ([]Order, bool) {
	var fn = func(order Order) bool {
		if order.Direction() != BOT {
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

// BidOrdersBySymbol returns all bid orders of a specific symbol from the orderbook.
func (ob OrderBook) BidOrdersBySymbol(symbol string) ([]Order, bool) {
	var fn = func(order Order) bool {
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

// AskOrders returns all ask orders from the orderbook.
func (ob OrderBook) AskOrders() ([]Order, bool) {
	var fn = func(order Order) bool {
		if order.Direction() != SLD {
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

// AskOrdersBySymbol returns all ask orders of a specific symbol from the orderbook.
func (ob OrderBook) AskOrdersBySymbol(symbol string) ([]Order, bool) {
	var fn = func(order Order) bool {
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

// OpenOrders returns all open orders from the orderbook.
func (ob OrderBook) OpenOrders() ([]Order, bool) {
	var fn = func(order Order) bool {
		if (order.Status() != OrderNew) || (order.Status() != OrderSubmitted) || (order.Status() != OrderPartiallyFilled) {
			return false
		}
		return true
	}

	orders, ok := ob.OrderBy(fn)
	return orders, ok
}

// CanceledOrders returns all canceled orders from the orderbook.
func (ob OrderBook) CanceledOrders() ([]Order, bool) {
	var fn = func(order Order) bool {
		if (order.Status() == OrderCanceled) || (order.Status() == OrderCancelPending) {
			return true
		}
		return false
	}

	orders, ok := ob.OrderBy(fn)
	return orders, ok
}

// History returns a slice of all historic orders.
func (ob OrderBook) History() (o []Order, ok bool) {
	if len(ob.history) == 0 {
		return ob.History, false
	}

	return ob.History, true
}

// Reset implements the Reseter interface and brings the OrderBook into a clean state.
func (ob OrderBook) Reset() error {
	ob.queue = nil
	ob.history = nil

	return nil
}
