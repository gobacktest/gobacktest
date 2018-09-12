package gobacktest

// DP sets the the precision of rounded floating numbers
// used after calculations to format.
const DP = 4 // DP

// Reseter provides a reseting interface.
type Reseter interface {
	Reset() error
}

// Backtest is the main struct which holds all elements of a test.
type Backtest struct {
	symbols    []string
	data       DataHandler
	strategy   StrategyHandler
	portfolio  PortfolioHandler
	exchange   ExecutionHandler
	statistic  StatisticHandler
	eventQueue []EventHandler
}

// New creates a default backtest with sensible defaults ready for use.
func New() *Backtest {
	return &Backtest{
		portfolio: &Portfolio{
			initialCash: 100000,
			sizeManager: &Size{DefaultSize: 100, DefaultValue: 1000},
			riskManager: &Risk{},
		},
		exchange: &Exchange{
			Symbol:      "TEST",
			Commission:  &FixedCommission{Commission: 0},
			ExchangeFee: &FixedExchangeFee{ExchangeFee: 0},
		},
		statistic: &Statistic{},
	}
}

// SetSymbols sets the symbols to include into the backtest.
func (t *Backtest) SetSymbols(symbols []string) {
	t.symbols = symbols
}

// SetData sets the data provider to be used within the backtest.
func (t *Backtest) SetData(data DataHandler) {
	t.data = data
}

// SetStrategy sets the strategy provider to be used within the backtest.
func (t *Backtest) SetStrategy(strategy StrategyHandler) {
	t.strategy = strategy
}

// SetPortfolio sets the portfolio provider to be used within the backtest.
func (t *Backtest) SetPortfolio(portfolio PortfolioHandler) {
	t.portfolio = portfolio
}

// SetExchange sets the execution provider to be used within the backtest.
func (t *Backtest) SetExchange(exchange ExecutionHandler) {
	t.exchange = exchange
}

// SetStatistic sets the statistic provider to be used within the backtest.
func (t *Backtest) SetStatistic(statistic StatisticHandler) {
	t.statistic = statistic
}

// Reset the backtest into a clean state with loaded data.
func (t *Backtest) Reset() error {
	t.eventQueue = nil
	t.data.Reset()
	t.portfolio.Reset()
	t.statistic.Reset()
	return nil
}

// Stats returns the statistic handler of the backtest.
func (t *Backtest) Stats() StatisticHandler {
	return t.statistic
}

// Run starts the backtest.
func (t *Backtest) Run() error {
	// setup before the backtest runs
	err := t.setup()
	if err != nil {
		return err
	}

	// poll event queue
	for event, ok := t.nextEvent(); true; event, ok = t.nextEvent() {
		// no event in the queue
		if !ok {
			// poll data stream
			data, ok := t.data.Next()
			// no more data, exit event loop
			if !ok {
				break
			}
			// found data event, add to event stream
			t.eventQueue = append(t.eventQueue, data)
			// start new event cycle
			continue
		}

		// processing event
		err := t.eventLoop(event)
		if err != nil {
			return err
		}
		// event in queue found, add to event history
		t.statistic.TrackEvent(event)
	}

	// teardown at the end of the backtest
	err = t.teardown()
	if err != nil {
		return err
	}

	return nil
}

// setup runs at the beginning of the backtest to perfom preparation operations.
func (t *Backtest) setup() error {
	// before first run, set portfolio cash
	t.portfolio.SetCash(t.portfolio.InitialCash())

	// make the data known to the strategy
	err := t.strategy.SetData(t.data)
	if err != nil {
		return err
	}

	// make the portfolio known to the strategy
	err = t.strategy.SetPortfolio(t.portfolio)
	if err != nil {
		return err
	}

	return nil
}

// teardown performs any cleaning operations at the end of the backtest.
func (t *Backtest) teardown() error {
	// no implementation yet
	return nil
}

// nextEvent gets the next event from the events queue.
func (t *Backtest) nextEvent() (e EventHandler, ok bool) {
	// if event queue empty return false
	if len(t.eventQueue) == 0 {
		return e, false
	}

	// return first element from the event queue
	e = t.eventQueue[0]
	t.eventQueue = t.eventQueue[1:]

	return e, true
}

// eventLoop directs the different events to their handler.
func (t *Backtest) eventLoop(e EventHandler) error {
	// type check for event type
	switch event := e.(type) {
	case DataEvent:
		// update portfolio to the last known price data
		t.portfolio.Update(event)
		// update statistics
		t.statistic.Update(event, t.portfolio)
		// check if any orders are filled before proceeding
		t.exchange.OnData(event)

		// run strategy with this data event
		signals, err := t.strategy.OnData(event)
		if err != nil {
			break
		}
		for _, signal := range signals {
			t.eventQueue = append(t.eventQueue, signal)
		}

	case *Signal:
		order, err := t.portfolio.OnSignal(event, t.data)
		if err != nil {
			break
		}
		t.eventQueue = append(t.eventQueue, order)

	case *Order:
		fill, err := t.exchange.OnOrder(event, t.data)
		if err != nil {
			break
		}
		t.eventQueue = append(t.eventQueue, fill)

	case *Fill:
		transaction, err := t.portfolio.OnFill(event, t.data)
		if err != nil {
			break
		}
		t.statistic.TrackTransaction(transaction)
	}

	return nil
}
