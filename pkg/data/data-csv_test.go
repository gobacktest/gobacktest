package data

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/dirkolbrich/gobacktest/pkg/backtest"
)

func TestCreateBarEventFromLine(t *testing.T) {
	// set the example time string in format yyyy-mm-dd
	var exampleTime, _ = time.Parse("2006-01-02", "2017-06-01")

	// testCases is a table for testing parsing bar data into a BarEvent
	var testCases = []struct {
		line     map[string]string // start input
		symbol   string            // symbol input
		expEvent backtest.BarEvent // expected bar return
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
			&backtest.Bar{
				Event:     backtest.Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				DataEvent: backtest.DataEvent{Metrics: make(map[string]float64)},
				Open:      float64(10),
				High:      float64(10),
				Low:       float64(10),
				Close:     float64(10),
				AdjClose:  float64(10),
				Volume:    100,
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
			(*backtest.Bar)(nil), //
			//errors.New("strconv.ParseFloat: parsing \"null\": invalid syntax"),
			&strconv.NumError{},
		},
	}

	for _, tc := range testCases {
		event, err := createBarEventFromLine(tc.line, tc.symbol)
		// if !reflect.DeepEqual(event, tc.expEvent) {
		if !reflect.DeepEqual(event, tc.expEvent) || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("createBarEventFromLine(%v, %v): \nexpected %T %p %#v %v\nactual   %T %p %#v %v",
				tc.line, tc.symbol, tc.expEvent, tc.expEvent, tc.expEvent, tc.expErr, event, event, event, err)
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
