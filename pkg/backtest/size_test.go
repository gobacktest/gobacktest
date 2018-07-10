package backtest

import (
	"errors"
	"reflect"
	"testing"
)

func TestSizeOrder(t *testing.T) {
	// testCases is a table for testing
	var testCases = []struct {
		msg       string // test message
		size      SizeHandler
		order     OrderEvent       // OrderEvent input
		data      DataEvent        // DataEvent input
		portfolio PortfolioHandler // the portfolio holdings
		expOrder  OrderEvent       // expected OrderEvent output
		expErr    error            // expected error output
	}{
		{"Empty SizeManager without default values:",
			&Size{},
			&Order{},
			&Bar{},
			&Portfolio{},
			&Order{},
			errors.New("cannot size order: no defaulSize or defaultValue set"),
		},
		{"buy order:",
			&Size{DefaultSize: 100, DefaultValue: 1000},
			&Order{direction: BOT},
			&Bar{},
			&Portfolio{},
			&Order{qty: 100, direction: BOT},
			nil,
		},
		{"sell order:",
			&Size{DefaultSize: 100, DefaultValue: 1000},
			&Order{direction: SLD},
			&Bar{},
			&Portfolio{},
			&Order{qty: 100, direction: SLD},
			nil,
		},
		{"exit order but no position in portfolio:",
			&Size{DefaultSize: 100, DefaultValue: 1000},
			&Order{direction: EXT},
			&Bar{},
			&Portfolio{},
			&Order{direction: EXT},
			errors.New("cannot exit order: no position to symbol in portfolio,"),
		},
		{"exit order with long position in portfolio:",
			&Size{DefaultSize: 100, DefaultValue: 1000},
			&Order{
				Event:     Event{symbol: "TEST.DE"},
				direction: EXT},
			&Bar{},
			&Portfolio{
				holdings: map[string]Position{
					"TEST.DE": {qty: 15},
				},
			},
			&Order{
				Event:     Event{symbol: "TEST.DE"},
				direction: SLD,
				qty:       15,
			},
			nil,
		},
		{"exit order with short position in portfolio:",
			&Size{DefaultSize: 100, DefaultValue: 1000},
			&Order{
				Event:     Event{symbol: "TEST.DE"},
				direction: EXT},
			&Bar{},
			&Portfolio{
				holdings: map[string]Position{
					"TEST.DE": {qty: -12},
				},
			},
			&Order{
				Event:     Event{symbol: "TEST.DE"},
				direction: BOT,
				qty:       12,
			},
			nil,
		},
	}

	for _, tc := range testCases {
		order, err := tc.size.SizeOrder(tc.order, tc.data, tc.portfolio)
		if !reflect.DeepEqual(order, tc.expOrder) || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("%v SizeOrder(%v %v %v): \nexpected %+v %v, \nactual   %+v %v",
				tc.msg, tc.order, tc.data, tc.portfolio, tc.expOrder, tc.expErr, order, err)
		}
	}
}

func TestSetDefaultSize(t *testing.T) {
	// testCases is a table for testing
	var testCases = []struct {
		msg    string // test message
		size   Size
		price  float64
		expQty int64 // expected error output
	}{
		{"Empty SizeManager without default values:",
			Size{},
			10,
			0,
		},
		{"price is higher than defaultValue:",
			Size{DefaultSize: 100, DefaultValue: 1000},
			15,
			66,
		},
		{"sprice is lower than defaultValue:",
			Size{DefaultSize: 100, DefaultValue: 1000},
			8,
			100,
		},
	}

	for _, tc := range testCases {
		qty := tc.size.setDefaultSize(tc.price)
		if qty != tc.expQty {
			t.Errorf("%v setDefaultSize(%v): \nexpected %v, \nactual   %v",
				tc.msg, tc.price, tc.expQty, qty)
		}
	}
}
