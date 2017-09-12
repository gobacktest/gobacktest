package internal

import (
	"testing"
)

func TestIsInvested(t *testing.T) {
	// testCases is a table for testing
	var testCases = []struct {
		msg       string    // error message
		symbol    string    // input string
		portfolio Portfolio // input portfolio
		expOk     bool      // expected return
	}{
		{"Portfolio is empty:",
			"TEST.DE",
			Portfolio{},
			false,
		},
		{"Portfolio should be invested with long:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 10},
				},
			},
			true,
		},
		{"Portfolio should be invested with short:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: -10},
				},
			},
			true,
		},
		{"Portfolio should not be invested:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 0},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		ok := tc.portfolio.IsInvested(tc.symbol)
		if ok != tc.expOk {
			t.Errorf("%v\nIsInvested(%v): \nexpected %#v, \nactual %#v", tc.msg, tc.symbol, tc.expOk, ok)
		}
	}
}

func TestIsLong(t *testing.T) {
	// testCases is a table for testing
	var testCases = []struct {
		msg       string    // error message
		symbol    string    // input string
		portfolio Portfolio // input portfolio
		expOk     bool      // expected return
	}{
		{"Portfolio is empty:",
			"TEST.DE",
			Portfolio{},
			false,
		},
		{"Portfolio should be long:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 10},
				},
			},
			true,
		},
		{"Portfolio is short:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: -10},
				},
			},
			false,
		},
		{"Portfolio is not invested:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 0},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		ok := tc.portfolio.IsLong(tc.symbol)
		if ok != tc.expOk {
			t.Errorf("%v\nIsLong(%v): \nexpected %#v, \nactual %#v", tc.msg, tc.symbol, tc.expOk, ok)
		}
	}
}

func TestIsShort(t *testing.T) {
	// testCases is a table for testing
	var testCases = []struct {
		msg       string    // error message
		symbol    string    // input string
		portfolio Portfolio // input portfolio
		expOk     bool      // expected return
	}{
		{"Portfolio is empty:",
			"TEST.DE",
			Portfolio{},
			false,
		},
		{"Portfolio is long:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 10},
				},
			},
			false,
		},
		{"Portfolio should be short:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: -10},
				},
			},
			true,
		},
		{"Portfolio is not invested:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 0},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		ok := tc.portfolio.IsShort(tc.symbol)
		if ok != tc.expOk {
			t.Errorf("%v\nIsShort(%v): \nexpected %#v, \nactual %#v", tc.msg, tc.symbol, tc.expOk, ok)
		}
	}
}
