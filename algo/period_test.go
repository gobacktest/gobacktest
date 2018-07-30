package algo

import (
	// "fmt"
	"reflect"
	"testing"
	"time"

	gbt "github.com/dirkolbrich/gobacktest"
)

func TestAlgoRunOnce(t *testing.T) {
	algo := RunOnce()

	ok, err := algo.Run(&gbt.Strategy{})
	if (ok == false) || (err != nil) {
		t.Errorf("first RunOnce(): \nexpected %v %#v, \nactual   %v %#v", true, nil, ok, err)
	}

	ok, err = algo.Run(&gbt.Strategy{})
	if (ok == true) || (err != nil) {
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
	// set up mock time
	times := testHelperTimeMap([]string{
		"2017-12-31",
		"2018-01-01",
		"2018-06-30",
		"2018-07-01",
	})

	var testCases = []struct {
		msg       string
		now       time.Time
		toCompare time.Time
		expOk     bool
		expErr    error
	}{
		{"test day change",
			times["2018-06-30"], times["2018-07-01"],
			true, nil,
		},
		{"test same day ",
			times["2018-06-30"], times["2018-06-30"],
			false, nil,
		},
		{"test year change",
			times["2017-12-31"], times["2018-01-01"],
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
	// set up mock Data Events
	mockdata := testHelperMockData([]string{
		"2018-06-30",
		"2018-07-01",
		"2018-07-02",
	})

	// set up data handler
	data := &gbt.Data{}
	data.SetStream(mockdata)
	event, _ := data.Next()

	// set up strategy
	strategy := &gbt.Strategy{}
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
	// set up mock time
	times := testHelperTimeMap([]string{
		"2016-12-31",
		"2017-01-01",
		"2017-12-31",
		"2018-01-01",
		"2018-06-28",
		"2018-06-29",
		"2018-06-30",
		"2018-07-01",
		"2018-07-02",
	})

	var testCases = []struct {
		msg       string
		now       time.Time
		toCompare time.Time
		expOk     bool
		expErr    error
	}{
		{"test week change",
			times["2018-07-01"], times["2018-07-02"],
			true, nil,
		},
		{"test same week on weekend",
			times["2018-06-30"], times["2018-07-01"],
			false, nil,
		},
		{"test same week during week",
			times["2018-06-28"], times["2018-06-29"],
			false, nil,
		},
		{"test year change, same week",
			times["2016-12-31"], times["2017-01-01"],
			false, nil,
		},
		{"test year change with week change",
			times["2017-12-31"], times["2018-01-01"],
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
	// set up mock time
	times := testHelperTimeMap([]string{
		"2017-12-31",
		"2018-01-01",
		"2018-01-02",
		"2018-02-01",
	})

	var testCases = []struct {
		msg       string
		now       time.Time
		toCompare time.Time
		expOk     bool
		expErr    error
	}{
		{"test month change",
			times["2018-01-02"], times["2018-02-01"],
			true, nil,
		},
		{"test same month ",
			times["2018-01-01"], times["2018-01-02"],
			false, nil,
		},
		{"test year change",
			times["2017-12-31"], times["2018-01-01"],
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
	// set up mock Data Events
	mockdata := testHelperMockData([]string{
		"2017-12-31",
		"2018-01-01",
		"2018-01-02",
		"2018-02-01",
	})

	// set up data handler
	data := &gbt.Data{}
	data.SetStream(mockdata)
	event, _ := data.Next()

	// set up strategy
	strategy := &gbt.Strategy{}
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
	// set up mock time
	times := testHelperTimeMap([]string{
		"2017-12-31",
		"2018-01-01",
		"2018-01-02",
		"2018-02-01",
	})

	var testCases = []struct {
		msg       string
		now       time.Time
		toCompare time.Time
		expOk     bool
		expErr    error
	}{
		{msg: "test year change",
			now: times["2017-12-31"], toCompare: times["2018-01-01"],
			expOk: true, expErr: nil,
		},
		{msg: "test same year, different day in month",
			now: times["2018-01-01"], toCompare: times["2018-01-02"],
			expOk: false, expErr: nil,
		},
		{msg: "test same year, different month",
			now: times["2018-01-01"], toCompare: times["2018-02-01"],
			expOk: false, expErr: nil,
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

func TestAlgoRunYearlyImplementation(t *testing.T) {
	// set up mock Data Events
	mockdata := testHelperMockData([]string{
		"2017-12-31",
		"2018-01-01",
		"2018-01-02",
		"2018-02-01",
	})

	// set up data handler
	data := &gbt.Data{}
	data.SetStream(mockdata)
	event, _ := data.Next()

	// set up strategy
	strategy := &gbt.Strategy{}
	strategy.SetData(data)
	strategy.SetEvent(event)

	// create Algo
	algo := RunYearly()

	// first data, no data in history
	ok, err := algo.Run(strategy)
	if (ok == true) || (err != nil) {
		t.Errorf("first data, no data in history: \nexpected %v %#v, \nactual   %v %#v", false, nil, ok, err)
	}

	// second data, one more data in history with year change
	// pull next event
	event, _ = data.Next()
	strategy.SetEvent(event)
	ok, err = algo.Run(strategy)
	if (ok == false) || (err != nil) {
		t.Errorf("second data, year and month change: \nexpected %v %#v, \nactual   %v %#v", true, nil, ok, err)
	}

	// third data, one more data in history without day change
	// pull next event
	event, _ = data.Next()
	strategy.SetEvent(event)
	ok, err = algo.Run(strategy)
	if (ok == true) || (err != nil) {
		t.Errorf("third data, no year change, no month change: \nexpected %v %#v, \nactual   %v %#v", false, nil, ok, err)
	}

	event, _ = data.Next()
	strategy.SetEvent(event)
	ok, err = algo.Run(strategy)
	if (ok == true) || (err != nil) {
		t.Errorf("third data, no year but month change: \nexpected %v %#v, \nactual   %v %#v", false, nil, ok, err)
	}
}
