package strategy

import (
	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

// Basic is a basic test strategy, which interprets every DataEvent as a signal to buy
type Basic struct {
}

// CalculateSignal handles the single Event
func (s *Basic) CalculateSignal(e bt.DataEventHandler, data bt.DataHandler, p bt.PortfolioHandler) (bt.SignalEvent, error) {
	// create Signal
	se := &bt.Signal{}

	// type switch for event type
	switch e := e.(type) {
	case *bt.Bar:
		// fill Signal
		se.Event = bt.Event{Timestamp: e.GetTime(), Symbol: e.GetSymbol()}
		se.Direction = "long"
	}

	return se, nil
}
