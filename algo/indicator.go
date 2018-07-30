package algo

import (
	"fmt"

	gbt "github.com/dirkolbrich/gobacktest"
	"github.com/dirkolbrich/gobacktest/ta"
)

// smaAlgo is an algo which calculates the simple moving average.
type smaAlgo struct {
	gbt.Algo
	period int
	sma    float64
}

// SMA returns a sma algo ready to use.
func SMA(i int) gbt.AlgoHandler {
	return &smaAlgo{period: i}
}

// Run runs the algo.
func (a *smaAlgo) Run(s gbt.StrategyHandler) (bool, error) {
	data, _ := s.Data()
	event, _ := s.Event()
	symbol := event.Symbol()

	// prepare list of floats
	list := data.List(symbol)
	var values []float64

	if len(list) < a.period {
		return false, fmt.Errorf("invalid value length for indicator sma")
	}

	for i := 0; i < a.period; i++ {
		values = append(values, list[len(list)-i-1].Price())
	}

	// calculate SMA
	a.sma = ta.Mean(values)
	// save the calculated sma to the event metrics
	event.Add(fmt.Sprintf("SMA%d", a.period), a.sma)

	return true, nil
}

// Value returns the value of this Algo.
func (a *smaAlgo) Value() float64 {
	return a.sma
}
