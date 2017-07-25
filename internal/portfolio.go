package internal

// PortfolioHandler interface is the basic building block for a portfolio.
type PortfolioHandler interface {
	updateSignal()
	updateFill()
}

// SimplePortfolio represent a simple portfolio struct.
type SimplePortfolio struct {
	// bars   []Bar
	events []EventHandler
	Cash   int64
}

func (sp SimplePortfolio) updateSignal() {

}

func (sp SimplePortfolio) updateFill() {

}
