package main

import (
	"fmt"

	gbt "github.com/dirkolbrich/gobacktest"
	"github.com/dirkolbrich/gobacktest/algo"
	"github.com/dirkolbrich/gobacktest/data"
)

func main() {
	// initiate new backtester
	test := gbt.New()

	// define and load symbols
	var symbols = []string{"TEST.DE"}
	test.SetSymbols(symbols)

	// create data provider and load data into the backtest
	data := &data.BarEventFromCSVFile{FileDir: "../testdata/test/"}
	data.Load(symbols)
	test.SetData(data)

	// create a new strategy with an algo stack and load into the backtest
	strategy := gbt.NewStrategy("basic")
	strategy.SetAlgo(
		algo.BoolAlgo(true),      // always return true, just a test
		algo.CreateSignal("buy"), // always create a buy signal on a data event
	)

	// create an asset and append to strategy
	strategy.SetChildren(gbt.NewAsset("TEST.DE"))

	// load the strategy into the backtest
	test.SetStrategy(strategy)

	// run the backtest
	err := test.Run()
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	// print the result of the test
	test.Stats().PrintResult()
}
