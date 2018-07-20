package algo

import (
	"fmt"

	gbt "github.com/dirkolbrich/gobacktest"
	"github.com/dirkolbrich/gobacktest/ta"
)

// SMA is an Algo which calculates the small moving average.
type SMA struct {
	gbt.Algo
	period int
	sma    float64
}

// NewSMA returns a sma algo ready to use.
func NewSMA(i int) *SMA {
	return &SMA{period: i}
}

// Run runs the algo.
func (a *SMA) Run(s gbt.StrategyHandler) (bool, error) {
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
		values = append(values, list[len(list)-i-1].LatestPrice())
	}

	// calculate SMA
	a.sma = ta.Mean(values)

	return true, nil
}
