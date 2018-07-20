package strategy

import (
	gbt "github.com/dirkolbrich/gobacktest"
)

// Basic is a basic test strategy, which interprets every DataEvent as a signal to buy
type Basic struct {
}

// CalculateSignal handles the single Event
func (s *Basic) CalculateSignal(e gbt.DataEvent, data gbt.DataHandler, p gbt.PortfolioHandler) (gbt.SignalEvent, error) {
	// create Signal
	se := &gbt.Signal{}

	// type switch for event type
	switch e := e.(type) {
	case *gbt.Bar:
		// fill Signal
		event := gbt.Event{}
		event.SetTime(e.Time())
		event.SetSymbol(e.Symbol())

		se.Event = event
		se.SetDirection(gbt.BOT)
	}

	return se, nil
}
