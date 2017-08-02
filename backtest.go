// Package gobacktest provides a simple stock backtesting framework.
package gobacktest

import (
	"log"

	"github.com/dirkolbrich/gobacktest/internal"
)

// Test is a basic back test struct
type Test struct {
	symbols    []string
	data       internal.DataHandler
	strategy   internal.StrategyHandler
	portfolio  internal.PortfolioHandler
	exchange   internal.ExecutionHandler
	eventQueue []internal.EventHandler
}

// New creates a default test backtest value for use.
func New() *Test {
	return &Test{}
}

// SetSymbols sets the symbols to include into the test
func (t *Test) SetSymbols(symbols []string) {
	t.symbols = symbols
}

// SetData sets the data provider to to be used within the test
func (t *Test) SetData(data internal.DataHandler) {
	t.data = data
}

// SetStrategy sets the strategy provider to to be used within the test
func (t *Test) SetStrategy(strategy internal.StrategyHandler) {
	t.strategy = strategy
}

// SetPortfolio sets the portfolio provider to to be used within the test
func (t *Test) SetPortfolio(portfolio internal.PortfolioHandler) {
	t.portfolio = portfolio
}

// SetExchange sets the execution provider to to be used within the test
func (t *Test) SetExchange(exchange internal.ExecutionHandler) {
	t.exchange = exchange
}

// Run starts the test.
func (t *Test) Run() {
	log.Println("Running backtest:")

	events := 0
	for dataEvent, ok := t.data.Next(); ok; dataEvent, ok = t.data.Next() {
		// add data event to event queue
		t.eventQueue = append(t.eventQueue, dataEvent)

		// if event queue has an event, start event loop
		for event, ok := t.nextEvent(); ok; event, ok = t.nextEvent() {
			events++
			// type switch for event type
			switch ev := event.(type) {
			case internal.BarEvent:
				signal, ok := t.strategy.CalculateSignal(ev)
				if !ok {
					continue
				}
				t.eventQueue = append(t.eventQueue, signal)

				// portfolio should be updated here as well
				// to the last known price data

			case internal.SignalEvent:
				order, ok := t.portfolio.OnSignal(ev)
				if !ok {
					continue
				}
				t.eventQueue = append(t.eventQueue, order)
			case internal.OrderEvent:
				fill, ok := t.exchange.ExecuteOrder(ev)
				if !ok {
					continue
				}
				t.eventQueue = append(t.eventQueue, fill)
			case internal.FillEvent:
				transaction, ok := t.portfolio.OnFill(ev)
				if !ok {
					continue
				}
				log.Printf("Transaction recorded: %#v\n", transaction)
			}
		}
	}

	// dataStream should be empty now
	log.Printf("counted %d events\n", events)
	// log.Printf("dataStream is empty: %v\n", t.data.StreamIsEmpty())
	// log.Printf("eventQueue is empty: %v\n", t.queueIsEmpty())
}

func (t *Test) nextEvent() (event internal.EventHandler, ok bool) {
	// if event queue empty return false
	if len(t.eventQueue) == 0 {
		return event, false
	}

	// return first element in event queue
	event = t.eventQueue[0]
	t.eventQueue = t.eventQueue[1:]

	return event, true
}

func (t *Test) queueIsEmpty() bool {
	// if event queue is empty return false
	if len(t.eventQueue) == 0 {
		return true
	}
	return false
}
