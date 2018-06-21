package backtest

import (
	"reflect"
	"testing"
)

func TestNewStrategy(t *testing.T) {
	var testCases = []struct {
		msg  string
		name string
		exp  *Strategy
	}{
		{"setup new strategy:",
			"test",
			&Strategy{
				Node{name: "test"},
				AlgoStack{},
			},
		},
	}

	for _, tc := range testCases {
		strategy := NewStrategy(tc.name)
		if !reflect.DeepEqual(strategy, tc.exp) {
			t.Errorf("%v NewStrategy(%s): \nexpected %#v, \nactual %#v",
				tc.msg, tc.name, tc.exp, strategy)
		}
	}
}

func TestStrategySingleSetAlgo(t *testing.T) {
	testStrategy := &Strategy{
		algos: AlgoStack{
			stack: []AlgoHandler{
				&TrueAlgo{},
			},
		},
	}
	strategy := &Strategy{}
	strategy = strategy.SetAlgo(&TrueAlgo{})

	if !reflect.DeepEqual(strategy, testStrategy) {
		t.Errorf("set single algo SetAlgo(): \nexpected %#v, \nactual %#v", testStrategy, strategy)
	}
}

func TestStrategyMultipleSetAlgo(t *testing.T) {
	testStrategy := &Strategy{
		algos: AlgoStack{
			stack: []AlgoHandler{
				&TrueAlgo{},
				&TrueAlgo{},
				&FalseAlgo{},
			},
		},
	}
	strategy := &Strategy{}
	strategy = strategy.SetAlgo(
		&TrueAlgo{},
		&TrueAlgo{},
		&FalseAlgo{},
	)

	if !reflect.DeepEqual(strategy, testStrategy) {
		t.Errorf("set single algo SetAlgo(): \nexpected %#v, \nactual %#v", testStrategy, strategy)
	}
}
