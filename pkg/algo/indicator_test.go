package algo

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

func TestSMAIntegration(t *testing.T) {
	// define dates
	var dates = []string{
		"2018-07-01",
		"2018-07-02",
		"2018-07-03",
		"2018-07-04",
		"2018-07-05",
		"2018-07-06",
		"2018-07-07",
		"2018-07-08",
		"2018-07-09",
		"2018-07-10",
	}
	// set up mock Data Events
	mockdata := []bt.DataEventHandler{}
	for i, d := range dates {
		time, _ := time.Parse("2006-01-02", d)
		bar := bt.Bar{
			Event: bt.Event{
				Symbol:    "Date" + strconv.Itoa(i),
				Timestamp: time,
			},
			Close: float64(i),
		}
		mockdata = append(mockdata, bar)
	}

	// set up data handler
	data := &bt.Data{}
	data.SetStream(mockdata)
	event, _ := data.Next()

	// set up strategy
	strategy := &bt.Strategy{}
	strategy.SetData(data)
	strategy.SetEvent(event)

	// create Algo
	algo := NewSMA(5)

	// first run, no data in history, sma can not be calculated
	ok, err := algo.Run(strategy)
	if ok || (err == nil) {
		t.Errorf("first run, no data in history: \nexpected %v %#v, \nactual   %v %#v",
			false, fmt.Errorf("invalid value length for indicator sma"),
			ok, err,
		)
	}

}
