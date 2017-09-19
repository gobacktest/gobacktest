package backtest

import (
	"fmt"
)

// StrategyHandler is a basic strategy interface
type StrategyHandler interface {
	CalculateSignal(DataEventHandler, DataHandler, PortfolioHandler) (SignalEvent, error)
}

// SimpleStrategy is a basic test strategy, which interprets every DataEvent as a signal to buy
type SimpleStrategy struct {
}

// CalculateSignal handles the single Event
func (s *SimpleStrategy) CalculateSignal(e DataEventHandler, data DataHandler, p PortfolioHandler) (SignalEvent, error) {
	// create Signal
	se := &Signal{}
	
	// type switch for event type
	switch e := e.(type) {
	case Bar:
		// fill Signal
		se.Event = Event{Timestamp: e.GetTime(), Symbol: e.GetSymbol()}
		se.Direction = "long"
	}

	return se, nil
}

// BuyAndHoldStrategy is a basic test strategy, which interprets the first DataEvent on a symbal
// as a signal to buy if the portfolio is not already invested.
type BuyAndHoldStrategy struct {
}

// CalculateSignal handles the single Event
func (s *BuyAndHoldStrategy) CalculateSignal(e DataEventHandler, data DataHandler, p PortfolioHandler) (SignalEvent, error) {
	// create Signal
	se := &Signal{}

	// type switch for event type
	switch e := e.(type) {
	case Bar:
		// check if already invested
		if _, ok := p.IsInvested(e.GetSymbol()); ok {
			return se, fmt.Errorf("already invested in %v, no signal created,", e.GetSymbol())
		}
		// fill Signal
		se.Event = Event{Timestamp: e.GetTime(), Symbol: e.GetSymbol()}
		se.Direction = "long"
	}

	return se, nil
}
