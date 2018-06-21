package backtest

// AlgoHandler defines the base algorythm functionality.
type AlgoHandler interface {
	Run() bool
}

// AlgoStack represents a single stack of algos.
type AlgoStack struct {
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
