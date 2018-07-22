package gobacktest

import (
	"errors"
	"reflect"
	"testing"
)

func TestMetricAdd(t *testing.T) {
	var testCases = []struct {
		msg    string
		metric Metric
		key    string
		value  float64
		expM   Metric
		expErr error
	}{
		{
			msg:    "empty key:",
			metric: Metric{},
			key:    "",
			value:  123.45,
			expM:   map[string]float64{},
			expErr: errors.New("invalid key given"),
		},
		{
			msg:    "simple key value:",
			metric: Metric{},
			key:    "test",
			value:  123.45,
			expM:   map[string]float64{"test": 123.45},
			expErr: nil,
		},
		{
			msg:    "add to existing key value:",
			metric: map[string]float64{"abcd": 678.90},
			key:    "test",
			value:  123.45,
			expM: map[string]float64{
				"test": 123.45,
				"abcd": 678.90,
			},
			expErr: nil,
		},
		{
			msg:    "overwrite existing key value:",
			metric: map[string]float64{"test": 678.90},
			key:    "test",
			value:  123.45,
			expM:   map[string]float64{"test": 123.45},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		err := tc.metric.Add(tc.key, tc.value)
		if !reflect.DeepEqual(tc.metric, tc.expM) || !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("%v Add(%v, %v): \nexpected %#v %#v\nactual   %#v %#v", tc.msg, tc.key, tc.value, tc.expM, tc.expErr, tc.metric, err)
		}
	}
}

func TestMetricGet(t *testing.T) {
	var testCases = []struct {
		msg    string
		metric Metric
		key    string
		expVal float64
		expOk  bool
	}{
		{
			msg:    "empty metric:",
			metric: Metric{},
			key:    "test",
			expVal: 0,
			expOk:  false,
		},
		{
			msg:    "simple metric:",
			metric: map[string]float64{"test": 123.45},
			key:    "test",
			expVal: 123.45,
			expOk:  true,
		},
		{
			msg: "multiple metric:",
			metric: map[string]float64{
				"test": 123.45,
				"abcd": 678.90,
			},
			key:    "test",
			expVal: 123.45,
			expOk:  true,
		},
	}

	for _, tc := range testCases {
		value, ok := tc.metric.Get(tc.key)
		if (value != tc.expVal) || (ok != tc.expOk) {
			t.Errorf("%v Metric(%v): \nexpected %#v %v\nactual   %#v %v", tc.msg, tc.key, tc.expVal, tc.expOk, value, ok)
		}
	}
}

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

func TestSignalSetDirection(t *testing.T) {
	var testCases = []struct {
		msg       string
		signal    Signal
		dir       OrderDirection
		expSignal Signal
	}{
		{"simple direction:",
			Signal{},
			BOT,
			Signal{direction: BOT},
		},
	}

	for _, tc := range testCases {
		tc.signal.SetDirection(tc.dir)
		if !reflect.DeepEqual(tc.signal, tc.expSignal) {
			t.Errorf("%v SetDirection(%v): \nexpected %#v, \nactual %#v",
				tc.msg, tc.dir, tc.expSignal, tc.signal)
		}
	}
}

func TestSignalGetDirection(t *testing.T) {
	var testCases = []struct {
		msg    string
		signal Signal
		expDir OrderDirection
	}{
		{"simple direction:",
			Signal{direction: BOT},
			BOT,
		},
	}

	for _, tc := range testCases {
		dir := tc.signal.Direction()
		if dir != tc.expDir {
			t.Errorf("%v Direction(): \nexpected %#v, \nactual %#v",
				tc.msg, tc.expDir, dir)
		}
	}
}

func TestFillSetDirection(t *testing.T) {
	var testCases = []struct {
		msg     string
		fill    Fill
		dir     OrderDirection
		expFill Fill
	}{
		{"simple direction:",
			Fill{},
			BOT,
			Fill{direction: BOT},
		},
	}

	for _, tc := range testCases {
		tc.fill.SetDirection(tc.dir)
		if !reflect.DeepEqual(tc.fill, tc.expFill) {
			t.Errorf("%v SetDirection(%v): \nexpected %#v, \nactual %#v",
				tc.msg, tc.dir, tc.expFill, tc.fill)
		}
	}
}

func TestFillSetQty(t *testing.T) {
	var testCases = []struct {
		msg     string
		fill    Fill
		qty     int64
		expFill Fill
	}{
		{"simple qty:",
			Fill{},
			100,
			Fill{qty: 100},
		},
	}

	for _, tc := range testCases {
		tc.fill.SetQty(tc.qty)
		if !reflect.DeepEqual(tc.fill, tc.expFill) {
			t.Errorf("%v SetQty(%v): \nexpected %#v, \nactual %#v",
				tc.msg, tc.qty, tc.expFill, tc.fill)
		}
	}
}

func TestFillValue(t *testing.T) {
	var testCases = []struct {
		msg  string
		fill Fill
		exp  float64
	}{
		{"Empty Fill:",
			Fill{qty: 0, price: 0},
			0,
		},
		{"Standard Fill:",
			Fill{qty: 10, price: 5},
			50,
		},
	}

	for _, tc := range testCases {
		float := tc.fill.Value()
		if float != tc.exp {
			t.Errorf("%v Value(): \nexpected %#v, \nactual %#v", tc.msg, tc.exp, float)
		}
	}
}

func TestFillNetValue(t *testing.T) {
	var testCases = []struct {
		msg  string
		fill Fill
		exp  float64
	}{
		{"Empty BOT Fill:",
			Fill{direction: BOT, qty: 0, price: 0, cost: 0},
			0,
		},
		{"Standard BOT Fill:",
			Fill{direction: BOT, qty: 10, price: 5, cost: 5},
			55,
		},
		{"Empty SLD Fill:",
			Fill{direction: SLD, qty: 0, price: 0, cost: 0},
			0,
		},
		{"Standard SLD Fill:",
			Fill{direction: SLD, qty: 10, price: 5, cost: 5},
			45,
		},
	}

	for _, tc := range testCases {
		float := tc.fill.NetValue()
		if float != tc.exp {
			t.Errorf("%v NetValue(): \nexpected %#v, \nactual %#v", tc.msg, tc.exp, float)
		}
	}
}
