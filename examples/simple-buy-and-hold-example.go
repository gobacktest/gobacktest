package main

import (
	gbt "github.com/dirkolbrich/gobacktest"
	"github.com/dirkolbrich/gobacktest/internal"
)

func main() {
	// define symbols
	var symbols = []string{"BAS.DE", "DBK.DE", "SZU.DE"}

	// load new backtester
	bt := gbt.New()
	bt.SetSymbols(symbols)

	// create data provider and load data into the backtest
	data := &internal.BarEventFromCSVFileData{FileDir: "../data/"}
	data.Load(symbols)
	bt.SetData(data)

	bt.Run()
}
