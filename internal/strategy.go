package internal

import (
	"fmt"
)

// StrategyHandler is a basic strategy interface
type StrategyHandler interface {
	CalculateSignal(DataEvent, DataHandler, PortfolioHandler) (SignalEvent, error)
}

// SimpleStrategy is a basic test strategy, which interprets every DataEvent as a signal to buy
type SimpleStrategy struct {
}

// CalculateSignal handles the single Event
func (s *SimpleStrategy) CalculateSignal(e DataEvent, data DataHandler, p PortfolioHandler) (se SignalEvent, err error) {
	// type switch for event type
	switch e := e.(type) {
	case bar:
		// create Signal
		se = &signal{
			event:     event{timestamp: e.Timestamp(), symbol: e.Symbol()},
			direction: "long",
		}
	}

	return se, nil
}

// BuyAndHoldStrategy is a basic test strategy, which interprets the first DataEvent on a symbal
// as a signal to buy if the portfolio is not already invested.
type BuyAndHoldStrategy struct {
}

// CalculateSignal handles the single Event
func (s *BuyAndHoldStrategy) CalculateSignal(e DataEvent, data DataHandler, p PortfolioHandler) (se SignalEvent, err error) {
	// type switch for event type
	switch e := e.(type) {
	case bar:
		// check if already invested
		if _, ok := p.IsInvested(e.Symbol()); ok {
			return se, fmt.Errorf("already invested in %v, no signal created,", e.Symbol())
		}
		// create Signal
		se = &signal{
			event:     event{timestamp: e.Timestamp(), symbol: e.Symbol()},
			direction: "long",
		}
	}

	return se, nil
}
