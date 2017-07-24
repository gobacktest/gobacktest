# gobacktest - fundamental stock analysis backtesting 

My attempt to create a event-driven backtesting framework to test stock trading strategies based on fundamental analysis. Preverably this package will be the core of a backend service exposed via a REST API.

## basic components

1. BackTester - general test case, bundles the follwing elements into a single test
2. DataHandler - interface to a set of data, e.g historical quotes, fundamental data etc.
3. StrategyHandler - generates a buy/sell signal based on the data
4. PortfolioHandler - generates orders and manages profit & loss
    + (RiskHandler) - manages the risk allocation of a portfolio
5. ExecutionHandler - sends orders to the broker and receives the “fills” or signals that the stock has been bought or sold

## infrastructure

An overviev of the infrastructure of a complete backtesting and trading environment

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

## resources

### articles

- Initial idea via a blog post [Python For Finance: Algorithmic Trading](https://www.datacamp.com/community/tutorials/finance-python-trading#backtesting) by Karlijn Willems [@willems_karlijn](https://twitter.com/willems_karlijn)
- Very good explanation of the internals of a back testing system by Michael Halls-Moore [@mhallsmoore](https://twitter.com/mhallsmoore) in the blog post series [Event-Driven-Backtesting-with-Python](https://www.quantstart.com/articles/Event-Driven-Backtesting-with-Python-Part-I).

### other backtesting frameworks

- [QuantConnect](https://www.quantconnect.com)
- [Quantopian](https://www.quantopian.com)
- [QuantRocket](https://www.quantrocket.com) - in development, available Q2/2018
- [Quandl](https://www.quandl.com) - financial data
- [QSTrader](https://www.quantstart.com/qstrader) - open-source backtesting framework from [QuantStart](https://www.quantstart.com)

### general information on quantitative finance

 - [Quantocracy](http://quantocracy.com) - forum for quant news
 - [QuantStart](https://www.quantstart.com) - articels and tutorials about quant finance

