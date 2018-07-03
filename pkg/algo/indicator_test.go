package algo

import (
	// "fmt"
	"errors"
	"reflect"
	"testing"

	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

func TestCalculateSMA(t *testing.T) {
	var testCases = []struct {
		msg    string
		period int
		data   []bt.DataEventHandler
		expSMA float64
		expErr error
	}{
		{"test zero data point",
			5, []bt.DataEventHandler{},
			0, errors.New("not enough data points to calculate SMA5, only 0 data points given"),
		},
		{"test too few data points",
			5,
			[]bt.DataEventHandler{
				&bt.Bar{},
			},
			0, errors.New("not enough data points to calculate SMA5, only 1 data points given"),
		},
		{"test single data point",
			1,
			[]bt.DataEventHandler{
				&bt.Bar{Close: 10},
			},
			10, nil,
		},
		{"test exact data points",
			3,
			[]bt.DataEventHandler{
				&bt.Bar{Close: 2},
				&bt.Bar{Close: 4},
				&bt.Bar{Close: 3},
			},
			3, nil,
		},
		{"test more than needed data points",
			3,
			[]bt.DataEventHandler{
				&bt.Bar{Close: 50},
				&bt.Bar{Close: 50},
				&bt.Bar{Close: 30},
				&bt.Bar{Close: 10},
				&bt.Bar{Close: 20},
			},
			20, nil,
		},
	}

	for _, tc := range testCases {
		sma, err := calculateSMA(tc.period, tc.data)
		if (sma != tc.expSMA) || (!reflect.DeepEqual(err, tc.expErr)) {
			t.Errorf("%v calculateSMA(%v, %v): \nexpected %v %+v, \nactual   %v %+v",
				tc.msg, tc.period, tc.data, tc.expSMA, tc.expErr, sma, err)
		}
	}
}
