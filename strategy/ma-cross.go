package strategy

import (
	// "fmt"

	gbt "github.com/dirkolbrich/gobacktest"
	"github.com/dirkolbrich/gobacktest/algo"
)

// MovingAverageCross is a test strategy, which interprets the SMA on a series of data events
// specified by ShortWindow (SW) and LongWindow (LW).
// If SW bigger tha LW and there is not already an invested BOT position, the strategy creates a buy signal.
// If SW falls below LW and there is an invested BOT position, the strategy creates an exit signal.
func MovingAverageCross(short, long int) *gbt.Strategy {
	// create a new strategy with an algo stack and load into the backtest
	strategy := gbt.NewStrategy("moving-average-cross")
	strategy.SetAlgo(
		algo.If(
			// condition
			algo.And(
				algo.BiggerThan(algo.SMA(short), algo.SMA(long)),
				algo.NotInvested(),
			),
			// action
			algo.CreateSignal("buy"), // create a buy signal
		),
		algo.If(
			// condition
			algo.And(
				algo.SmallerThan(algo.SMA(short), algo.SMA(long)),
				algo.IsInvested(),
			),
			// action
			algo.CreateSignal("exit"), // create a sell signal
		),
	)

	return strategy
}
