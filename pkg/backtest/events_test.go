package backtest

import (
	"testing"
)

func TestTickLatestPrice(t *testing.T) {
	var testCases = []struct {
		msg  string
		tick Tick
		exp  float64
	}{
		{"Empty Tick:",
			Tick{Bid: 0, Ask: 0},
			0,
		},
		{"Standard Tick:",
			Tick{Bid: 10, Ask: 5},
			7.5,
		},
	}

	for _, tc := range testCases {
		float := tc.tick.LatestPrice()
		if float != tc.exp {
			t.Errorf("%v LatestPrice(): \nexpected %#v, \nactual %#v", tc.msg, tc.exp, float)
		}
	}
}

func TestFillValue(t *testing.T) {
	var testCases = []struct {
		msg  string
		fill Fill
		exp  float64
	}{
		{"Empty Fill:",
			Fill{Qty: 0, Price: 0},
			0,
		},
		{"Standard Fill:",
			Fill{Qty: 10, Price: 5},
			50,
		},
	}

	for _, tc := range testCases {
		float := tc.fill.Value()
		if float != tc.exp {
			t.Errorf("%v Value(): \nexpected %#v, \nactual %#v", tc.msg, tc.exp, float)
		}
	}
}

func TestFillNetValue(t *testing.T) {
	var testCases = []struct {
		msg  string
		fill Fill
		exp  float64
	}{
		{"Empty BOT Fill:",
			Fill{Direction: "BOT", Qty: 0, Price: 0, Cost: 0},
			0,
		},
		{"Standard BOT Fill:",
			Fill{Direction: "BOT", Qty: 10, Price: 5, Cost: 5},
			55,
		},
		{"Empty SLD Fill:",
			Fill{Direction: "SLD", Qty: 0, Price: 0, Cost: 0},
			0,
		},
		{"Standard SLD Fill:",
			Fill{Direction: "SLD", Qty: 10, Price: 5, Cost: 5},
			45,
		},
	}

	for _, tc := range testCases {
		float := tc.fill.NetValue()
		if float != tc.exp {
			t.Errorf("%v NetValue(): \nexpected %#v, \nactual %#v", tc.msg, tc.exp, float)
		}
	}
}
