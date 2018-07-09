package backtest

import (
	"math"
	"time"
	// "github.com/shopspring/decimal"
)

// Position represents the holdings position
type Position struct {
	timestamp   time.Time
	symbol      string
	qty         int64   // current qty of the position, positive on BOT position, negativ on SLD position
	qtyBOT      int64   // how many BOT
	qtySLD      int64   // how many SLD
	avgPrice    float64 // average price without cost
	avgPriceNet float64 // average price including cost
	avgPriceBOT float64 // average price BOT, without cost
	avgPriceSLD float64 // average price SLD, without cost
	value       float64 // qty * price
	valueBOT    float64 // qty BOT * price
	valueSLD    float64 // qty SLD * price
	netValue    float64 // current value - cost
	netValueBOT float64 // current BOT value + cost
	netValueSLD float64 // current SLD value - cost
	marketPrice float64 // last known market price
	marketValue float64 // qty * price
	commission  float64
	exchangeFee float64
	cost        float64 // commission + fees
	costBasis   float64 // absolute qty * avgPriceNet

	realProfitLoss   float64
	unrealProfitLoss float64
	totalProfitLoss  float64
}

// Create a new position based on a fill event
func (p *Position) Create(fill FillEvent) {
	p.timestamp = fill.Time()
	p.symbol = fill.Symbol()

	p.update(fill)
}

// Update a position on a new fill event
func (p *Position) Update(fill FillEvent) {
	p.timestamp = fill.Time()

	p.update(fill)
}

// UpdateValue updates the current market value of a position
func (p *Position) UpdateValue(data DataEvent) {
	p.timestamp = data.Time()

	latest := data.LatestPrice()
	p.updateValue(latest)
}

// internal function to update a position on a new fill event
func (p *Position) update(fill FillEvent) {
	// convert fill to internally used decimal numbers
	fillQty := float64(fill.Qty())
	fillPrice := fill.Price()
	fillCommission := fill.Commission()
	fillExchangeFee := fill.ExchangeFee()
	fillCost := fill.Cost()
	fillNetValue := fill.NetValue()

	// convert position to internally used decimal numbers
	qty := float64(p.qty)
	qtyBot := float64(p.qtyBOT)
	qtySld := float64(p.qtySLD)
	avgPrice := p.avgPrice
	avgPriceNet := p.avgPriceNet
	avgPriceBot := p.avgPriceBOT
	avgPriceSld := p.avgPriceSLD
	value := p.value
	valueBot := p.valueBOT
	valueSld := p.valueSLD
	netValue := p.netValue
	netValueBot := p.netValueBOT
	netValueSld := p.netValueSLD
	commission := p.commission
	exchangeFee := p.exchangeFee
	cost := p.cost
	costBasis := p.costBasis
	realProfitLoss := p.realProfitLoss

	switch fill.Direction() {
	case "BOT":
		if p.qty >= 0 { // position is long, adding to position
			costBasis += fillNetValue
		} else { // position is short, closing partially out
			// costBasis + abs(fillQty) / qty * costBasis
			costBasis += math.Abs(fillQty) / qty * costBasis
			// realProfitLoss + fillQty * (avgPriceNet - fillPrice) - fillCost
			realProfitLoss += fillQty*(avgPriceNet-fillPrice) - fillCost
		}

		// update average price for bought stock without cost
		// ( (abs(qty) * avgPrice) + (fillQty * fillPrice) ) / (abs(qty) + fillQty)
		avgPrice = ((math.Abs(qty) * avgPrice) + (fillQty * fillPrice)) / (math.Abs(qty) + fillQty)
		// (abs(qty) * avgPriceNet + fillNetValue) / (abs(qty) * fillQty)
		avgPriceNet = (math.Abs(qty)*avgPriceNet + fillNetValue) / (math.Abs(qty) + fillQty)
		// ( (qty + avgPriceBot) + (fillQty * fillPrice) ) / fillQty
		avgPriceBot = ((qtyBot * avgPriceBot) + (fillQty * fillPrice)) / (qtyBot + fillQty)

		// update position qty
		qty += fillQty
		qtyBot += fillQty

		// update bought value
		valueBot = qtyBot * avgPriceBot
		netValueBot += fillNetValue

	case "SLD":
		if p.qty > 0 { // position is long, closing partially out
			costBasis -= math.Abs(fillQty) / qty * costBasis
			// realProfitLoss + fillQty * (fillPrice - avgPriceNet) - fillCost
			realProfitLoss += math.Abs(fillQty)*(fillPrice-avgPriceNet) - fillCost
		} else { // position is short, adding to position
			costBasis -= fillNetValue
		}

		// update average price for bought stock without cost
		// ( (abs(qty) * avgPrice) + (fillQty * fillPrice) ) / (abs(qty) + fillQty)
		avgPrice = (math.Abs(qty)*avgPrice + fillQty*fillPrice) / (math.Abs(qty) + fillQty)
		// (abs(qty) * avgPriceNet + fillNetValue) / (abs(qty) * fillQty)
		avgPriceNet = (math.Abs(qty)*avgPriceNet + fillNetValue) / (math.Abs(qty) + fillQty)
		// avgPriceSld + (fillQty * fillPrice) / fillQty
		avgPriceSld = (qtySld*avgPriceSld + fillQty*fillPrice) / (qtySld + fillQty)

		// update position qty
		qty -= fillQty
		qtySld += fillQty

		// update sold value
		valueSld = qtySld * avgPriceSld
		netValueSld += fillNetValue
	}

	commission += fillCommission
	exchangeFee += fillExchangeFee
	cost += fillCost
	value = valueSld - valueBot
	netValue = value - cost

	// convert from internal decimal to float
	p.qty = int64(qty)
	p.qtyBOT = int64(qtyBot)
	p.qtySLD = int64(qtySld)
	p.avgPrice = math.Round(avgPrice*math.Pow10(DP)) / math.Pow10(DP)
	p.avgPriceBOT = math.Round(avgPriceBot*math.Pow10(DP)) / math.Pow10(DP)
	p.avgPriceSLD = math.Round(avgPriceSld*math.Pow10(DP)) / math.Pow10(DP)
	p.avgPriceNet = math.Round(avgPriceNet*math.Pow10(DP)) / math.Pow10(DP)
	p.value = math.Round(value*math.Pow10(DP)) / math.Pow10(DP)
	p.valueBOT = math.Round(valueBot*math.Pow10(DP)) / math.Pow10(DP)
	p.valueSLD = math.Round(valueSld*math.Pow10(DP)) / math.Pow10(DP)
	p.netValue = math.Round(netValue*math.Pow10(DP)) / math.Pow10(DP)
	p.netValueBOT = math.Round(netValueBot*math.Pow10(DP)) / math.Pow10(DP)
	p.netValueSLD = math.Round(netValueSld*math.Pow10(DP)) / math.Pow10(DP)
	p.commission = commission
	p.exchangeFee = exchangeFee
	p.cost = cost
	p.costBasis = math.Round(costBasis*math.Pow10(DP)) / math.Pow10(DP)
	p.realProfitLoss = math.Round(realProfitLoss*math.Pow10(DP)) / math.Pow10(DP)

	p.updateValue(fill.Price())
}

// internal function to updates the current market value and profit/loss of a position
func (p *Position) updateValue(l float64) {
	// convert to internally used decimal numbers
	latest := l
	qty := float64(p.qty)
	costBasis := p.costBasis

	// update market value
	marketPrice := latest
	p.marketPrice = marketPrice
	// abs(qty) * current
	marketValue := math.Abs(qty) * latest
	p.marketValue = marketValue

	// qty * current - costBasis
	unrealProfitLoss := qty*latest - costBasis
	p.unrealProfitLoss = math.Round(unrealProfitLoss*math.Pow10(DP)) / math.Pow10(DP)

	realProfitLoss := p.realProfitLoss
	totalProfitLoss := realProfitLoss + unrealProfitLoss
	p.totalProfitLoss = math.Round(totalProfitLoss*math.Pow10(DP)) / math.Pow10(DP)
}
