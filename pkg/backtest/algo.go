package backtest

// AlgoHandler defines the base algorythm functionality.
type AlgoHandler interface {
	Run() bool
}

// Algo is a base algo structure, implements AlgoHandler
type Algo struct {
	// determines if the algo runs always, even if a preceding algo fails
	RunAlways bool
}

// Run implements the Algo interface.
func (a Algo) Run() bool {
	return true
}

// AlgoStack represents a single stack of algos.
type AlgoStack struct {
	Algo
	stack []AlgoHandler
}

// Run implements the Algo interface on the AlgoStack, which makes it itself an Algo.
func (as AlgoStack) Run() bool {
	for _, algo := range as.stack {
		if !algo.Run() {
			return false
		}
	}
	return true
}
