package backtest

import (
	"reflect"
	"testing"
)

func TestGetEquityPoint(t *testing.T) {
	var statCases = map[string]Statistic{
		"multiple": {
			equity: []equityPoint{
				{equity: 100, equityReturn: 0.1},
				{equity: 110, equityReturn: 0.2},
				{equity: 120, equityReturn: 0.3},
			},
		},
		"single": {
			equity: []equityPoint{
				{equity: 150, equityReturn: 0.25},
			},
		},
		"empty": {
			equity: []equityPoint{},
		},
	}

	// define test cases struct
	type testCase struct {
		msg   string
		stat  Statistic
		expEP equityPoint
		expOk bool
	}

	// set up test cases for getting first equity point
	var testCasesFirst = []testCase{
		{"testing first for multiple entryPoints",
			statCases["multiple"],
			equityPoint{equity: 100, equityReturn: 0.1},
			true},
		{"testing first for single entryPoints",
			statCases["single"],
			equityPoint{equity: 150, equityReturn: 0.25},
			true},
		{"testing first for nil entryPoints",
			statCases["empty"],
			equityPoint{},
			false},
	}

	for _, tc := range testCasesFirst {
		ep, ok := tc.stat.firstEquityPoint()
		if !reflect.DeepEqual(ep, tc.expEP) || (ok != tc.expOk) {
			t.Errorf("%v firstEquityPoint(): \nexpected %#v %v, \nactual   %#v %v",
				tc.msg, tc.expEP, tc.expOk, ep, ok)
		}
	}

	// set up test cases for getting last equity point
	var testCasesLast = []testCase{
		{"testing last for multiple entryPoints",
			statCases["multiple"],
			equityPoint{equity: 120, equityReturn: 0.3},
			true},
		{"testing last for single entryPoints",
			statCases["single"],
			equityPoint{equity: 150, equityReturn: 0.25},
			true},
		{"testing last for nil entryPoints",
			statCases["empty"],
			equityPoint{},
			false},
	}

	for _, tc := range testCasesLast {
		ep, ok := tc.stat.lastEquityPoint()
		if !reflect.DeepEqual(ep, tc.expEP) || (ok != tc.expOk) {
			t.Errorf("%v firstEquityPoint(): \nexpected %#v %v, \nactual   %#v %v",
				tc.msg, tc.expEP, tc.expOk, ep, ok)
		}
	}
}
