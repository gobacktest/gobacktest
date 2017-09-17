package backtest

import (
	"reflect"
	"testing"
)

func TestEvaluateOrder(t *testing.T) {
	// set RiskHandler
	var r = &Risk{}

	// orderEventTests is a table for testing parsing bar data into a BarEvent
	var testCases = []struct {
		order     OrderEvent          // OrderEvent input
		data      DataEvent           // DataEvent input
		positions map[string]position // the portfolio holdings
		expOrder  OrderEvent          // expected FillEvent return
		expErr    error               // expected error output
	}{
		{
			&order{},
			&bar{},
			map[string]position{},
			&order{},
			nil,
		},
	}

	for _, tc := range testCases {
		order, err := r.EvaluateOrder(tc.order, tc.data, tc.positions)
		if !reflect.DeepEqual(order, tc.expOrder) || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("ExecuteOrder(%v): \nexpected %+v %v, \nactual   %+v %v",
				tc.order, tc.expOrder, tc.expErr, order, err)
		}
	}
}
