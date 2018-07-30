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
	var symbols = []string{"SDF.DE"}
	test.SetSymbols(symbols)

	// create data provider and load data into the backtest
	data := &data.BarEventFromCSVFile{FileDir: "../testdata/bar/"}
	data.Load(symbols)
	test.SetData(data)

	// create a new strategy with an algo stack and load into the backtest
	strategy := gbt.NewStrategy("moving-average-cross")
	strategy.SetAlgo(
		algo.If(
			// condition
			algo.And(
				algo.BiggerThan(algo.SMA(50), algo.SMA(200)),
				algo.NotInvested(),
			),
			// action
			algo.CreateSignal("buy"), // create a buy signal
		),
		algo.If(
			// condition
			algo.And(
				algo.SmallerThan(algo.SMA(50), algo.SMA(200)),
				algo.IsInvested(),
			),
			// action
			algo.CreateSignal("exit"), // create a sell signal
		),
	)

	// create an asset and append to strategy
	strategy.SetChildren(gbt.NewAsset("SDF.DE"))

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
