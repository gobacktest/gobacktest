package internal

import (
	"reflect"
	"testing"
	"time"
)

// set ExecutionHandler with symbol
var e := &Exchange(Symbol: "XETRA")

// set the example time string in format yyyy-mm-dd
var exampleTimeString = "2017-06-01"
var exampleTime, _ = time.Parse("2006-01-02", exampleTimeString)

// orderEventTests is a table for testing parsing bar data into a BarEvent
var orderEventTests = []struct {
	order   OrderEvent // start input
	expFill FillEvent  // expected FillEvent return
	expErr  error      // expected error output
}{
	{OrderEvent{
		Event:     Event{timestamp: exampleTime, symbol: "BAS.DE"},
		Direction: "buy",  // buy or sell
		Qty:       10},
	FillEvent{
		Event:       Event{timestamp: exampleTime, symbol: "BAS.DE"},
		Exchange:    "XETRA",
		Direction:   "BOT", // BOT for buy or SLD for sell
		Qty:         10,
		Price:       10
		Commission:  10
		ExchangeFee: 1.2
		Cost:        11.2,
		Net:         111.2},
	nil},
}

func TestExecuteOrder(t *testing.T) {
	for _, tt := range orderEventTests {
		fill, ok := e.ExecuteOrder(tt.order)
		if (fill != tt.expFill) || (reflect.TypeOf(err) != reflect.TypeOf(tt.expErr)) {
			t.Errorf("ExecuteOrder(%v): expected %+v %v, actual %+v %v",
				tt.order, tt.expFill, tt.expErr, fill, ok)
		}
	}
}
