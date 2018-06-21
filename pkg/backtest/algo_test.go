package backtest

import (
	"testing"
)

// mocking algos
// BasicAlgo is a base Algo which always returns true
type TrueAlgo struct {
}

// Run runs the algo, returns true
func (a TrueAlgo) Run() bool {
	return true
}

// FalseAlgo is an algo which always fails
type FalseAlgo struct {
}

// Run runs the algo, returns false
func (a FalseAlgo) Run() bool {
	return false
}

func TestAlgoStackRun(t *testing.T) {
	var testCases = []struct {
		msg string
		as  AlgoStack
		exp bool
	}{
		{"Empty AlgoStack:",
			AlgoStack{},
			true,
		},
		{"Simple AlgoStack with single true Algo:",
			AlgoStack{
				stack: []AlgoHandler{
					TrueAlgo{},
				},
			},
			true,
		},
		{"Simple AlgoStack with single false Algo:",
			AlgoStack{
				stack: []AlgoHandler{
					FalseAlgo{},
				},
			},
			false,
		},
		{"AlgoStack with multiple true Algos:",
			AlgoStack{
				stack: []AlgoHandler{
					TrueAlgo{},
					TrueAlgo{},
					TrueAlgo{},
				},
			},
			true,
		},
		{"AlgoStack with multiple Algos, one is false:",
			AlgoStack{
				stack: []AlgoHandler{
					TrueAlgo{},
					FalseAlgo{},
					TrueAlgo{},
				},
			},
			false,
		},
		{"AlgoStack with true sub AlgoStack:",
			AlgoStack{
				stack: []AlgoHandler{
					TrueAlgo{},
					AlgoStack{
						stack: []AlgoHandler{
							TrueAlgo{},
							TrueAlgo{},
						},
					},
					TrueAlgo{},
				},
			},
			true,
		},
		{"AlgoStack with false sub AlgoStack:",
			AlgoStack{
				stack: []AlgoHandler{
					TrueAlgo{},
					AlgoStack{
						stack: []AlgoHandler{
							TrueAlgo{},
							FalseAlgo{},
						},
					},
					TrueAlgo{},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		ok := tc.as.Run()
		if ok != tc.exp {
			t.Errorf("%v Run(): \nexpected %#v, \nactual %#v", tc.msg, tc.exp, ok)
		}
	}
}
