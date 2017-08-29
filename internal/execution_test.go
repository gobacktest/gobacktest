package internal

import (
	"reflect"
	"testing"
)

// set ExecutionHandler with symbol
var e = &Exchange{Symbol: "TEST"}

// set the example time string in format yyyy-mm-dd
//var exampleTimeString = "2017-06-01"
//var exampleTime, _ = time.Parse("2006-01-02", exampleTimeString)

// orderEventTests is a table for testing parsing bar data into a BarEvent
var orderEventTests = []struct {
	order   OrderEvent // OrderEvent input
	data    DataEvent  // DataEvent input
	expFill FillEvent  // expected FillEvent return
	expErr  error      // expected error output
}{
	{&order{
		event:     event{timestamp: exampleTime, symbol: "TEST.DE"},
		direction: "buy", // buy or sell
		qty:       10},
		&bar{},
		&fill{
			event:       event{timestamp: exampleTime, symbol: "TEST.DE"},
			exchange:    "TEST",
			direction:   "BOT", // BOT for buy or SLD for sell
			qty:         10,
			price:       10,
			commission:  10,
			exchangeFee: 1,
			cost:        11,
			net:         111},
		nil},
}

func TestExecuteOrder(t *testing.T) {
	for _, tt := range orderEventTests {
		fill, err := e.ExecuteOrder(tt.order, tt.data)
		if (fill != tt.expFill) || (reflect.TypeOf(err) != reflect.TypeOf(tt.expErr)) {
			t.Errorf("ExecuteOrder(%v): \nexpected %+v %v, \nactual   %+v %v",
				tt.order, tt.expFill, tt.expErr, fill, err)
		}
	}
}
