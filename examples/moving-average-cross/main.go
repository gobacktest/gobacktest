package main

import (
	"fmt"
	"time"

	"github.com/dirkolbrich/gobacktest/pkg/backtest"
	"github.com/dirkolbrich/gobacktest/pkg/data"
	"github.com/dirkolbrich/gobacktest/pkg/strategy"
)

func main() {
	// define symbols
	var symbols = []string{"SDF.DE"}

	// initiate new backtester and load symbols
	test := backtest.New()
	test.SetSymbols(symbols)

	startDataLoad := time.Now()
	// create data provider and load data into the backtest
	data := &data.BarEventFromCSVFile{FileDir: "../testdata/bar/"}
	data.Load(symbols)
	test.SetData(data)

	stopDataLoad := time.Now()
	fmt.Printf("Loading data took %v ms \n", stopDataLoad.Sub(startDataLoad).Seconds()*1000)

	// set portfolio with initial cash and default size and risk manager
	portfolio := &backtest.Portfolio{}
	portfolio.SetInitialCash(10000)

	sizeManager := &backtest.Size{DefaultSize: 200, DefaultValue: 2500}
	portfolio.SetSizeManager(sizeManager)

	riskManager := &backtest.Risk{}
	portfolio.SetRiskManager(riskManager)

	test.SetPortfolio(portfolio)

	// create strategy provider and load into the backtest
	strategy := &strategy.MovingAverageCross{ShortWindow: 50, LongWindow: 200}
	test.SetStrategy(strategy)

	// create execution provider and load into the backtest
	exchange := &backtest.Exchange{Symbol: "XTRA", ExchangeFee: 1.00}
	test.SetExchange(exchange)

	// choose a statistic and load into the backtest
	statistic := &backtest.Statistic{}
	test.SetStatistic(statistic)

	startRun := time.Now()
	// run the backtest
	test.Run()

	stopRun := time.Now()
	fmt.Printf("Running backtest took %v ms \n", stopRun.Sub(startRun).Seconds()*1000)

	// print the result of the test
	test.Stats().PrintResult()
}
