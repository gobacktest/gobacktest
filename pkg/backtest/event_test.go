package backtest

import (
	"reflect"
	"testing"
)

func TestTickLatestPrice(t *testing.T) {
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
		float := tc.tick.LatestPrice()
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
