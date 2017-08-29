package main

import (
	"github.com/dirkolbrich/gobacktest"
	"github.com/dirkolbrich/gobacktest/internal"
)

func main() {
	// define symbols
	var symbols = []string{"TEST.DE"}

	// initiate new backtester and load symbols
	test := gobacktest.New()
	test.SetSymbols(symbols)

	// create data provider and load data into the backtest
	data := &internal.BarEventFromCSVFileData{FileDir: "../data/test/"}
	data.Load(symbols)
	test.SetData(data)

	// set portfolio with inital cash and default size and risk manager
	portfolio := &internal.Portfolio{Cash: 10000}

	sizeManager := &internal.Size{DefaultSize: 100, DefaultValue: 1000}
	portfolio.SetSizeManager(sizeManager)

	riskManager := &internal.Risk{}
	portfolio.SetRiskManager(riskManager)

	test.SetPortfolio(portfolio)

	// create strategy provider and load into the backtest
	strategy := &internal.SimpleStrategy{}
	test.SetStrategy(strategy)

	// create execution provider and load into the backtest
	exchange := &internal.Exchange{Symbol: "TEST", ExchangeFee: 1.00}
	test.SetExchange(exchange)

	// run the backtest
	test.Run()

}
