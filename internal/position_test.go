package internal

import (
	"reflect"
	"testing"
)

// initialize new Position ready for use
var p = new(position)

// set the example time string in format yyyy-mm-dd
// var exampleTimeString = "2017-06-01"
// var exampleTime, _ = time.Parse("2006-01-02", exampleTimeString)

// createPositionTests is a table for testing creation of a Position
var createPositionTests = []struct {
	fill   FillEvent // input
	expPos *position // expected Position
}{
	{&fill{
		event:       event{timestamp: exampleTime, symbol: "TEST.DE"},
		exchange:    "TEST",
		direction:   "BOT", // BOT for buy or SLD for sell
		qty:         10,
		price:       10,
		commission:  5,
		exchangeFee: 1,
		cost:        6,
		net:         94},
		&position{
			timestamp:   exampleTime,
			symbol:      "TEST.DE",
			qty:         10,
			avgPrice:    10,
			value:       100,
			marketPrice: 10,
			marketValue: 100,
			commission:  5,
			exchangeFee: 1,
			cost:        6,
			netValue:    94},
	},
}

func TestCreate(t *testing.T) {
	for _, tt := range createPositionTests {
		p.Create(tt.fill)
		if !reflect.DeepEqual(p, tt.expPos) {
			t.Fatalf("Create(%v): \nexpected %p %#v, \nactual   %p %#v", tt.fill, tt.expPos, tt.expPos, p, p)
		}
	}
}
