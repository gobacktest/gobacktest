// Copyright 2017-present The gobacktest Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gobacktest

// DP sets the the precision of rounded floating numbers
// used after calculations to format the result.
const DP = 4 // DP

// Reseter defines a reseting interface.
type Reseter interface {
	Reset() error
}

// Backtest is the main struct which holds all elements of a test.
type Backtest struct {
	symbols   []string
	events    EventHandler
	data      DataHandler
	strategy  StrategyHandler
	portfolio PortfolioHandler
	orderbook OrderBookHandler
	exchange  ExecutionHandler
	statistic StatisticHandler
}

// New creates a backtest with sensible defaults ready for use.
func New() *Backtest {
	return &Backtest{
		events: EventStore{},
		portfolio: &Portfolio{
			initialCash: 100000,
			sizeManager: &Size{DefaultSize: 100, DefaultValue: 1000},
			riskManager: &Risk{},
		},
		orderbook: &OrderBook{},
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
	t.events.Reset()
	t.data.Reset()
	t.portfolio.Reset()
	t.statistic.Reset()
	return nil
}

// Statistics returns the statistic handler of the backtest.
func (t Backtest) Statistics() StatisticHandler {
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
	for event, ok := t.events.NextFromQueue(); true; event, ok = t.events.NextFromQueue() {
		// no event in the queue
		if !ok {
			// poll data stream
			data, ok := t.data.NextFromStream()
			// no more data, exit event loop
			if !ok {
				break
			}
			// found data event, add to event stream
			t.events.AppendToQueue(data)
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

// eventLoop directs the different events to their handler.
func (t Backtest) eventLoop(e Event) error {
	// type check for event type
	switch event := e.(type) {
	case DataEvent:
		// update portfolio to the last known price data
		t.portfolio.Update(event.Data)
		// update statistics
		t.statistic.Update(event.Data, t.portfolio)
		// check if any orders are filled before proceeding
		t.exchange.OnData(event.Data)

		// run strategy with this data event
		signals, err := t.strategy.OnData(event.Data)
		if err != nil {
			break
		}
		for _, signal := range signals {
			t.events.AppendToQueue(signal)
		}

	case SignalEvent:
		order, err := t.portfolio.OnSignal(event.Signal, t.data)
		if err != nil {
			break
		}
		t.events.AppendToQueue(order)

	case OrderEvent:
		fill, err := t.exchange.OnOrder(event.Order, t.data)
		if err != nil {
			break
		}
		t.events.AppendToQueue(fill)

	case FillEvent:
		transaction, err := t.portfolio.OnFill(event.Fill, t.data)
		if err != nil {
			break
		}
		t.statistic.TrackTransaction(transaction)
	}

	return nil
}
