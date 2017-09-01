package internal

import "github.com/dirkolbrich/gobacktest/internal/utils"

// PortfolioHandler is the combined interface building block for a portfolio.
type PortfolioHandler interface {
	OnSignaler
	OnFiller
}

// OnSignaler as an intercafe for the OnSignal method
type OnSignaler interface {
	OnSignal(SignalEvent, DataEvent) (OrderEvent, error)
}

// OnFiller as an intercafe for the OnFill method
type OnFiller interface {
	OnFill(FillEvent) (FillEvent, error)
}

// Portfolio represent a simple portfolio struct.
type Portfolio struct {
	Cash         float64
	holdings     map[string]position
	transactions []FillEvent
	sizeManager  SizeHandler
	riskManager  RiskHandler
}

// SetSizeManager sets the size manager to be used with the portfolio
func (p *Portfolio) SetSizeManager(size SizeHandler) {
	p.sizeManager = size
}

// SetRiskManager sets the risk manager to be used with the portfolio
func (p *Portfolio) SetRiskManager(risk RiskHandler) {
	p.riskManager = risk
}

// OnSignal handles an incomming signal event
func (p *Portfolio) OnSignal(signal SignalEvent, current DataEvent) (OrderEvent, error) {
	// log.Printf("Portfolio receives Signal: %#v \n", s)

	// set order action
	var action string
	switch signal.Direction() {
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

	initialOrder := &order{
		event: event{
			timestamp: signal.Timestamp(),
			symbol:    signal.Symbol(),
		},
		direction: action,
		// Qty should be set by PositionSizer
		orderType: orderType,
		limit:     limit,
	}

	sizedOrder, err := p.sizeManager.SizeOrder(initialOrder, current, p.holdings)
	if err != nil {
	}

	order, err := p.riskManager.EvaluateOrder(sizedOrder, current, p.holdings)
	if err != nil {
	}

	return order, nil
}

// OnFill handles an incomming fill event
func (p *Portfolio) OnFill(fill FillEvent) (FillEvent, error) {
	// log.Printf("Portfolio receives Fill: %#v \n", f)

	// Check for nil map, else initialise the map
	if p.holdings == nil {
		p.holdings = make(map[string]position)
	}

	// check if portfolio has already a holding of the symbol from this fill
	if pos, ok := p.holdings[fill.Symbol()]; ok {
		// log.Printf("holding to this symbol exists: %+v \n", pos)
		// update existing Position
		pos.Update(fill)
	} else {
		// log.Println("No holding to this transaction")
		// create new position
		pos := position{}
		pos.Create(fill)
		p.holdings[fill.Symbol()] = pos
	}

	// update cash
	if fill.Direction() == "BOT" {
		p.Cash = utils.Round(p.Cash-fill.Net(), 3)
	} else {
		// direction is "SLD"
		p.Cash = utils.Round(p.Cash+fill.Net(), 3)
	}

	// add to transactions
	p.transactions = append(p.transactions, fill)

	return fill, nil
}
