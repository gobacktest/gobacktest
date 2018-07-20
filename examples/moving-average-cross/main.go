package main

import (
	"fmt"
	"time"

	"github.com/dirkolbrich/gobacktest"
	"github.com/dirkolbrich/gobacktest/data"
	"github.com/dirkolbrich/gobacktest/strategy"
)

func main() {
	// initiate new backtester
	test := gobacktest.New()

	// define and load symbols
	symbols := []string{"SDF.DE"}
	test.SetSymbols(symbols)

	// create data provider and load data into the backtest
	startDataLoad := time.Now()

	data := &data.BarEventFromCSVFile{FileDir: "../testdata/bar/"}
	data.Load(symbols)
	test.SetData(data)

	stopDataLoad := time.Now()
	fmt.Printf("Loading data took %v ms \n", stopDataLoad.Sub(startDataLoad).Seconds()*1000)

	// set default portfolio and redefine size manager
	portfolio := gobacktest.NewPortfolio()

	sizeManager := &gobacktest.Size{DefaultSize: 200, DefaultValue: 2500}
	portfolio.SetSizeManager(sizeManager)

	test.SetPortfolio(portfolio)

	// create strategy provider and load into the backtest
	strategy := &strategy.MovingAverageCross{ShortWindow: 50, LongWindow: 200}
	test.SetStrategy(strategy)

	startRun := time.Now()
	// run the backtest
	test.Run()

	stopRun := time.Now()
	fmt.Printf("Running backtest took %v ms \n", stopRun.Sub(startRun).Seconds()*1000)

	// print the result of the test
	test.Stats().PrintResult()
}
