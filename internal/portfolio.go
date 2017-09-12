package internal

import "github.com/dirkolbrich/gobacktest/internal/utils"

// PortfolioHandler is the combined interface building block for a portfolio.
type PortfolioHandler interface {
	OnSignaler
	OnFiller
	Investor
}

// OnSignaler as an intercafe for the OnSignal method
type OnSignaler interface {
	OnSignal(SignalEvent, DataHandler) (OrderEvent, error)
}

// OnFiller as an intercafe for the OnFill method
type OnFiller interface {
	OnFill(FillEvent, DataHandler) (FillEvent, error)
}

// Investor is an inteface to check if a portfolio has a position of a symbol
type Investor interface {
	IsInvested(string) (position, bool)
	IsLong(string) (position, bool)
	IsShort(string) (position, bool)
}

// Portfolio represent a simple portfolio struct.
type Portfolio struct {
	InitialCash  float64
	cash         float64
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
func (p *Portfolio) OnSignal(signal SignalEvent, data DataHandler) (OrderEvent, error) {
	// log.Printf("Portfolio receives Signal: %#v \n", s)

	// set order type
	orderType := "MKT" // default Market, should be set by risk manager
	var limit float64

	initialOrder := &order{
		event: event{
			timestamp: signal.Timestamp(),
			symbol:    signal.Symbol(),
		},
		direction: signal.Direction(),
		// Qty should be set by PositionSizer
		orderType: orderType,
		limit:     limit,
	}

	// fetch latest known price for the symbol
	latest := data.Latest(signal.Symbol())

	sizedOrder, err := p.sizeManager.SizeOrder(initialOrder, latest, p)
	if err != nil {
	}

	order, err := p.riskManager.EvaluateOrder(sizedOrder, latest, p.holdings)
	if err != nil {
	}

	return order, nil
}

// OnFill handles an incomming fill event
func (p *Portfolio) OnFill(fill FillEvent, data DataHandler) (FillEvent, error) {
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

	// before first trade, set cash
	if len(p.transactions) == 0 {
		p.cash = p.InitialCash
	}

	// update cash
	if fill.Direction() == "BOT" {
		p.cash = utils.Round(p.cash-fill.NetValue(), 3)
	} else {
		// direction is "SLD"
		p.cash = utils.Round(p.cash+fill.NetValue(), 3)
	}

	// add to transactions
	p.transactions = append(p.transactions, fill)

	return fill, nil
}

// IsInvested checks if the portfolio has an open position on the given symbol
func (p Portfolio) IsInvested(symbol string) (pos position, ok bool) {
	pos, ok = p.holdings[symbol]
	if ok && (pos.qty != 0) {
		return pos, true
	}
	return pos, false
}

// IsLong checks if the portfolio has an open long position on the given symbol
func (p Portfolio) IsLong(symbol string) (pos position, ok bool) {
	pos, ok = p.holdings[symbol]
	if ok && (pos.qty > 0) {
		return pos, true
	}
	return pos, false
}

// IsShort checks if the portfolio has an open short position on the given symbol
func (p Portfolio) IsShort(symbol string) (pos position, ok bool) {
	pos, ok = p.holdings[symbol]
	if ok && (pos.qty < 0) {
		return pos, true
	}
	return pos, false
}
