package backtest

// StrategyHandler is a basic strategy interface
type StrategyHandler interface {
	CalculateSignal(DataEventHandler, DataHandler, PortfolioHandler) (SignalEvent, error)
}

// Strategy implements a sub node, used as a strategy building block
type Strategy struct {
	Node
	algos AlgoStack
}

// NewStrategy return a new strategy node ready to use
func NewStrategy(name string) *Strategy {
	var s = &Strategy{}
	s.SetName(name)
	return s
}

// SetAlgo sets the algo stack for the Strategy
func (s *Strategy) SetAlgo(algos ...AlgoHandler) *Strategy {
	for _, algo := range algos {
		s.algos.stack = append(s.algos.stack, algo)
	}
	return s
}

// Run the algos of this Strategy Node, overwrite base Node method functionality
func (s Strategy) Run() {
	// run the algo stack
	s.algos.Run()

	// check for children and run their algos
	if children, ok := s.Children(); ok {
		for _, child := range children {
			child.Run()
		}
	}
}
