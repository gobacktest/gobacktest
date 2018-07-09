package algo

import (
	// "fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

func TestAlgoRunOnce(t *testing.T) {
	algo := RunOnce()

	ok, err := algo.Run(&bt.Strategy{})
	if !ok || (err != nil) {
		t.Errorf("first RunOnce(): \nexpected %v %#v, \nactual   %v %#v", true, nil, ok, err)
	}

	ok, err = algo.Run(&bt.Strategy{})
	if ok || (err != nil) {
		t.Errorf("second RunOnce(): \nexpected %v %#v, \nactual   %v %#v", false, nil, ok, err)
	}
}

func TestRunPeriodWithOptions(t *testing.T) {
	var testCases = []struct {
		msg     string
		rp      runPeriod
		options []string
		expRp   runPeriod
	}{
		{"test zero option",
			runPeriod{},
			[]string{""},
			runPeriod{},
		},
		{"test invalid option",
			runPeriod{},
			[]string{"test"},
			runPeriod{},
		},
		{"test onFirstDate option",
			runPeriod{},
			[]string{"onFirstDate"},
			runPeriod{runOnFirstDate: true},
		},
		{"test onLastDate option",
			runPeriod{},
			[]string{"onLastDate"},
			runPeriod{runOnLastDate: true},
		},
		{"test runEndOfPeriod option",
			runPeriod{},
			[]string{"endOfPeriod"},
			runPeriod{runEndOfPeriod: true},
		},
		{"test multiple options",
			runPeriod{},
			[]string{"onLastDate", "endOfPeriod"},
			runPeriod{runOnLastDate: true, runEndOfPeriod: true},
		},
	}

	for _, tc := range testCases {
		rp := runPeriodWithOptions(tc.options...)
		if !reflect.DeepEqual(rp, tc.expRp) {
			t.Errorf("%v runPeriodWithOptions(%v): \nexpected %+v, \nactual   %+v",
				tc.msg, tc.options, tc.expRp, rp)
		}
	}
}

func TestAlgoRunDailyCompare(t *testing.T) {
	var dates = []string{
		"2017-12-31",
		"2018-01-01",
		"2018-06-30",
		"2018-07-01",
	}
	// set up mock time
	times := []time.Time{}
	for _, d := range dates {
		time, _ := time.Parse("2006-01-02", d)
		times = append(times, time)
	}

	var testCases = []struct {
		msg       string
		now       time.Time
		toCompare time.Time
		expOk     bool
		expErr    error
	}{
		{"test day change",
			times[2], times[3],
			true, nil,
		},
		{"test same day ",
			times[2], times[2],
			false, nil,
		},
		{"test year change",
			times[0], times[1],
			true, nil,
		},
	}

	algo := runDaily{}
	for _, tc := range testCases {
		ok, err := algo.CompareDates(tc.now, tc.toCompare)
		if (ok != tc.expOk) || (!reflect.DeepEqual(err, tc.expErr)) {
			t.Errorf("%v compareDatesDaily(%v, %v): \nexpected %v %+v, \nactual   %v %+v",
				tc.msg, tc.now, tc.toCompare, tc.expOk, tc.expErr, ok, err)
		}
	}
}

func TestAlgoRunDailyImplementation(t *testing.T) {
	// define dates to test
	var dates = []string{
		"2018-06-30",
		"2018-07-01",
		"2018-07-02",
	}
	// set up mock Data Events
	mockdata := []bt.DataEventHandler{}
	for i, d := range dates {
		time, _ := time.Parse("2006-01-02", d)
		symbol := "Date" + strconv.Itoa(i)

		event := &bt.Event{}
		event.SetSymbol(symbol)
		event.SetTime(time)

		bar := &bt.Bar{
			Event: *event,
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
	algo := RunDaily()

	// first data, no data in history
	ok, err := algo.Run(strategy)
	if ok || (err != nil) {
		t.Errorf("first run, no data in history: \nexpected %v %#v, \nactual   %v %#v", false, nil, ok, err)
	}

	// second data, one more data in history with day change
	// pull next event
	event, _ = data.Next()
	strategy.SetEvent(event)
	ok, err = algo.Run(strategy)
	if !ok || (err != nil) {
		t.Errorf("second run, day change: \nexpected %v %#v, \nactual   %v %#v", true, nil, ok, err)
	}

	// third data, one more data in history without day change
	// pull next event
	event, _ = data.Next()
	strategy.SetEvent(event)
	ok, err = algo.Run(strategy)
	if !ok || (err != nil) {
		t.Errorf("third run, no day change: \nexpected %v %#v, \nactual   %v %#v", false, nil, ok, err)
	}
}

func TestAlgoRunWeeklyCompare(t *testing.T) {
	var dates = []string{
		"2016-12-31",
		"2017-01-01",
		"2017-12-31",
		"2018-01-01",
		"2018-06-28",
		"2018-06-29",
		"2018-06-30",
		"2018-07-01",
		"2018-07-02",
	}
	// set up mock time
	times := []time.Time{}
	for _, d := range dates {
		time, _ := time.Parse("2006-01-02", d)
		times = append(times, time)
	}

	var testCases = []struct {
		msg       string
		now       time.Time
		toCompare time.Time
		expOk     bool
		expErr    error
	}{
		{"test week change",
			times[7], times[8],
			true, nil,
		},
		{"test same week on weekend",
			times[6], times[7],
			false, nil,
		},
		{"test same week during week",
			times[4], times[5],
			false, nil,
		},
		{"test year change, same week",
			times[0], times[1],
			false, nil,
		},
		{"test year change with week change",
			times[2], times[3],
			true, nil,
		},
	}

	algo := runWeekly{}
	for _, tc := range testCases {
		ok, err := algo.CompareDates(tc.now, tc.toCompare)
		if (ok != tc.expOk) || (!reflect.DeepEqual(err, tc.expErr)) {
			t.Errorf("%v compareDatesWeekly(%v, %v): \nexpected %v %+v, \nactual   %v %+v",
				tc.msg, tc.now, tc.toCompare, tc.expOk, tc.expErr, ok, err)
		}
	}
}

func TestAlgoRunMonthlyCompare(t *testing.T) {
	var dates = []string{
		"2017-12-31",
		"2018-01-01",
		"2018-01-02",
		"2018-02-01",
	}
	// set up mock time
	times := []time.Time{}
	for _, d := range dates {
		time, _ := time.Parse("2006-01-02", d)
		times = append(times, time)
	}

	var testCases = []struct {
		msg       string
		now       time.Time
		toCompare time.Time
		expOk     bool
		expErr    error
	}{
		{"test month change",
			times[2], times[3],
			true, nil,
		},
		{"test same month ",
			times[1], times[2],
			false, nil,
		},
		{"test year change",
			times[0], times[1],
			true, nil,
		},
	}

	algo := runMonthly{}
	for _, tc := range testCases {
		ok, err := algo.CompareDates(tc.now, tc.toCompare)
		if (ok != tc.expOk) || (!reflect.DeepEqual(err, tc.expErr)) {
			t.Errorf("%v compareDatesMonthly(%v, %v): \nexpected %v %+v, \nactual   %v %+v",
				tc.msg, tc.now, tc.toCompare, tc.expOk, tc.expErr, ok, err)
		}
	}
}

func TestAlgoRunMonthlyImplementation(t *testing.T) {
	// define dates to test
	var dates = []string{
		"2017-12-31",
		"2018-01-01",
		"2018-01-02",
		"2018-02-01",
	}
	// set up mock Data Events
	mockdata := []bt.DataEventHandler{}
	for i, d := range dates {
		time, _ := time.Parse("2006-01-02", d)
		symbol := "Date" + strconv.Itoa(i)

		event := &bt.Event{}
		event.SetSymbol(symbol)
		event.SetTime(time)

		bar := &bt.Bar{
			Event: *event,
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
	algo := RunMonthly()

	// first data, no data in history
	ok, err := algo.Run(strategy)
	if ok || (err != nil) {
		t.Errorf("first run, no data in history: \nexpected %v %#v, \nactual   %v %#v", false, nil, ok, err)
	}

	// second data, one more data in history with year change
	// pull next event
	event, _ = data.Next()
	strategy.SetEvent(event)
	ok, err = algo.Run(strategy)
	if !ok || (err != nil) {
		t.Errorf("second run, year and month change: \nexpected %v %#v, \nactual   %v %#v", true, nil, ok, err)
	}

	// third data, one more data in history without month change
	// pull next event
	event, _ = data.Next()
	strategy.SetEvent(event)
	ok, err = algo.Run(strategy)
	if ok || (err != nil) {
		t.Errorf("third run, no month changy: \nexpected %v %#v, \nactual   %v %#v", false, nil, ok, err)
	}

	// fourth data, one more data in history without day change
	// pull next event
	event, _ = data.Next()
	strategy.SetEvent(event)
	ok, err = algo.Run(strategy)
	if !ok || (err != nil) {
		t.Errorf("fourth run, month change: \nexpected %v %#v, \nactual   %v %#v", true, nil, ok, err)
	}
}

func TestAlgoRunYearlyCompare(t *testing.T) {
	var dates = []string{
		"2017-12-31",
		"2018-01-01",
		"2018-01-02",
	}
	// set up mock time
	times := []time.Time{}
	for _, d := range dates {
		time, _ := time.Parse("2006-01-02", d)
		times = append(times, time)
	}

	var testCases = []struct {
		msg       string
		now       time.Time
		toCompare time.Time
		expOk     bool
		expErr    error
	}{
		{"test year change",
			times[0], times[1],
			true, nil,
		},
		{"test same year ",
			times[1], times[2],
			false, nil,
		},
	}

	algo := runYearly{}
	for _, tc := range testCases {
		ok, err := algo.CompareDates(tc.now, tc.toCompare)
		if (ok != tc.expOk) || (!reflect.DeepEqual(err, tc.expErr)) {
			t.Errorf("%v CompareDatesYearly(%v, %v): \nexpected %v %+v, \nactual   %v %+v",
				tc.msg, tc.now, tc.toCompare, tc.expOk, tc.expErr, ok, err)
		}
	}
}
