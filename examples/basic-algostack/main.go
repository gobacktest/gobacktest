package main

import (
	"fmt"
	"github.com/dirkolbrich/gobacktest/pkg/algo"
	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
	"github.com/dirkolbrich/gobacktest/pkg/data"
)

func main() {
	// initiate new backtester
	test := bt.New()

	// define and load symbols
	var symbols = []string{"TEST.DE"}
	test.SetSymbols(symbols)

	// create data provider and load data into the backtest
	data := &data.BarEventFromCSVFile{FileDir: "../testdata/test/"}
	data.Load(symbols)
	test.SetData(data)

	// create new strategy with algo stack and load into the backtest
	strategy := bt.NewStrategy("basic")
	strategy.SetAlgo(
		algo.TrueAlgo{},
	)

	// create assets and append to strategy
	strategy.SetChildren(bt.NewAsset("TEST.DE"))
	fmt.Printf("%#v\n", strategy)

	// load strategy into the backtest
	test.SetStrategy(strategy)

	// run the backtest
	err := test.Run()
	if err != nil {
		fmt.Println(err)
	}

	// print the result of the test
	test.Stats().PrintResult()
}
