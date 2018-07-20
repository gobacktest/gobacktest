package main

import (
	"github.com/dirkolbrich/gobacktest"
	"github.com/dirkolbrich/gobacktest/data"
	"github.com/dirkolbrich/gobacktest/strategy"
)

func main() {
	// initiate new backtester
	test := gobacktest.New()

	// define and load symbols
	symbols := []string{"TEST.DE"}
	test.SetSymbols(symbols)

	// create data provider and load data into the backtest
	data := &data.BarEventFromCSVFile{FileDir: "../testdata/test/"}
	data.Load(symbols)
	test.SetData(data)

	// create strategy provider and load into the backtest
	strategy := &strategy.Basic{}
	test.SetStrategy(strategy)

	// run the backtest
	test.Run()

	// print the result of the test
	test.Stats().PrintResult()
}
