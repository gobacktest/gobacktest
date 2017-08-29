// Package gobacktest provides a simple stock backtesting framework.
package gobacktest

import (
	"fmt"
	"log"

	"github.com/dirkolbrich/gobacktest/internal"
)

// Test is a basic back test struct
type Test struct {
	symbols      []string
	data         internal.DataHandler
	strategy     internal.StrategyHandler
	portfolio    internal.PortfolioHandler
	exchange     internal.ExecutionHandler
	eventQueue   []internal.Event
	eventHistory []internal.Event
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

	for data, ok := t.data.Next(); ok; data, ok = t.data.Next() {
		// add data event to event queue
		t.eventQueue = append(t.eventQueue, data)
		log.Printf("Portfolio: %+v\n", t.portfolio)

		// if event queue has an event, start event loop
		for event, ok := t.nextEvent(); ok; event, ok = t.nextEvent() {
			// add event to history
			t.eventHistory = append(t.eventHistory, event)

			// run event loop
			t.eventLoop(event)
		}
	}

	log.Printf("counted %d events\n", len(t.eventHistory))
}

// nextEvent gets the next event from the events queue
func (t *Test) nextEvent() (event internal.Event, ok bool) {
	// if event queue empty return false
	if len(t.eventQueue) == 0 {
		return event, false
	}

	// return first element from the event queue
	event = t.eventQueue[0]
	t.eventQueue = t.eventQueue[1:]

	return event, true
}

// eventLoop
func (t *Test) eventLoop(e internal.Event) {
	// symbol for this event
	fmt.Printf("%#v", e)
	symbol := e.Symbol()

	switch event := e.(type) {
	case internal.DataEvent:
		signal, ok := t.strategy.CalculateSignal(event)
		if !ok {
			break
		}
		t.eventQueue = append(t.eventQueue, signal)

		// portfolio should be updated here as well
		// to the last known price data

	case internal.SignalEvent:
		// get latest data event for this symbol
		current := t.data.Current(symbol)
		order, ok := t.portfolio.OnSignal(event, current)
		if !ok {
			break
		}
		t.eventQueue = append(t.eventQueue, order)

	case internal.OrderEvent:
		current := t.data.Current(symbol)
		fill, err := t.exchange.ExecuteOrder(event, current)
		if err != nil {
			break
		}
		t.eventQueue = append(t.eventQueue, fill)
	case internal.FillEvent:
		_, ok := t.portfolio.OnFill(event)
		if !ok {
			break
		}
		// log.Printf("Transaction recorded: %#v\n", transaction)
	}
}
