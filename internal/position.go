package internal

import "time"

type Position struct {
	timestamp          time.Time
	symbol             string
	qty                int64
	avgPrice           float64 // average price this position is aquired
	value              float64 // qty * price
	marketPrice        float64 // last known market price
	marketValue        float64 // qty * price
	commission         float64
	exchangeFee        float64
	cost               float64 // value - commision - fees
	netValue           float64 // current value - cost
	realisedProfitLoss float64
}
