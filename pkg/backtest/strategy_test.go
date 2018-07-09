package backtest

import (
	"reflect"
	"testing"
)

func TestSetData(t *testing.T) {
	var testCases = []struct {
		msg      string
		strategy *Strategy
		data     DataHandler
		exp      *Strategy
		expErr   error
	}{
		{"set data:",
			&Strategy{
				Node:  Node{name: "test", root: true},
				algos: AlgoStack{},
			},
			&Data{},
			&Strategy{
				Node:  Node{name: "test", root: true},
				algos: AlgoStack{},
				data:  &Data{},
			},
			nil,
		},
		{"set data with child strategy:",
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Strategy{
							Node:  Node{name: "sub", root: false},
							algos: AlgoStack{},
						},
					},
				},
				algos: AlgoStack{},
			},
			&Data{},
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Strategy{
							Node:  Node{name: "sub", root: false},
							algos: AlgoStack{},
							data:  &Data{},
						},
					},
				},
				algos: AlgoStack{},
				data:  &Data{},
			},
			nil,
		},
	}

	for _, tc := range testCases {
		err := tc.strategy.SetData(tc.data)
		if !reflect.DeepEqual(tc.strategy, tc.exp) || (err != tc.expErr) {
			t.Errorf("%v SetData(%s): \nexpected: %#v %#v, \nactual:   %#v %#v",
				tc.msg, tc.data, tc.exp, tc.expErr, tc.strategy, err)
		}
	}
}

func TestSetPortfolio(t *testing.T) {
	var testCases = []struct {
		msg       string
		strategy  *Strategy
		portfolio PortfolioHandler
		exp       *Strategy
		expErr    error
	}{
		{"set data:",
			&Strategy{
				Node:  Node{name: "test", root: true},
				algos: AlgoStack{},
			},
			&Portfolio{},
			&Strategy{
				Node:      Node{name: "test", root: true},
				algos:     AlgoStack{},
				portfolio: &Portfolio{},
			},
			nil,
		},
		{"set data with child strategy:",
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Strategy{
							Node:  Node{name: "sub", root: false},
							algos: AlgoStack{},
						},
					},
				},
				algos: AlgoStack{},
			},
			&Portfolio{},
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Strategy{
							Node:      Node{name: "sub", root: false},
							algos:     AlgoStack{},
							portfolio: &Portfolio{},
						},
					},
				},
				algos:     AlgoStack{},
				portfolio: &Portfolio{},
			},
			nil,
		},
	}

	for _, tc := range testCases {
		err := tc.strategy.SetPortfolio(tc.portfolio)
		if !reflect.DeepEqual(tc.strategy, tc.exp) || (err != tc.expErr) {
			t.Errorf("%v SetPortfolio(%s): \nexpected: %#v %#v, \nactual:   %#v %#v",
				tc.msg, tc.portfolio, tc.exp, tc.expErr, tc.strategy, err)
		}
	}
}

func TestNewStrategy(t *testing.T) {
	var testCases = []struct {
		msg  string
		name string
		exp  *Strategy
	}{
		{"setup new strategy:",
			"test",
			&Strategy{
				Node:  Node{name: "test", root: true},
				algos: AlgoStack{},
			},
		},
	}

	for _, tc := range testCases {
		strategy := NewStrategy(tc.name)
		if !reflect.DeepEqual(strategy, tc.exp) {
			t.Errorf("%v NewStrategy(%s): \nexpected: %#v, \nactual:   %#v",
				tc.msg, tc.name, tc.exp, strategy)
		}
	}
}

func TestStrategies(t *testing.T) {
	var testCases = []struct {
		msg      string
		strategy *Strategy
		exp      []StrategyHandler
		expOk    bool
	}{
		{"test no children:",
			&Strategy{
				Node:  Node{name: "test", root: true},
				algos: AlgoStack{},
			},
			nil,
			false,
		},
		{"test no sub strategy, only assets:",
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Asset{
							Node: Node{name: "asset", root: false},
						},
					},
				},
				algos: AlgoStack{},
			},
			nil,
			false,
		},
		{"test single sub strategy:",
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Strategy{
							Node:  Node{name: "sub", root: false},
							algos: AlgoStack{},
						},
					},
				},
				algos: AlgoStack{},
			},
			[]StrategyHandler{
				&Strategy{
					Node:  Node{name: "sub", root: false},
					algos: AlgoStack{},
				},
			},
			true,
		},
		{"test multiple sub strategies:",
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Strategy{
							Node:  Node{name: "subA", root: false},
							algos: AlgoStack{},
						},
						&Strategy{
							Node:  Node{name: "subB", root: false},
							algos: AlgoStack{},
						},
					},
				},
				algos: AlgoStack{},
			},
			[]StrategyHandler{
				&Strategy{
					Node:  Node{name: "subA", root: false},
					algos: AlgoStack{},
				},
				&Strategy{
					Node:  Node{name: "subB", root: false},
					algos: AlgoStack{},
				},
			},
			true,
		},
		{"test multiple sub strategies and multiple assets:",
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Strategy{
							Node:  Node{name: "subA", root: false},
							algos: AlgoStack{},
						},
						&Strategy{
							Node:  Node{name: "subB", root: false},
							algos: AlgoStack{},
						},
						&Asset{
							Node: Node{name: "assetA", root: false},
						},
						&Asset{
							Node: Node{name: "assetB", root: false},
						},
					},
				},
				algos: AlgoStack{},
			},
			[]StrategyHandler{
				&Strategy{
					Node:  Node{name: "subA", root: false},
					algos: AlgoStack{},
				},
				&Strategy{
					Node:  Node{name: "subB", root: false},
					algos: AlgoStack{},
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		strategies, ok := tc.strategy.Strategies()
		if !reflect.DeepEqual(strategies, tc.exp) || (ok != tc.expOk) {
			t.Errorf("%v Strategies(): \nexpected %#v %v, \nactual  %#v %v",
				tc.msg, tc.exp, tc.expOk, strategies, ok)
		}
	}
}

func TestAssets(t *testing.T) {
	var testCases = []struct {
		msg      string
		strategy *Strategy
		exp      []*Asset
		expOk    bool
	}{
		{"test no children:",
			&Strategy{
				Node:  Node{name: "test", root: true},
				algos: AlgoStack{},
			},
			nil,
			false,
		},
		{"test no assets, only sub strategy:",
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Strategy{
							Node:  Node{name: "subA", root: true},
							algos: AlgoStack{},
						},
					},
				},
				algos: AlgoStack{},
			},
			nil,
			false,
		},
		{"test single asset:",
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Asset{
							Node: Node{name: "asset", root: false},
						},
					},
				},
				algos: AlgoStack{},
			},
			[]*Asset{
				{
					Node: Node{name: "asset", root: false},
				},
			},
			true,
		},
		{"test multiple assetss:",
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Asset{
							Node: Node{name: "assetA", root: false},
						},
						&Asset{
							Node: Node{name: "assetB", root: false},
						},
					},
				},
				algos: AlgoStack{},
			},
			[]*Asset{
				{
					Node: Node{name: "assetA", root: false},
				},
				{
					Node: Node{name: "assetB", root: false},
				},
			},
			true,
		},
		{"test multiple sub strategies and multiple assets:",
			&Strategy{
				Node: Node{name: "test", root: true,
					children: []NodeHandler{
						&Strategy{
							Node:  Node{name: "subA", root: false},
							algos: AlgoStack{},
						},
						&Strategy{
							Node:  Node{name: "subB", root: false},
							algos: AlgoStack{},
						},
						&Asset{
							Node: Node{name: "assetA", root: false},
						},
						&Asset{
							Node: Node{name: "assetB", root: false},
						},
					},
				},
				algos: AlgoStack{},
			},
			[]*Asset{
				{
					Node: Node{name: "assetA", root: false},
				},
				{
					Node: Node{name: "assetB", root: false},
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		strategies, ok := tc.strategy.Assets()
		if !reflect.DeepEqual(strategies, tc.exp) || (ok != tc.expOk) {
			t.Errorf("%v Strategies(): \nexpected %#v %v, \nactual  %#v %v",
				tc.msg, tc.exp, tc.expOk, strategies, ok)
		}
	}
}
func TestStrategySingleSetAlgo(t *testing.T) {
	testStrategy := &Strategy{
		algos: AlgoStack{
			stack: []AlgoHandler{
				&MockAlgo{ret: true},
			},
		},
	}
	strategy := &Strategy{}
	strategy = strategy.SetAlgo(&MockAlgo{ret: true})

	if !reflect.DeepEqual(strategy, testStrategy) {
		t.Errorf("set single algo SetAlgo(): \nexpected %#v, \nactual %#v", testStrategy, strategy)
	}
}

func TestStrategyMultipleSetAlgo(t *testing.T) {
	testStrategy := &Strategy{
		algos: AlgoStack{
			stack: []AlgoHandler{
				&MockAlgo{ret: true},
				&MockAlgo{ret: true},
				&MockAlgo{ret: false},
			},
		},
	}
	strategy := &Strategy{}
	strategy = strategy.SetAlgo(
		&MockAlgo{ret: true},
		&MockAlgo{ret: true},
		&MockAlgo{ret: false},
	)

	if !reflect.DeepEqual(strategy, testStrategy) {
		t.Errorf("set multiple algos SetAlgo(): \nexpected %#v, \nactual %#v", testStrategy, strategy)
	}
}
