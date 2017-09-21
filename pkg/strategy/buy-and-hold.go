package strategy

import (
	"fmt"

	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

// BuyAndHold is a basic test strategy, which interprets the first DataEvent on a symbal
// as a signal to buy if the portfolio is not already invested.
type BuyAndHold struct {
}

// CalculateSignal handles the single Event
func (s *BuyAndHold) CalculateSignal(e bt.DataEventHandler, data bt.DataHandler, p bt.PortfolioHandler) (bt.SignalEvent, error) {
	// create Signal
	se := &bt.Signal{}

	// type switch for event type
	switch e := e.(type) {
	case bt.Bar:
		// check if already invested
		if _, ok := p.IsInvested(e.GetSymbol()); ok {
			return se, fmt.Errorf("already invested in %v, no signal created,", e.GetSymbol())
		}
		// fill Signal
		se.Event = bt.Event{Timestamp: e.GetTime(), Symbol: e.GetSymbol()}
		se.Direction = "long"
	}

	return se, nil
}
