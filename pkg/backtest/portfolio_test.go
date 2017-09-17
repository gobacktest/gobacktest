package backtest

import (
	"testing"
)

func TestIsInvested(t *testing.T) {
	// testCases is a table for testing
	var testCases = []struct {
		msg       string    // error message
		symbol    string    // input string
		portfolio Portfolio // input portfolio
		expPos    position  // expected position return
		expOk     bool      // expected bool return
	}{
		{"Portfolio is empty:",
			"TEST.DE",
			Portfolio{},
			position{},
			false,
		},
		{"Portfolio should be invested with long:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 10},
				},
			},
			position{},
			true,
		},
		{"Portfolio should be invested with short:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: -10},
				},
			},
			position{},
			true,
		},
		{"Portfolio should not be invested:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 0},
				},
			},
			position{},
			false,
		},
	}

	for _, tc := range testCases {
		pos, ok := tc.portfolio.IsInvested(tc.symbol)
		if (pos != tc.expPos) && (ok != tc.expOk) {
			t.Errorf("%v\nIsInvested(%v): \nexpected %#v %v, \nactual %#v %v", tc.msg, tc.symbol, tc.expOk, tc.expPos, pos, ok)
		}
	}
}

func TestIsLong(t *testing.T) {
	// testCases is a table for testing
	var testCases = []struct {
		msg       string    // error message
		symbol    string    // input string
		portfolio Portfolio // input portfolio
		expPos    position  // expected position return
		expOk     bool      // expected bool return
	}{
		{"Portfolio is empty:",
			"TEST.DE",
			Portfolio{},
			position{},
			false,
		},
		{"Portfolio should be long:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 10},
				},
			},
			position{},
			true,
		},
		{"Portfolio is short:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: -10},
				},
			},
			position{},
			false,
		},
		{"Portfolio is not invested:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 0},
				},
			},
			position{},
			false,
		},
	}

	for _, tc := range testCases {
		pos, ok := tc.portfolio.IsLong(tc.symbol)
		if (pos != tc.expPos) && (ok != tc.expOk) {
			t.Errorf("%v\nIsLong(%v): \nexpected %#v %v, \nactual %#v %v", tc.msg, tc.symbol, tc.expOk, tc.expPos, pos, ok)
		}
	}
}

func TestIsShort(t *testing.T) {
	// testCases is a table for testing
	var testCases = []struct {
		msg       string    // error message
		symbol    string    // input string
		portfolio Portfolio // input portfolio
		expPos    position  // expected position return
		expOk     bool      // expected bool return
	}{
		{"Portfolio is empty:",
			"TEST.DE",
			Portfolio{},
			position{},
			false,
		},
		{"Portfolio is long:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 10},
				},
			},
			position{},
			false,
		},
		{"Portfolio should be short:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: -10},
				},
			},
			position{},
			true,
		},
		{"Portfolio is not invested:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 0},
				},
			},
			position{},
			false,
		},
	}

	for _, tc := range testCases {
		pos, ok := tc.portfolio.IsShort(tc.symbol)
		if (pos != tc.expPos) && (ok != tc.expOk) {
			t.Errorf("%v\nIsShort(%v): \nexpected %#v %v, \nactual %#v %v", tc.msg, tc.symbol, tc.expOk, tc.expPos, pos, ok)
		}
	}
}
