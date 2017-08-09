package main

import (
    "github.com/dirkolbrich/gobacktest"
    "github.com/dirkolbrich/gobacktest/internal"
)

func main() {
    // define symbols
    var symbols = []string{"SDF.DE"}

    // initiate new backtester and load symbols
    test := gobacktest.New()
    test.SetSymbols(symbols)

    // set portfolio with inital cash and risk manager
    portfolio := &internal.Portfolio{Cash: 10000}
    riskManager := &internal.Risk{}
    portfolio.SetRiskManager(riskManager)
    test.SetPortfolio(portfolio)

    // create data provider and load data into the backtest
    data := &internal.BarEventFromCSVFileData{FileDir: "../data/test/"}
    data.Load(symbols)
    test.SetData(data)

    // create strategy provider and load into the backtest
    strategy := &internal.SimpleStrategy{}
    test.SetStrategy(strategy)

    // create execution provider and load into the backtest
    exchange := &internal.Exchange{}
    test.SetExchange(exchange)

    // run the backtest
    test.Run()

}
