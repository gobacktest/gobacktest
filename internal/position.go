package internal

import (
	"time"

	"github.com/dirkolbrich/gobacktest/internal/utils"
)

// Position represents the holdings position
type position struct {
	timestamp        time.Time
	symbol           string
	qty              int64
	avgPrice         float64 // average price this position is aquired
	value            float64 // qty * price
	marketPrice      float64 // last known market price
	marketValue      float64 // qty * price
	commission       float64
	exchangeFee      float64
	cost             float64 // value - commision - fees
	netValue         float64 // current value - cost
	realProfitLoss   float64
	unrealProfitLoss float64
	totalProfitLoss  float64
}

// Create a new position based on a fill event
func (p *position) Create(fill FillEvent) {
	p.timestamp = fill.Timestamp()
	p.symbol = fill.Symbol()

	p.qty = fill.Qty()
	p.avgPrice = fill.Price()
	p.value = float64(fill.Qty()) * fill.Price()

	p.marketPrice = fill.Price()
	p.marketValue = p.value

	p.commission = fill.Commission()
	p.exchangeFee = fill.ExchangeFee()
	p.cost = fill.Cost()

	p.netValue = p.value - p.cost
}

// Update a position on a new fill event
func (p *position) Update(fill FillEvent) {
	p.timestamp = fill.Timestamp()

	p.avgPrice = (float64(p.qty)*p.avgPrice + float64(fill.Qty())*fill.Price()) / float64(p.qty+fill.Qty())
	p.qty += fill.Qty()
	p.value = float64(p.qty) * p.avgPrice

	p.marketPrice = fill.Price()
	p.marketValue = p.value

	p.commission = utils.Round(p.commission+fill.Commission(), 4)
	p.exchangeFee = utils.Round(p.exchangeFee+fill.ExchangeFee(), 4)
	p.cost = utils.Round(p.cost+fill.Cost(), 3)

	p.netValue = p.value - p.cost
}

// UpdateValue updates the current market value of a position
func (p *position) UpdateValue(current DataEvent) {
	p.timestamp = current.Timestamp()

	p.marketPrice = current.LatestPrice()
	p.marketValue = float64(p.qty) * p.marketPrice
}
