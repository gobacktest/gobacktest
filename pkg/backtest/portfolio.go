package backtest

import (
	"github.com/shopspring/decimal"
)

// PortfolioHandler is the combined interface building block for a portfolio.
type PortfolioHandler interface {
	OnSignaler
	OnFiller
	Investor
	Updater
	Casher
	Valuer
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

// Updater handles the updating of the portfolio on data events
type Updater interface {
	Update(DataEventHandler)
}

// Casher handles basic portolio info
type Casher interface {
	SetInitialCash(float64)
	InitialCash() float64
	SetCash(float64)
	Cash() float64
}

// Valuer returns the values of the portfolio
type Valuer interface {
	Value() float64
}

// Portfolio represent a simple portfolio struct.
type Portfolio struct {
	initialCash  float64
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

	initialOrder := &Order{
		Event: Event{
			Timestamp: signal.GetTime(),
			Symbol:    signal.GetSymbol(),
		},
		Direction: signal.GetDirection(),
		// Qty should be set by PositionSizer
		OrderType: orderType,
		Limit:     limit,
	}

	// fetch latest known price for the symbol
	latest := data.Latest(signal.GetSymbol())

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
	// Check for nil map, else initialise the map
	if p.holdings == nil {
		p.holdings = make(map[string]position)
	}

	// fmt.Printf("holdings before: %+v\n", p.holdings)

	// check if portfolio has already a holding of the symbol from this fill
	if pos, ok := p.holdings[fill.GetSymbol()]; ok {
		// log.Printf("holding to this symbol exists: %+v \n", pos)
		// update existing Position
		pos.Update(fill)
		p.holdings[fill.GetSymbol()] = pos
	} else {
		// log.Println("No holding to this transaction")
		// create new position
		pos := position{}
		pos.Create(fill)
		p.holdings[fill.GetSymbol()] = pos
	}

	// fmt.Printf("holdings after: %+v\n", p.holdings)

	// update cash
	if fill.GetDirection() == "BOT" {
		p.cash = p.cash - fill.NetValue()
	} else {
		// direction is "SLD"
		p.cash = p.cash + fill.NetValue()
	}

	// add fill to transactions
	p.transactions = append(p.transactions, fill)

	return fill, nil
}

// IsInvested checks if the portfolio has an open position on the given symbol
func (p Portfolio) IsInvested(symbol string) (position, bool) {
	pos, ok := p.holdings[symbol]
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

// Update updates the holding on a data event
func (p *Portfolio) Update(d DataEventHandler) {
	if pos, ok := p.IsInvested(d.GetSymbol()); ok {
		pos.UpdateValue(d)
		p.holdings[d.GetSymbol()] = pos
	}
}

// SetInitialCash sets the initial cash value of the portfolio
func (p *Portfolio) SetInitialCash(initial float64) {
	p.initialCash = initial
}

// InitialCash returns the initial cash value of the portfolio
func (p Portfolio) InitialCash() float64 {
	return p.initialCash
}

// SetCash sets the current cash value of the portfolio
func (p *Portfolio) SetCash(cash float64) {
	p.cash = cash
}

// Cash returns the current cash value of the portfolio
func (p Portfolio) Cash() float64 {
	return p.cash
}

// Value return the current vlue of the portfolio
func (p Portfolio) Value() float64 {
	holdingValue := decimal.NewFromFloat(0)
	for _, pos := range p.holdings {
		marketValue := decimal.NewFromFloat(pos.marketValue)
		holdingValue = holdingValue.Add(marketValue)
	}

	cash := decimal.NewFromFloat(p.cash)
	value, _ := cash.Add(holdingValue).Round(4).Float64()
	return value
}
