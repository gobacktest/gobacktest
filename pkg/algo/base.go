package algo

import (
	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

// boolAlgo is a base Algo which return a set bool value
type boolAlgo struct {
	bt.Algo
	bool
}

// BoolAlgo returns a simple true/false algo ready to use.
func BoolAlgo(b bool) bt.AlgoHandler {
	return &boolAlgo{bool: b}
}

// Run runs the algo, returns the bool value of the algo
func (a boolAlgo) Run(s bt.StrategyHandler) (bool, error) {
	return a.bool, nil
}
