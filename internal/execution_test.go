package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestExecuteOrder(t *testing.T) {
	// set the example time string in format yyyy-mm-dd
	var exampleTime, _ = time.Parse("2006-01-02", "2017-06-01")

	// set ExecutionHandler with symbol
	var e = &Exchange{Symbol: "TEST", ExchangeFee: 1.00}

	// orderEventTests is a table for testing parsing bar data into a BarEvent
	var testCases = []struct {
		order   OrderEvent  // OrderEvent input
		data    DataHandler // DataEvent input
		expFill FillEvent   // expected FillEvent return
		expErr  error       // expected error output
	}{
		{
			&order{
				event:     event{timestamp: exampleTime, symbol: "TEST.DE"},
				direction: "buy", // buy or sell
				qty:       10,
			},
			&Data{
				latest: map[string]DataEvent{
					"TEST.DE": &bar{closePrice: 10},
				},
			},
			&fill{
				event:       event{timestamp: exampleTime, symbol: "TEST.DE"},
				exchange:    "TEST",
				direction:   "BOT", // BOT for buy or SLD for sell
				qty:         10,
				price:       10,
				commission:  9.90,
				exchangeFee: 1,
				cost:        10.90,
			},
			nil,
		},
		{
			&order{
				event:     event{timestamp: exampleTime, symbol: "TEST.DE"},
				direction: "sell", // buy or sell
				qty:       10,
			},
			&Data{
				latest: map[string]DataEvent{
					"TEST.DE": &bar{closePrice: 10},
				},
			},
			&fill{
				event:       event{timestamp: exampleTime, symbol: "TEST.DE"},
				exchange:    "TEST",
				direction:   "SLD", // BOT for buy or SLD for sell
				qty:         10,
				price:       10,
				commission:  9.90,
				exchangeFee: 1,
				cost:        10.90,
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
