package gobacktest

import (
	"testing"
)

func TestTickPrice(t *testing.T) {
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
		float := tc.tick.Price()
		if float != tc.exp {
			t.Errorf("%v LatestPrice(): \nexpected %#v, \nactual %#v", tc.msg, tc.exp, float)
		}
	}
}
