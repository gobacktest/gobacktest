package ta

import (
	"errors"
	"reflect"
	"testing"
)

func TestMean(t *testing.T) {
	var testCases = []struct {
		msg     string
		values  []float64
		expMean float64
	}{
		{"test zero data point",
			[]float64{},
			0,
		},
		{"test simple values",
			[]float64{2, 3},
			2.5,
		},
		{"test 10 values",
			[]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			4.5,
		},
		{"test float values",
			[]float64{2, 3, 5},
			3.3333333333333335,
		},
	}

	for _, tc := range testCases {
		mean := Mean(tc.values)
		if mean != tc.expMean {
			t.Errorf("%v Mean(%v): \nexpected %v, \nactual   %v",
				tc.msg, tc.values, tc.expMean, mean)
		}
	}
}

func TestSMA(t *testing.T) {
	var testCases = []struct {
		msg    string
		values []float64
		period int
		expSMA []float64
		expErr error
	}{
		{"test zero values",
			[]float64{},
			0,
			nil,
			errors.New("no values given"),
		},
		{"test length of values less than period ",
			[]float64{0, 1, 2},
			5,
			nil,
			errors.New("invalid length of values, given 3, needs 5"),
		},
		{"test simple values",
			[]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			5,
			[]float64{2, 3, 4, 5, 6, 7},
			nil,
		},
		{"test complex values",
			[]float64{0, 0.5, 1, 1.5, 2},
			3,
			[]float64{0.5, 1, 1.5},
			nil,
		},
	}

	for _, tc := range testCases {
		sma, err := SMA(tc.values, tc.period)
		if (!reflect.DeepEqual(sma, tc.expSMA)) || (!reflect.DeepEqual(err, tc.expErr)) {
			t.Errorf("%v SMA(%v, %v): \nexpected %#v %v, \nactual   %#v %v",
				tc.msg, tc.values, tc.period, tc.expSMA, tc.expErr, sma, err)
		}
	}
}

func TestEMA(t *testing.T) {
	var testCases = []struct {
		msg    string
		values []float64
		period int
		expEMA []float64
		expErr error
	}{
		{"test zero values",
			[]float64{},
			0,
			nil,
			errors.New("no values given"),
		},
		{"test length of values less than period ",
			[]float64{0, 1, 2},
			5,
			nil,
			errors.New("invalid length of values, given 3, needs 5"),
		},
		{"test simple values",
			[]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			5,
			[]float64{2, 3, 4, 5, 6, 7},
			nil,
		},
		{"test complex values",
			[]float64{26.0, 54.0, 8.0, 77.0, 61.0, 39.0, 44.0, 91.0, 98.0, 17.0},
			5,
			[]float64{45.2, 43.13333333333333, 43.422222222222224, 59.28148148148148, 72.18765432098765, 53.7917695473251},
			nil,
		},
	}

	for _, tc := range testCases {
		ema, err := EMA(tc.values, tc.period)
		if (!reflect.DeepEqual(ema, tc.expEMA)) || (!reflect.DeepEqual(err, tc.expErr)) {
			t.Errorf("%v EMA(%v, %v): \nexpected %#v %v, \nactual   %#v %v",
				tc.msg, tc.values, tc.period, tc.expEMA, tc.expErr, ema, err)
		}
	}
}
