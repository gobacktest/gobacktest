package backtest

// AlgoHandler defines the base algorythm functionality.
type AlgoHandler interface {
	Run(StrategyHandler) (bool, error)
	Always() bool
	SetAlways()
}

// Algo is a base algo structure, implements AlgoHandler
type Algo struct {
	// determines if the algo runs always, even if a preceding algo fails
	runAlways bool
}

// Run implements the Algo interface.
func (a Algo) Run(_ StrategyHandler) (bool, error) {
	return true, nil
}

// Always returns the runAlways property.
func (a Algo) Always() bool {
	return a.runAlways
}

// SetAlways set the runAlways property.
func (a *Algo) SetAlways() {
	a.runAlways = true
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

// RunAlways set the runAlways property on the AlgoHandler
func RunAlways(a AlgoHandler) AlgoHandler {
	a.SetAlways()
	return a
}
