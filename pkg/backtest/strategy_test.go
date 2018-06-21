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
