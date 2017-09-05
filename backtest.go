// Package gobacktest provides a simple stock backtesting framework.
package gobacktest

import (
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
func (t *Test) Run() error {
	log.Println("Running backtest:")

	// poll event queue
	for event, ok := t.nextEvent(); true; event, ok = t.nextEvent() {
		// no event in queue
		if !ok {
			// poll data stream
			data, ok := t.data.Next()
			// no  data event, exit event loop
			if !ok {
				break
			}
			// found data, add to event stream
			t.eventQueue = append(t.eventQueue, data)
			// start new event polling cycle
			continue
		}
		// event in queue found, processing
		err := t.eventLoop(event)
		if err != nil {
			return err
		}

		// add event to history
		t.eventHistory = append(t.eventHistory, event)
	}

	log.Printf("counted %d events\n", len(t.eventHistory))

	return nil
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
func (t *Test) eventLoop(e internal.Event) error {
	// type check for event type
	switch event := e.(type) {
	case internal.DataEvent:
		signal, err := t.strategy.CalculateSignal(event, t.data)
		if err != nil {
			break
		}
		t.eventQueue = append(t.eventQueue, signal)

		// portfolio should be updated here as well
		// to the last known price data

	case internal.SignalEvent:
		order, err := t.portfolio.OnSignal(event, t.data)
		if err != nil {
			break
		}
		t.eventQueue = append(t.eventQueue, order)

	case internal.OrderEvent:
		fill, err := t.exchange.ExecuteOrder(event, t.data)
		if err != nil {
			break
		}
		t.eventQueue = append(t.eventQueue, fill)
	case internal.FillEvent:
		_, err := t.portfolio.OnFill(event, t.data)
		if err != nil {
			break
		}
		// log.Printf("Transaction recorded: %#v\n", transaction)
	}

	return nil
}
