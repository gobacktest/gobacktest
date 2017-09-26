package backtest

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// CalculateSMA calculates the simple moving average of a given slice of data points.
func CalculateSMA(i int, data []DataEventHandler) (sma float64, err error) {
	// // fmt.Printf("Calculting SMA%d for %#v\n", i, e)
	// check for available data points
	if len(data) < i {
		return 0, fmt.Errorf("not enough data points to calculate SMA%d, only %d data points given, ", i, len(data))
	}

	sumClosePrices := decimal.New(0, 0)
	for j := 0; j < i; j++ {
		// take element from the end of the slice
		closePrice := decimal.NewFromFloat(data[len(data)-1-j].LatestPrice())
		sumClosePrices = sumClosePrices.Add(closePrice)
	}

	interval := decimal.New(int64(i), 0)

	sma, _ = sumClosePrices.Div(interval).Round(DP).Float64()
	return sma, nil
}
