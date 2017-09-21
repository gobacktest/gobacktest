package backtest

import (
	"time"

	"github.com/shopspring/decimal"
)

// Position represents the holdings position
type position struct {
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
func (p *position) Create(fill FillEvent) {
	p.timestamp = fill.GetTime()
	p.symbol = fill.GetSymbol()

	p.update(fill)
}

// Update a position on a new fill event
func (p *position) Update(fill FillEvent) {
	p.timestamp = fill.GetTime()

	p.update(fill)
}

// UpdateValue updates the current market value of a position
func (p *position) UpdateValue(data DataEventHandler) {
	p.timestamp = data.GetTime()

	latest := data.LatestPrice()
	p.updateValue(latest)
}

// internal function to update a position on a new fill event
func (p *position) update(fill FillEvent) {
	// convert fill to internally used decimal numbers
	fillQty := decimal.New(fill.GetQty(), 0)
	fillPrice := decimal.NewFromFloat(fill.GetPrice())
	fillCommission := decimal.NewFromFloat(fill.GetCommission())
	fillExchangeFee := decimal.NewFromFloat(fill.GetExchangeFee())
	fillCost := decimal.NewFromFloat(fill.GetCost())
	fillNetValue := decimal.NewFromFloat(fill.NetValue())

	// convert position to internally used decimal numbers
	qty := decimal.New(p.qty, 0)
	qtyBot := decimal.New(p.qtyBOT, 0)
	qtySld := decimal.New(p.qtySLD, 0)
	avgPrice := decimal.NewFromFloat(p.avgPrice)
	avgPriceNet := decimal.NewFromFloat(p.avgPriceNet)
	avgPriceBot := decimal.NewFromFloat(p.avgPriceBOT)
	avgPriceSld := decimal.NewFromFloat(p.avgPriceSLD)
	value := decimal.NewFromFloat(p.value)
	valueBot := decimal.NewFromFloat(p.valueBOT)
	valueSld := decimal.NewFromFloat(p.valueSLD)
	netValue := decimal.NewFromFloat(p.netValue)
	netValueBot := decimal.NewFromFloat(p.netValueBOT)
	netValueSld := decimal.NewFromFloat(p.netValueSLD)
	commission := decimal.NewFromFloat(p.commission)
	exchangeFee := decimal.NewFromFloat(p.exchangeFee)
	cost := decimal.NewFromFloat(p.cost)
	costBasis := decimal.NewFromFloat(p.costBasis)
	realProfitLoss := decimal.NewFromFloat(p.realProfitLoss)

	switch fill.GetDirection() {
	case "BOT":
		if p.qty >= 0 { // position is long, adding to position
			costBasis = costBasis.Add(fillNetValue)
		} else { // position is short, closing partially out
			costBasis = costBasis.Add(fillQty.Abs().Div(qty).Mul(costBasis))
			// realProfitLoss + fillQty * (avgPriceNet - fillPrice) - fillCost
			realProfitLoss = realProfitLoss.Add(fillQty.Mul(avgPriceNet.Sub(fillPrice))).Sub(fillCost)
		}

		// update average price for bought stock without cost

		// ( (abs(qty) * avgPrice) + (fillQty * fillPrice) ) / (abs(qty) + fillQty)
		avgPrice = qty.Abs().Mul(avgPrice).Add(fillQty.Mul(fillPrice)).Div(qty.Abs().Add(fillQty))
		// (abs(qty) * avgPriceNet + fillNetValue) / (abs(qty) * fillQty)
		avgPriceNet = qty.Abs().Mul(avgPriceNet).Add(fillNetValue).Div(qty.Abs().Add(fillQty))
		// ( (qty + avgPriceBot) + (fillQty * fillPrice) ) / fillQty
		avgPriceBot = qtyBot.Mul(avgPriceBot).Add(fillQty.Mul(fillPrice)).Div(qtyBot.Add(fillQty))

		// update position qty
		qty = qty.Add(fillQty)
		qtyBot = qtyBot.Add(fillQty)

		// update bought value
		valueBot = qtyBot.Mul(avgPriceBot)
		netValueBot = netValueBot.Add(fillNetValue)

	case "SLD":
		if p.qty > 0 { // position is long, closing partially out
			costBasis = costBasis.Sub(fillQty.Abs().Div(qty).Mul(costBasis))
			// realProfitLoss + fillQty * (fillPrice - avgPriceNet) - fillCost
			realProfitLoss = realProfitLoss.Add(fillQty.Abs().Mul(fillPrice.Sub(avgPriceNet))).Sub(fillCost)
		} else { // position is short, adding to position
			costBasis = costBasis.Sub(fillNetValue)
		}

		// update average price for bought stock without cost
		// ( (abs(qty) * avgPrice) + (fillQty * fillPrice) ) / (abs(qty) + fillQty)
		avgPrice = qty.Abs().Mul(avgPrice).Add(fillQty.Mul(fillPrice)).Div(qty.Abs().Add(fillQty))
		// (abs(qty) * avgPriceNet + fillNetValue) / (abs(qty) * fillQty)
		avgPriceNet = qty.Abs().Mul(avgPriceNet).Add(fillNetValue).Div(qty.Abs().Add(fillQty))
		// avgPriceSld + (fillQty * fillPrice) / fillQty
		avgPriceSld = qtySld.Mul(avgPriceSld).Add(fillQty.Mul(fillPrice)).Div(qtySld.Add(fillQty))

		// update position qty
		qty = qty.Sub(fillQty)
		qtySld = qtySld.Add(fillQty)

		// update sold value
		valueSld = qtySld.Mul(avgPriceSld)
		netValueSld = netValueSld.Add(fillNetValue)
	}

	commission = commission.Add(fillCommission)
	exchangeFee = exchangeFee.Add(fillExchangeFee)
	cost = cost.Add(fillCost)

	value = valueSld.Sub(valueBot)
	netValue = value.Sub(cost)

	// convert from internal decimal to float
	p.qty = qty.IntPart()
	p.qtyBOT = qtyBot.IntPart()
	p.qtySLD = qtySld.IntPart()
	p.avgPrice, _ = avgPrice.Round(DP).Float64()
	p.avgPriceBOT, _ = avgPriceBot.Round(DP).Float64()
	p.avgPriceSLD, _ = avgPriceSld.Round(DP).Float64()
	p.avgPriceNet, _ = avgPriceNet.Round(DP).Float64()
	p.value, _ = value.Round(DP).Float64()
	p.valueBOT, _ = valueBot.Round(DP).Float64()
	p.valueSLD, _ = valueSld.Round(DP).Float64()
	p.netValue, _ = netValue.Round(DP).Float64()
	p.netValueBOT, _ = netValueBot.Round(DP).Float64()
	p.netValueSLD, _ = netValueSld.Round(DP).Float64()
	p.commission, _ = commission.Round(DP).Float64()
	p.exchangeFee, _ = exchangeFee.Round(DP).Float64()
	p.cost, _ = cost.Round(DP).Float64()
	p.costBasis, _ = costBasis.Round(DP).Float64()
	p.realProfitLoss, _ = realProfitLoss.Round(DP).Float64()

	p.updateValue(fill.GetPrice())
}

// internal function to updates the current market value and profit/loss of a position
func (p *position) updateValue(l float64) {
	// convert to internally used decimal numbers
	latest := decimal.NewFromFloat(l)
	qty := decimal.New(p.qty, 0)
	costBasis := decimal.NewFromFloat(p.costBasis)

	// update market value
	marketPrice := latest
	p.marketPrice, _ = marketPrice.Round(DP).Float64()
	// abs(qty) * current
	marketValue := qty.Abs().Mul(latest)
	p.marketValue, _ = marketValue.Round(DP).Float64()

	// qty * current - costBasis
	unrealProfitLoss := qty.Mul(latest).Sub(costBasis)
	p.unrealProfitLoss, _ = unrealProfitLoss.Round(DP).Float64()

	realProfitLoss := decimal.NewFromFloat(p.realProfitLoss)
	totalProfitLoss := realProfitLoss.Add(unrealProfitLoss)
	p.totalProfitLoss, _ = totalProfitLoss.Round(DP).Float64()
}
