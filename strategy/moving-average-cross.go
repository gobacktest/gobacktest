package strategy

import (
	"fmt"

	gbt "github.com/dirkolbrich/gobacktest"
	"github.com/dirkolbrich/gobacktest/ta"
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
func (s *MovingAverageCross) CalculateSignal(e gbt.DataEvent, data gbt.DataHandler, p gbt.PortfolioHandler) (gbt.SignalEvent, error) {
	// create empty Signal
	se := &gbt.Signal{}

	// type switch for event type
	switch e := e.(type) {
	case *gbt.Bar:
		// get last price values
		var values []float64
		list := data.List(e.Symbol())
		for i, v := range list {
			values[i] = v.LatestPrice()
		}

		// calculate and set SMA for short window
		smaShort, err := ta.SMA(values, s.ShortWindow)
		if err != nil {
			return se, err
		}
		e.Metric.Add(fmt.Sprintf("SMA%d", s.ShortWindow), smaShort[0])

		// calculate and set SMA for long window
		smaLong, err := ta.SMA(values, s.LongWindow)
		if err != nil {
			return se, err
		}
		e.Metric.Add(fmt.Sprintf("SMA%d", s.LongWindow), smaLong[0])

		// check if already invested
		_, invested := p.IsInvested(e.Symbol())

		if (smaShort[0] > smaLong[0]) && invested {
			return se, fmt.Errorf("buy signal but already invested in %v, no signal created,", e.Symbol())
		}

		if (smaShort[0] > smaLong[0]) && !invested {
			// buy signal, populate the signal event
			event := gbt.Event{}
			event.SetTime(e.Time())
			event.SetSymbol(e.Symbol())

			se.Event = event
			se.SetDirection(gbt.BOT)
		}

		if (smaShort[0] <= smaLong[0]) && !invested {
			return se, fmt.Errorf("sell signal but not invested in %v, no signal created,", e.Symbol())
		}

		if (smaShort[0] <= smaLong[0]) && invested {
			// sell signal, populate the signal event
			event := gbt.Event{}
			event.SetTime(e.Time())
			event.SetSymbol(e.Symbol())

			se.Event = event
			se.SetDirection(gbt.EXT)
		}

	}
	return se, nil
}
