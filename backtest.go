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

// New creates a test backtest value for use.
func New(symbols []string) Test {
	test := Test{}
	test.symbols = symbols
	data := internal.Data{}
	d, err := test.loadData(data, symbols)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	test.data = d
	log.Printf("test.data: [%T] %+v \n", test.data, test.data)
	log.Printf("test: [%T] %+v \n", test, test)
	return test
}

// Run starts the test.
func (t Test) Run() {
	log.Println("Hello test.")
	log.Printf("t.symbols: [%T] %v \n", t.symbols, t.symbols)
	log.Printf("t.data: [%T] %v \n", t.data, t.data)
}

// load data into internal data struct
func (t Test) loadData(dh internal.DataHandler, symbols []string) (data internal.DataHandler, err error) {

	if len(symbols) == 0 {
		// bt.data.LoadAll()
		log.Println("No symbols given")
	} else {
		log.Printf("%d symbols given\n", len(t.symbols))
		err := dh.Load(t.symbols)
		if err != nil {
			// handle error
			log.Fatal(err)
		}
		log.Printf("dh: [%T] %+v \n", dh, dh)
	}

	return dh, nil
}

func (t Test) continueLoopCondition() bool {
	return true
}
