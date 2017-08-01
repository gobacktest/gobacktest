package internal

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
	OnFill()
}

// Portfolio represent a simple portfolio struct.
type Portfolio struct {
	Cash      int64
	positions map[string]Position
}

// OnSignal handles an incomming signal event
func (p *Portfolio) OnSignal(s SignalEvent) (order OrderEvent, ok bool) {
	return order, true
}

// OnFill handles an incomming fill event
func (p *Portfolio) OnFill() {

}
