package backtest

// AlgoHandler defines the base algorythm functionality.
type AlgoHandler interface {
	Run(StrategyHandler) (bool, error)
}

// Algo is a base algo structure, implements AlgoHandler
type Algo struct {
	// determines if the algo runs always, even if a preceding algo fails
	RunAlways bool
}

// Run implements the Algo interface.
func (a Algo) Run(_ StrategyHandler) (bool, error) {
	return true, nil
}

// AlgoStack represents a single stack of algos.
type AlgoStack struct {
	Algo
	stack []AlgoHandler
}

// Run implements the Algo interface on the AlgoStack, which makes it itself an Algo.
func (as AlgoStack) Run(s StrategyHandler) (bool, error) {
	for _, algo := range as.stack {
		if ok, err := algo.Run(s); !ok {
			return false, err
		}
	}
	return true, nil
}
