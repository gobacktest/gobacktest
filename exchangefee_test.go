package gobacktest

import (
	"reflect"
	"testing"
)

func TestFixedExchangeFee(t *testing.T) {
	var testCases = []struct {
		msg    string
		e      ExchangeFeeHandler
		expFee float64
		expErr error
	}{
		{"testing fixed exchange fee for zero fee :",
			&FixedExchangeFee{ExchangeFee: 0},
			0, nil,
		},
		{"testing fixed exchange fee for fixed fee:",
			&FixedExchangeFee{ExchangeFee: 1.0},
			1, nil,
		},
	}

	for _, tc := range testCases {
		exchangeFee, err := tc.e.Fee()
		if exchangeFee != tc.expFee || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("%v Fee(): \nexpected %#v, \nactual %#v", tc.msg, tc.expFee, exchangeFee)
		}
	}
}
