package internal

import (
	"reflect"
	"testing"
)

// set ExecutionHandler with symbol
var e = &Exchange{Symbol: "TEST", ExchangeFee: 1.00}

// orderEventTests is a table for testing parsing bar data into a BarEvent
var orderEventTests = []struct {
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
			net:         110.90,
		},
		nil,
	},
}

func TestExecuteOrder(t *testing.T) {
	for _, tt := range orderEventTests {
		fill, err := e.ExecuteOrder(tt.order, tt.data)
		if !reflect.DeepEqual(fill, tt.expFill) || (reflect.TypeOf(err) != reflect.TypeOf(tt.expErr)) {
			t.Fatalf("ExecuteOrder(%v): \nexpected %+v %v, \nactual   %+v %v",
				tt.order, tt.expFill, tt.expErr, fill, err)
		}
	}
}
