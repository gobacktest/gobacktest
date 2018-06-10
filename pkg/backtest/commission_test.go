package backtest

import (
	"reflect"
	"testing"
)

func TestFixedCommission(t *testing.T) {
	var testCases = []struct {
		msg    string
		c      CommissionHandler
		qty    float64
		price  float64
		expCom float64
		expErr error
	}{
		{"testing fixed commission for empty parameters / no trade:",
			&FixedCommission{Commission: 10},
			0, 0,
			0, nil,
		},
		{"testing fixed commission for zero commission:",
			&FixedCommission{Commission: 0},
			100, 100,
			0, nil,
		},
		{"testing fixed commission for fixed commission:",
			&FixedCommission{Commission: 10},
			100, 100,
			10, nil,
		},
	}

	for _, tc := range testCases {
		commission, err := tc.c.Calculate(tc.qty, tc.price)
		if commission != tc.expCom || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("%v Calculate(): \nexpected %#v, \nactual %#v", tc.msg, tc.expCom, commission)
		}
	}
}

func TestTresholdFixedCommission(t *testing.T) {
	var testCases = []struct {
		msg    string
		c      CommissionHandler
		qty    float64
		price  float64
		expCom float64
		expErr error
	}{
		{"testing treshold commission for empty parameters / no trade:",
			&TresholdFixedCommission{Commission: 10, MinValue: 1000},
			0, 0,
			0, nil,
		},
		{"testing treshold commission for value below minimum:",
			&TresholdFixedCommission{Commission: 10, MinValue: 1000},
			10, 10,
			0, nil,
		},
		{"testing treshold commission for value above minimum:",
			&TresholdFixedCommission{Commission: 10, MinValue: 1000},
			100, 100,
			10, nil,
		},
	}

	for _, tc := range testCases {
		commission, err := tc.c.Calculate(tc.qty, tc.price)
		if commission != tc.expCom || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("%v Calculate(): \nexpected %#v, \nactual %#v", tc.msg, tc.expCom, commission)
		}
	}
}

func TestPercentageCommission(t *testing.T) {
	var testCases = []struct {
		msg    string
		c      CommissionHandler
		qty    float64
		price  float64
		expCom float64
		expErr error
	}{
		{"testing percentage commission for empty parameters / no trade:",
			&PercentageCommission{Commission: 0.01},
			0, 0,
			0, nil,
		},
		{"testing percentage commission for value:",
			&PercentageCommission{Commission: 0.01},
			10, 10,
			1, nil,
		},
	}

	for _, tc := range testCases {
		commission, err := tc.c.Calculate(tc.qty, tc.price)
		if commission != tc.expCom || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("%v Calculate(): \nexpected %#v, \nactual %#v", tc.msg, tc.expCom, commission)
		}
	}
}

func TestValueCommission(t *testing.T) {
	var testCases = []struct {
		msg    string
		c      CommissionHandler
		qty    float64
		price  float64
		expCom float64
		expErr error
	}{
		{"testing value commission for empty parameters / no trade:",
			&ValueCommission{Commission: 0.01, MinCommission: 10, MaxCommission: 100},
			0, 0,
			0, nil,
		},
		{"testing value commission for value below minimum:",
			&ValueCommission{Commission: 0.01, MinCommission: 10, MaxCommission: 100},
			10, 10,
			10, nil,
		},
		{"testing value commission for value above maximum:",
			&ValueCommission{Commission: 0.01, MinCommission: 10, MaxCommission: 100},
			1000, 1000,
			100, nil,
		},
		{"testing value commission for value between minimum and maximum:",
			&ValueCommission{Commission: 0.01, MinCommission: 10, MaxCommission: 100},
			100, 50,
			50, nil,
		},
	}

	for _, tc := range testCases {
		commission, err := tc.c.Calculate(tc.qty, tc.price)
		if commission != tc.expCom || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("%v Calculate(): \nexpected %#v, \nactual %#v", tc.msg, tc.expCom, commission)
		}
	}
}
