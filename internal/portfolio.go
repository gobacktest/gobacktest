package internal

import (
	"log"
	"time"
)

// PortfolioHandler is the combined interface building block for a portfolio.
type PortfolioHandler interface {
	OnSignaler
	OnFiller
}

// OnSignaler as an intercafe for the OnSignal method
type OnSignaler interface {
	OnSignal(SignalEvent) (OrderEvent, bool)
}

// OnFiller as an intercafe for the OnFill method
type OnFiller interface {
	OnFill(FillEvent) (FillEvent, bool)
}

// Portfolio represent a simple portfolio struct.
type Portfolio struct {
	Cash        float64
	holdings    map[string]Position
	riskManager RiskHandler
}

// SetRiskManager sets the risk manager to to be used within the portfolio
func (p *Portfolio) SetRiskManager(risk RiskHandler) {
	p.riskManager = risk
}

// OnSignal handles an incomming signal event
func (p *Portfolio) OnSignal(s SignalEvent) (order OrderEvent, ok bool) {
	log.Printf("Portfolio receives Signal: %#v \n", s)

	// set order action
	var action string
	switch s.Direction {
	case "long":
		action = "buy"
	case "short":
		action = "sell"
	case "exit": // all shares should be sold or bought, depending on position
		action = "sell"
	}

	// set order type
	orderType := "market" // default, should be set by risk manager
	var limit float64

	initialOrder := OrderEvent{
		Timestamp: time.Now(),
		Symbol:    s.Symbol,
		Direction: action,
		Qty:       s.SuggestedQty,
		OrderType: orderType,
		Limit:     limit,
	}

	order, ok = p.riskManager.EvaluateOrder(initialOrder, p.holdings)

	return order, ok
}

// OnFill handles an incomming fill event
func (p *Portfolio) OnFill(f FillEvent) (fill FillEvent, ok bool) {
	log.Printf("Portfolio receives Fill: %#v \n", f)

	// Check for nil map, else initialise the map
	if p.holdings == nil {
		p.holdings = make(map[string]Position)
	}

	// check if portfolio has already a holding of the symbol from this fill
	if pos, ok := p.holdings[f.Symbol]; ok {
		log.Printf("holding to this symbol exists: %+v \n", pos)
		// update existing Position
		p.holdings[f.Symbol] = p.updatePosition(pos, f)
		p.Cash -= f.Cost
	} else {
		log.Println("No holding to this transaction")
		// create new Position
		p.holdings[f.Symbol] = p.createPosition(f)
	}

	return f, true
}

// create a new position
func (p *Portfolio) createPosition(f FillEvent) Position {
	pos := Position{}
	pos.timestamp = time.Now()
	pos.symbol = f.Symbol
	pos.qty = f.Qty
	pos.avgPrice = f.Price
	pos.value = float64(f.Qty) * f.Price

	pos.marketPrice = f.Price
	pos.marketValue = pos.value

	pos.commission = f.Commission
	pos.exchangeFee = f.ExchangeFee
	pos.cost = f.Cost

	pos.netValue = pos.value - pos.cost

	return pos
}

// update a new position
func (p *Portfolio) updatePosition(pos Position, f FillEvent) Position {
	pos.timestamp = time.Now()
	pos.avgPrice = (float64(pos.qty)*pos.avgPrice + float64(f.Qty)*f.Price) / float64(pos.qty+f.Qty)
	pos.qty += f.Qty
	pos.value = float64(pos.qty) * pos.avgPrice

	pos.marketPrice = f.Price
	pos.marketValue = pos.value

	pos.commission += f.Commission
	pos.exchangeFee += f.ExchangeFee
	pos.cost += f.Cost

	pos.netValue = pos.value - pos.cost

	return pos
}
