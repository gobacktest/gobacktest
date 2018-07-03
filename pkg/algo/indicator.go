package algo

import (
	"fmt"

	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
	"github.com/shopspring/decimal"
)

// SMA is an Algo which calculates small moving average.
type SMA struct {
	bt.Algo
	period int
	sma    float64
}

// NewSMA returns a sma algo ready to use.
func NewSMA(i int) *SMA {
	return &SMA{period: i}
}

// Run runs the algo.
func (a *SMA) Run(s bt.StrategyHandler) (bool, error) {
	data, _ := s.Data()
	event, _ := s.Event()
	symbol := event.GetSymbol()

	// calculate SMA
	sma, err := calculateSMA(a.period, data.List(symbol))
	if err != nil {
		return false, err
	}
	a.sma = sma

	return true, nil
}

// calculateSMA calculates the simple moving average of a given slice of data points.
func calculateSMA(i int, data []bt.DataEventHandler) (float64, error) {
	// check for available data points
	if len(data) < i {
		return 0, fmt.Errorf("not enough data points to calculate SMA%d, only %d data points given", i, len(data))
	}

	sumClosePrices := decimal.New(0, 0)
	// start from the end of the slice
	for j := 0; j < i; j++ {
		// take element from the end of the slice
		closePrice := decimal.NewFromFloat(data[len(data)-1-j].LatestPrice())
		// fmt.Printf("i: %v j: %v close: %v ", i, j, closePrice)
		sumClosePrices = sumClosePrices.Add(closePrice)
		// fmt.Printf("sumClosePrices: %v\n", sumClosePrices)
	}
	interval := decimal.New(int64(i), 0)

	sma, _ := sumClosePrices.Div(interval).Round(bt.DP).Float64()
	return sma, nil
}
