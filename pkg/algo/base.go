package algo

import (
	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

// TrueAlgo is a base Algo which always returns true
type TrueAlgo struct {
	bt.Algo
}

// Run runs the algo, returns true
func (a TrueAlgo) Run() bool {
	return true
}

// FalseAlgo is an algo which always fails
type FalseAlgo struct {
	bt.Algo
}

// Run runs the algo, returns false
func (a FalseAlgo) Run() bool {
	return false
}
