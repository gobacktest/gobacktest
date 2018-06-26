package backtest

// StrategyHandler is a basic strategy interface
type StrategyHandler interface {
	SetData(d DataHandler) error
	SetPortfolio(p PortfolioHandler) error
	// CalculateSignal(DataEventHandler, DataHandler, PortfolioHandler) (SignalEvent, error)
	OnData(DataEventHandler) (SignalEvent, error)
}

// Strategy implements a sub node, used as a strategy building block
type Strategy struct {
	Node
	algos     AlgoStack
	data      DataHandler
	portfolio PortfolioHandler
}

// NewStrategy return a new strategy node ready to use
func NewStrategy(name string) *Strategy {
	var s = &Strategy{}
	s.SetName(name)
	s.SetRoot(true)
	return s
}

// SetData sets the data property
func (s *Strategy) SetData(data DataHandler) error {
	s.data = data

	// check for children and if one of those is a strategy
	if children, ok := s.Children(); ok {
		for _, child := range children {
			// type switch for child type
			switch c := child.(type) {
			case *Strategy:
				err := c.SetData(data)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// SetPortfolio sets the portfolio property
func (s *Strategy) SetPortfolio(portfolio PortfolioHandler) error {
	s.portfolio = portfolio

	// check for children and if one of those is a strategy
	if children, ok := s.Children(); ok {
		for _, child := range children {
			// type switch for child type
			switch c := child.(type) {
			case *Strategy:
				err := c.SetPortfolio(portfolio)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// OnData handles the single Event
func (s *Strategy) OnData(event DataEventHandler) (SignalEvent, error) {
	// create Signal
	se := &Signal{}

	// type switch for event type
	switch e := event.(type) {
	case *Bar:
		// fill Signal
		se.Event = Event{Timestamp: e.GetTime(), Symbol: e.GetSymbol()}
		se.Direction = "long"
	}

	return se, nil
}

// // CalculateSignal handles the single Event
// func (s *Strategy) CalculateSignal(event DataEventHandler, data DataHandler, p PortfolioHandler) (SignalEvent, error) {
// 	// create Signal
// 	se := &Signal{}

// 	// type switch for event type
// 	switch e := event.(type) {
// 	case *Bar:
// 		// fill Signal
// 		se.Event = Event{Timestamp: e.GetTime(), Symbol: e.GetSymbol()}
// 		se.Direction = "long"
// 	}

// 	return se, nil
// }

// SetAlgo sets the algo stack for the Strategy
func (s *Strategy) SetAlgo(algos ...AlgoHandler) *Strategy {
	for _, algo := range algos {
		s.algos.stack = append(s.algos.stack, algo)
	}
	return s
}

// Run the algos of this Strategy Node, overwrite base Node method functionality
func (s Strategy) Run() error {
	// run the algo stack
	s.algos.Run()

	// check for children and run their algos
	if children, ok := s.Children(); ok {
		for _, child := range children {
			child.Run()
		}
	}

	return nil
}
