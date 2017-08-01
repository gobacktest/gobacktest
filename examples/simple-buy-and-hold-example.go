package main

import (
	gbt "github.com/dirkolbrich/gobacktest"
	"github.com/dirkolbrich/gobacktest/internal"
)

func main() {
	// define symbols
	var symbols = []string{"BAS.DE", "DBK.DE"}

	// load new backtester
	bt := gbt.New()
	bt.SetSymbols(symbols)

	// set portfolio with inital cash
	portfolio := &internal.Portfolio{Cash: 10000}
	bt.SetPortfolio(portfolio)

	// create data provider and load data into the backtest
	data := &internal.BarEventFromCSVFileData{FileDir: "../data/"}
	data.Load(symbols)
	bt.SetData(data)

	// create strategy provider and load into the backtest
	strategy := &internal.SimpleStrategy{}
	bt.SetStrategy(strategy)

	// create execution provider and load into the backtest
	exchange := &internal.Exchange{}
	bt.SetExchange(exchange)

	bt.Run()
}
