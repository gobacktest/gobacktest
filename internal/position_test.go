package internal

import (
	"reflect"
	"testing"
	"time"
)

// initialize new Position ready for use
var p = new(Position)

// set the example time string in format yyyy-mm-dd
var exampleTimeString = "2017-06-01"
var exampleTime, _ = time.Parse("2006-01-02", exampleTimeString)

// createPositionTests is a table for testing creation of a Position
var createPositionTests = []struct {
	fill   FillEvent // input
	expPos Position  // expected Position
}{
	{FillEvent{
		Event:       Event{timestamp: exampleTime, symbol: "BAS.DE"},
		Exchange:    "XETRA",
		Direction:   "BOT", // BOT for buy or SLD for sell
		Qty:         10,
		Price:       10
		Commission:  5
		ExchangeFee: 1.2
		Cost:        6.2,
		Net:         106.2},
	Position{
		timestamp:   exampleTime,
		symbol:      "BAS.DE",
		qty:         10,
		avgPrice:    10,
		value:       100,
		marketPrice: 10,
		marketValue: 100,
		commission:  5,
		exchangeFee: 1.2,
		cost:        6.2,
		netValue:         },
	},
}

func TestCreate(t *testing.T) {
	for _, tt := range createPositionTests {
		p.Create(tt.fill)
		if (p != tt.expPos) {
			t.Errorf("Create(%v): expected %+v, actual %+v", tt.fill, tt.expPos, p)
		}
	}
}
