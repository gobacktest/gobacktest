package backtest

import (
	"reflect"
	"testing"
	"time"
)

func TestExecuteOrder(t *testing.T) {
	// set the example time string in format yyyy-mm-dd
	var exampleTime, _ = time.Parse("2006-01-02", "2017-06-01")

	// set ExecutionHandler with symbol
	var e = &Exchange{
		Symbol:      "TEST",
		Commission:  &FixedCommission{Commission: 0},
		ExchangeFee: &FixedExchangeFee{ExchangeFee: 1.0},
	}

	// orderEventTests is a table for testing parsing bar data into a BarEvent
	var testCases = []struct {
		order   OrderEvent  // OrderEvent input
		data    DataHandler // DataEvent input
		expFill FillEvent   // expected FillEvent return
		expErr  error       // expected error output
	}{
		{
			&Order{
				Event:     Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Direction: "buy", // buy or sell
				Qty:       10,
			},
			&Data{
				latest: map[string]DataEventHandler{
					"TEST.DE": &Bar{Close: 10},
				},
			},
			&Fill{
				Event:       Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange:    "TEST",
				Direction:   "BOT", // BOT for buy or SLD for sell
				Qty:         10,
				Price:       10,
				Commission:  0,
				ExchangeFee: 1,
				Cost:        1,
			},
			nil,
		},
		{
			&Order{
				Event:     Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Direction: "sell", // buy or sell
				Qty:       10,
			},
			&Data{
				latest: map[string]DataEventHandler{
					"TEST.DE": &Bar{Close: 10},
				},
			},
			&Fill{
				Event:       Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange:    "TEST",
				Direction:   "SLD", // BOT for buy or SLD for sell
				Qty:         10,
				Price:       10,
				Commission:  0,
				ExchangeFee: 1,
				Cost:        1,
			},
			nil,
		},
	}

	for _, tc := range testCases {
		fill, err := e.ExecuteOrder(tc.order, tc.data)
		if !reflect.DeepEqual(fill, tc.expFill) || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("ExecuteOrder(%v): \nexpected %+v %v, \nactual   %+v %v",
				tc.order, tc.expFill, tc.expErr, fill, err)
		}
	}
}
