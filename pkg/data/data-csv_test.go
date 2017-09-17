package data

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateBarEventFromLine(t *testing.T) {
	// set the example time string in format yyyy-mm-dd
	var exampleTime, _ = time.Parse("2006-01-02", "2017-06-01")

	// testCases is a table for testing parsing bar data into a BarEvent
	var testCases = []struct {
		line     map[string]string // start input
		symbol   string            // symbol input
		expEvent BarEvent          // expected bar return
		expErr   error             // expected error output
	}{
		{
			map[string]string{
			"Date":      "2017-06-01",
			"Open":      "10",
			"High":      "10",
			"Low":       "10",
			"Close":     "10",
			"Adj Close": "10",
			"Volume":    "100",
			},
			"TEST.DE",
			bar{
				event:         event{timestamp: exampleTime, symbol: "TEST.DE"},
				openPrice:     float64(10),
				highPrice:     float64(10),
				lowPrice:      float64(10),
				closePrice:    float64(10),
				adjClosePrice: float64(10),
				volume:        100,
			},
			nil},
		{
			map[string]string{
			"Date":      "2017-06-01",
			"Open":      "null", // field in csv ist marked null, means no data
			"High":      "null",
			"Low":       "null",
			"Close":     "null",
			"Adj Close": "null",
			"Volume":    "null",
			},
			"TEST.DE",
			bar{
				event: event{timestamp: exampleTime, symbol: "TEST.DE"},
			}, // other values are nil
			nil},
	}

	for _, tc := range testCases {
		event, err := createBarEventFromLine(tc.line, tc.symbol)
		if !reflect.DeepEqual(event, tc.expEvent) || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Logf("createBarEventFromLine(%v, %v): \nexpected %#v %v, \nactual   %#v %v",
				tc.line, tc.symbol, tc.expEvent, tc.expErr, event, err)
		}
	}
}

func BenchmarkCreateBarEventFromLine(b *testing.B) {
	// barDataTests is a table for testing parsing bar data into a BarEvent
	var barDataBenchLine = map[string]string{
		"Date":      "2017-06-01",
		"Open":      "10.50",
		"High":      "15.00",
		"Low":       "9.00",
		"Close":     "12.00",
		"Adj Close": "12.00",
		"Volume":    "100"}
	var barDataBenchSymbol = "BAS.DE"

	for i := 0; i < b.N; i++ {
		createBarEventFromLine(barDataBenchLine, barDataBenchSymbol)
	}

}
