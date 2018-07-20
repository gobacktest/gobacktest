package strategy

import (
	"fmt"

	gbt "github.com/dirkolbrich/gobacktest"
)

// BuyAndHold is a basic test strategy, which interprets the first DataEvent on a symbal
// as a signal to buy if the portfolio is not already invested.
type BuyAndHold struct {
}

// CalculateSignal handles the single Event
func (s *BuyAndHold) CalculateSignal(e gbt.DataEvent, data gbt.DataHandler, p gbt.PortfolioHandler) (gbt.SignalEvent, error) {
	// create Signal
	se := &gbt.Signal{}

	// type switch for event type
	switch e := e.(type) {
	case *gbt.Bar:
		// check if already invested
		if _, ok := p.IsInvested(e.Symbol()); ok {
			return se, fmt.Errorf("already invested in %v, no signal created,", e.Symbol())
		}
		// fill Signal
		event := gbt.Event{}
		event.SetTime(e.Time())
		event.SetSymbol(e.Symbol())

		se.Event = event
		se.SetDirection(gbt.BOT)
	}

	return se, nil
}
