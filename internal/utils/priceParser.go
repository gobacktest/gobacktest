package utils

import (
	"math"
)

const priceMultiplier = 10000000

/*
PriceParser is designed to abstract away the underlying number used as a price.
Due to efficiency and floating point precision limitations, this backtester uses
an integer to represent all prices. This means that $0.10 is, internally, 100,000,000.
Because such large numbers are rather unwieldy for humans,
the PriceParser will take a "normal" 2 decimal places numbers as input, and show
normal" 2 decimal places numbers as output when requested to `display()`
For consistency's sake, PriceParser should be used for ALL prices that enter
the backtester system. Numbers should also always be parsed correctly to view.
*/
type PriceParser struct{}

// Parse a float to int
func (p *PriceParser) Parse(f float64) int64 {
	return int64(f) * priceMultiplier
}

// Display an int as float with 2 decimal places
func (p *PriceParser) Display(i int64) float64 {
	return round(float64(i/priceMultiplier), 2)
}

// DisplayDecimal an int as float with custom decimal places
func (p *PriceParser) DisplayDecimal(i int64, decimal int) float64 {
	return round(float64(i/priceMultiplier), decimal)
}

func round(f float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	return math.Floor((f*shift)+.5) / shift
}
