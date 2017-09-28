package strategy

import (
	"fmt"

	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

// MovingAverageCross is a test strategy, which interprets the SMA on a series of data events
// specified by ShortWindow (SW) and LongWindow (LW).
// If SW bigger tha LW and there is not already an invested BOT position, the strategy creates a buy signal.
// If SW falls below LW and there is an invested BOT position, the strategy creates an exit signal.
type MovingAverageCross struct {
	ShortWindow int
	LongWindow  int
}

// CalculateSignal handles the single Event
func (s *MovingAverageCross) CalculateSignal(e bt.DataEventHandler, data bt.DataHandler, p bt.PortfolioHandler) (bt.SignalEvent, error) {
	// create empty Signal
	se := &bt.Signal{}

	// type switch for event type
	switch e := e.(type) {
	case *bt.Bar:
		// calculate and set SMA for short window
		smaShort, err := bt.CalculateSMA(s.ShortWindow, data.List(e.GetSymbol()))
		if err != nil {
			return se, err
		}
		e.Metrics[fmt.Sprintf("SMA%d", s.ShortWindow)] = smaShort

		// calculate and set SMA for long window
		smaLong, err := bt.CalculateSMA(s.LongWindow, data.List(e.GetSymbol()))
		if err != nil {
			return se, err
		}
		e.Metrics[fmt.Sprintf("SMA%d", s.LongWindow)] = smaLong

		// check if already invested
		_, invested := p.IsInvested(e.GetSymbol())

		if (smaShort > smaLong) && invested {
			return se, fmt.Errorf("buy signal but already invested in %v, no signal created,", e.GetSymbol())
		}

		if (smaShort > smaLong) && !invested {
			// buy signal, populate the signal event
			se.Event = bt.Event{Timestamp: e.GetTime(), Symbol: e.GetSymbol()}
			se.Direction = "long"
		}

		if (smaShort <= smaLong) && !invested {
			return se, fmt.Errorf("sell signal but not invested in %v, no signal created,", e.GetSymbol())
		}

		if (smaShort <= smaLong) && invested {
			// sell signal, populate the signal event
			se.Event = bt.Event{Timestamp: e.GetTime(), Symbol: e.GetSymbol()}
			se.Direction = "exit"
		}

	}
	return se, nil
}
