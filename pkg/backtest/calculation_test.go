package backtest

import (
	"errors"
	"reflect"
	"testing"
)

func TestCalculateSMA(t *testing.T) {
	// testCases is a table for testing the calculation fof a simple moving average
	var testCases = []struct {
		msg      string
		interval int
		data     []DataEventHandler // slice of data points
		expSMA   float64            // expected Position
		expErr   error              // expected error
	}{
		{msg: "Test simple SMA:",
			interval: 2,
			data: []DataEventHandler{
				&Bar{Event: Event{Symbol: "TEST.DE"}, Close: 5},
				&Bar{Event: Event{Symbol: "TEST.DE"}, Close: 10},
			},
			expSMA: 7.5,
			expErr: nil,
		},
		{msg: "Test SMA with lots of data points:",
			interval: 3,
			data: []DataEventHandler{
				&Bar{Event: Event{Symbol: "TEST.DE"}, Close: 5},
				&Bar{Event: Event{Symbol: "TEST.DE"}, Close: 4},
				&Bar{Event: Event{Symbol: "TEST.DE"}, Close: 6},
				&Bar{Event: Event{Symbol: "TEST.DE"}, Close: 2},
				&Bar{Event: Event{Symbol: "TEST.DE"}, Close: 10},
			},
			expSMA: 6,
			expErr: nil,
		},
		{msg: "Test SMA not enough data points:",
			interval: 2,
			data:     []DataEventHandler{},
			expSMA:   0,
			expErr:   errors.New("not enough data points to calculate SMA2"),
		},
	}

	for _, tc := range testCases {
		sma, err := CalculateSMA(tc.interval, tc.data)
		if (sma != tc.expSMA) || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("%s CalculateSMA(): \nexpected %#v %v, \nactual   %#v %v", tc.msg, tc.expSMA, tc.expErr, sma, err)
		}
	}
}
