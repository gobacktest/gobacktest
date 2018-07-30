package algo

import (
	gbt "github.com/dirkolbrich/gobacktest"
)

type ifAlgo struct {
	gbt.Algo
	condition, action gbt.AlgoHandler
}

// If triggers an action algo if a condition algo is true.
func If(condition, action gbt.AlgoHandler) gbt.AlgoHandler {
	return &ifAlgo{
		condition: condition,
		action:    action,
	}
}

// Run runs the algo, returns the bool value of the algo
func (algo ifAlgo) Run(s gbt.StrategyHandler) (bool, error) {
	ok, err := algo.condition.Run(s)
	if err != nil {
		return false, err
	}

	if ok {
		ok, err = algo.action.Run(s)
		if err != nil {
			return false, err
		}
	}

	return true, nil // always return true
}

type andAlgo struct {
	gbt.Algo
	a, b gbt.AlgoHandler
}

// And compares if both condition algos are true.
func And(a, b gbt.AlgoHandler) gbt.AlgoHandler {
	return &andAlgo{
		a: a,
		b: b,
	}
}

// Run runs the algo, returns the bool value of the algo.
func (algo andAlgo) Run(s gbt.StrategyHandler) (bool, error) {
	okA, err := algo.a.Run(s)
	if err != nil {
		return false, err
	}

	okB, err := algo.b.Run(s)
	if err != nil {
		return false, err
	}

	if !okA || !okB {
		return false, nil
	}

	return true, nil
}

type orAlgo struct {
	gbt.Algo
	a, b gbt.AlgoHandler
}

// Or compares if any or both of two algos are true.
func Or(a, b gbt.AlgoHandler) gbt.AlgoHandler {
	return &orAlgo{
		a: a,
		b: b,
	}
}

// Run runs the algo, returns the bool value of the algo.
func (algo orAlgo) Run(s gbt.StrategyHandler) (bool, error) {
	okA, err := algo.a.Run(s)
	if err != nil {
		return false, err
	}

	okB, err := algo.b.Run(s)
	if err != nil {
		return false, err
	}

	if !okA && !okB {
		return false, nil
	}

	return true, nil
}

type xorAlgo struct {
	gbt.Algo
	a, b gbt.AlgoHandler
}

// Xor compares if any one but only one of two algos is true.
func Xor(a, b gbt.AlgoHandler) gbt.AlgoHandler {
	return &xorAlgo{
		a: a,
		b: b,
	}
}

// Run runs the algo, returns the bool value of the algo.
func (algo xorAlgo) Run(s gbt.StrategyHandler) (bool, error) {
	okA, err := algo.a.Run(s)
	if err != nil {
		return false, err
	}

	okB, err := algo.b.Run(s)
	if err != nil {
		return false, err
	}

	if (!okA && !okB) || (okA && okB) {
		return false, nil
	}

	return true, nil
}
