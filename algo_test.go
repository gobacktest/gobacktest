package gobacktest

import (
	"testing"
)

// MockAlgo is a mocked base Algo
type MockAlgo struct {
	Algo
	ret bool
}

// Run runs the algo
func (ma MockAlgo) Run(s StrategyHandler) (bool, error) {
	return ma.ret, nil
}

func TestAlgoRun(t *testing.T) {
	var testCases = []struct {
		msg      string
		algo     Algo
		strategy *Strategy
		expOk    bool
		expErr   error
	}{
		{"simple true test:",
			Algo{}, &Strategy{},
			true, nil,
		},
	}

	for _, tc := range testCases {
		ok, err := tc.algo.Run(tc.strategy)
		if (ok != tc.expOk) || (err != tc.expErr) {
			t.Errorf("%v Run(): \nexpected %v %v, \nactual   %v %v", tc.msg, tc.expOk, tc.expErr, ok, err)
		}
	}
}

func TestAlgoStackRun(t *testing.T) {
	var testCases = []struct {
		msg      string
		as       AlgoStack
		strategy *Strategy
		expOk    bool
		expErr   error
	}{
		{"Empty AlgoStack:",
			AlgoStack{}, &Strategy{}, true, nil,
		},
		{"Simple AlgoStack with single true Algo:",
			AlgoStack{
				stack: []AlgoHandler{
					&MockAlgo{ret: true},
				},
			},
			&Strategy{}, true, nil,
		},
		{"Simple AlgoStack with single false Algo:",
			AlgoStack{
				stack: []AlgoHandler{
					&MockAlgo{ret: false},
				},
			},
			&Strategy{}, false, nil,
		},
		{"AlgoStack with multiple true Algos:",
			AlgoStack{
				stack: []AlgoHandler{
					&MockAlgo{ret: true},
					&MockAlgo{ret: true},
					&MockAlgo{ret: true},
				},
			},
			&Strategy{}, true, nil,
		},
		{"AlgoStack with multiple Algos, one is false:",
			AlgoStack{
				stack: []AlgoHandler{
					&MockAlgo{ret: true},
					&MockAlgo{ret: false},
					&MockAlgo{ret: true},
				},
			},
			&Strategy{}, false, nil,
		},
		{"AlgoStack with true sub AlgoStack:",
			AlgoStack{
				stack: []AlgoHandler{
					&MockAlgo{ret: true},
					&AlgoStack{
						stack: []AlgoHandler{
							&MockAlgo{ret: true},
							&MockAlgo{ret: true},
						},
					},
					&MockAlgo{ret: true},
				},
			},
			&Strategy{}, true, nil,
		},
		{"AlgoStack with false sub AlgoStack:",
			AlgoStack{
				stack: []AlgoHandler{
					&MockAlgo{ret: true},
					&AlgoStack{
						stack: []AlgoHandler{
							&MockAlgo{ret: true},
							&MockAlgo{ret: false},
						},
					},
					&MockAlgo{ret: true},
				},
			},
			&Strategy{}, false, nil,
		},
	}

	for _, tc := range testCases {
		ok, err := tc.as.Run(tc.strategy)
		if (ok != tc.expOk) || (err != tc.expErr) {
			t.Errorf("%v Run(): \nexpected %v %v, \nactual   %v %v", tc.msg, tc.expOk, tc.expErr, ok, err)
		}
	}
}
