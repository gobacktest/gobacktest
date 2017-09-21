package backtest

// StrategyHandler is a basic strategy interface
type StrategyHandler interface {
	CalculateSignal(DataEventHandler, DataHandler, PortfolioHandler) (SignalEvent, error)
}
