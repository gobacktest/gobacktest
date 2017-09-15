package internal

import (
	"reflect"
	"testing"
	"time"
)

// set the example time string in formar yyyy-mm-dd
var exampleTimeString = "2017-06-01"
var exampleTime, _ = time.Parse("2006-01-02", exampleTimeString)

// barDataTests is a table for testing parsing bar data into a BarEvent
var barDataTests = []struct {
	line     map[string]string // start input
	symbol   string            // symbol input
	expEvent BarEvent          // expected bar return
	expErr   error             // expected error output
}{
	{map[string]string{
		"Date":      exampleTimeString,
		"Open":      "10",
		"High":      "10",
		"Low":       "10",
		"Close":     "10",
		"Adj Close": "10",
		"Volume":    "100"},
		"TEST.DE",
		bar{
			event:         event{timestamp: exampleTime, symbol: "TEST.DE"},
			openPrice:     float64(10),
			highPrice:     float64(10),
			lowPrice:      float64(10),
			closePrice:    float64(10),
			adjClosePrice: float64(10),
			volume:        100},
		nil},
	{map[string]string{
		"Date":      exampleTimeString,
		"Open":      "null", // field in csv ist marked null, means no data
		"High":      "null",
		"Low":       "null",
		"Close":     "null",
		"Adj Close": "null",
		"Volume":    "null"},
		"TEST.DE",
		bar{
			event: event{timestamp: exampleTime, symbol: "TEST.DE"},
		}, // other values are nil
		nil},
}

func TestCreateBarEventFromLine(t *testing.T) {
	for _, tt := range barDataTests {
		event, err := createBarEventFromLine(tt.line, tt.symbol)
		if !reflect.DeepEqual(event, tt.expEvent) || (reflect.TypeOf(err) != reflect.TypeOf(tt.expErr)) {
			t.Logf("createBarEventFromLine(%v, %v): \nexpected %#v %v, \nactual   %#v %v",
				tt.line, tt.symbol, tt.expEvent, tt.expErr, event, err)
		}
	}
}

// barDataTests is a table for testing parsing bar data into a BarEvent
var barDataBenchLine = map[string]string{
	"Date":      exampleTimeString,
	"Open":      "10.50",
	"High":      "15.00",
	"Low":       "9.00",
	"Close":     "12.00",
	"Adj Close": "12.00",
	"Volume":    "100"}
var barDataBenchSymbol = "BAS.DE"

func BenchmarkCreateBarEventFromLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createBarEventFromLine(barDataBenchLine, barDataBenchSymbol)
	}

}
