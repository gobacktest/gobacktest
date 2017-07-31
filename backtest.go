// Package gobacktest provides a simple stock backtesting framework.
package gobacktest

import (
	"log"

	"github.com/dirkolbrich/gobacktest/internal"
)

// Test is a basic back test struct
type Test struct {
	symbols []string
	data    internal.DataHandler
	// strategy     internal.StrategyHandler
	// eventStream []internal.EventHandler
	// portfolio    internal.PortfolioHandler
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

// Run starts the test.
func (t *Test) Run() {
	log.Println("Hello test.")

	// view the first entries of the data stream
	for i, v := range t.data.Stream() {
		// type switch for event type in stream
		switch val := v.(type) {
		case internal.BarEvent:
			if i < 10 {
				log.Printf("%d: dataEvent: %s %s\n", i, val.Date.Format("2006-01-02"), val.Symbol)
			}
		}
	}
}

func (t Test) continueLoopCondition() bool {
	return true
}
