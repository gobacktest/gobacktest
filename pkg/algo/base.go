package algo

import (
	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

// TrueAlgo is a base Algo which always returns true
type TrueAlgo struct {
	bt.Algo
}

// Run runs the algo, returns true
func (a TrueAlgo) Run(s bt.StrategyHandler) (bool, error) {
	return true, nil
}

// FalseAlgo is an algo which always fails
type FalseAlgo struct {
	bt.Algo
}

// Run runs the algo, returns false
func (a FalseAlgo) Run(s bt.StrategyHandler) (bool, error) {
	return false, nil
}
