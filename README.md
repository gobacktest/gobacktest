_**Heads up:** this is a framework in developement, not ready for any useful action. A lot of the functionality is still missing._

_You can read along and follow the development of this project. And if you like, give me some tips or discussion points for improvement._

---

# gobacktest - fundamental stock analysis backtesting 

My attempt to create a event-driven backtesting framework to test stock trading strategies based on fundamental analysis. Preferably this package will be the core of a backend service exposed via a REST API.

## usage

Example tests are in the `/examples` folder.

```golang
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

```

---

## basic components

These are the basic components of an event-driven framework. 

1. BackTester - general test case, bundles the follwing elements into a single test
2. DataHandler - interface to a set of data, e.g historical quotes, fundamental data etc.
3. StrategyHandler - generates a buy/sell signal based on the data
4. PortfolioHandler - generates orders and manages profit & loss
    + (RiskHandler) - manages the risk allocation of a portfolio
5. ExecutionHandler - sends orders to the broker and receives the “fills” or signals that the stock has been bought or sold
6. EventHandler - the different types of events, which travel through this system - data event, signal event, order event and fill event

---

## infrastructure

An overviev of the infrastructure of a complete backtesting and trading environment. Taken from the production roadmap of [QuantRocket](https://www.quantrocket.com/#product-roadmap).

- General
    + API gateway
    + configuration loader
    + logging service
    + cron service
- Data
    + database backup and download service
    + securities master services
    + historical market data service
    + fundamental data service
    + earnings data service
    + dividend data service
    + real-time market data service
    + exchange calendar service
- Strategy
    + performance analysis service - tearsheet
- Portfolio
    + account and portfolio service
    + risk management service
- Execution
    + trading platform gateway service
    + order management and trade ledger service
    + backtesting and trading engine

---

## resources

### articles

These links to articles are a good starting point to understand the intentions and basic functions of an event-driven backtesting framework.

- Initial idea via a blog post [Python For Finance: Algorithmic Trading](https://www.datacamp.com/community/tutorials/finance-python-trading#backtesting) by Karlijn Willems [@willems_karlijn](https://twitter.com/willems_karlijn).
- Very good explanation of the internals of a backtesting system by Michael Halls-Moore [@mhallsmoore](https://twitter.com/mhallsmoore) in the blog post series [Event-Driven-Backtesting-with-Python](https://www.quantstart.com/articles/Event-Driven-Backtesting-with-Python-Part-I).

### other backtesting frameworks

- [QuantConnect](https://www.quantconnect.com)
- [Quantopian](https://www.quantopian.com)
- [QuantRocket](https://www.quantrocket.com) - in development, available Q2/2018
- [Quandl](https://www.quandl.com) - financial data
- [QSTrader](https://www.quantstart.com/qstrader) - open-source backtesting framework from [QuantStart](https://www.quantstart.com)

### general information on quantitative finance

 - [Quantocracy](http://quantocracy.com) - forum for quant news
 - [QuantStart](https://www.quantstart.com) - articels and tutorials about quant finance

