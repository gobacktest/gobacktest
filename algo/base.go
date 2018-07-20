package algo

import (
	gbt "github.com/dirkolbrich/gobacktest"
)

// boolAlgo is a base Algo which return a set bool value
type boolAlgo struct {
	gbt.Algo
	bool
}

// BoolAlgo returns a simple true/false algo ready to use.
func BoolAlgo(b bool) gbt.AlgoHandler {
	return &boolAlgo{bool: b}
}

// Run runs the algo, returns the bool value of the algo
func (a boolAlgo) Run(s gbt.StrategyHandler) (bool, error) {
	return a.bool, nil
}
