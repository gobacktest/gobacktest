# marty - fundamental stock analysis backtesting 

My attempt to create a backtesting framework to test stock trading strategies based on fundamental analysis. Preverably this package will be the backend exposed via a REST API.

## basic components

1. DataHandler - interface to a set of data, e.g historical quotes, fundamental data etc.
2. StrategyHandler - generates a buy/sell signal based on the data
3. PortfolioHandler - generates orders and manages profit & loss
4. ExecutionHandler - sends orders to the broker and receives the “fills” or signals that the stock has been bought or sold

## resources

### articles

Initial idea via a blog post [Python For Finance: Algorithmic Trading](https://www.datacamp.com/community/tutorials/finance-python-trading#backtesting) by Karlijn Willems [@willems_karlijn](https://twitter.com/willems_karlijn)

Very good explanation of the internals of a back testing system by Michael Halls-Moore [@mhallsmoore](https://twitter.com/mhallsmoore) in the blog post series [Event-Driven-Backtesting-with-Python](https://www.quantstart.com/articles/Event-Driven-Backtesting-with-Python-Part-I).

### other backtesting frameworks

[Quantopian](https://www.quantopian.com)

## trivia

Whats`s with the name marty?
Well, Marty McFly travels "Back to the Future" to see the consequences of his actions. 
