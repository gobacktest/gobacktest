package algo

import (
	"fmt"
	"reflect"
	"testing"

	gbt "github.com/dirkolbrich/gobacktest"
)

func TestSMAIntegration(t *testing.T) {
	// set up mock Data Events
	mockdata := testHelperMockData([]string{
		"2018-07-01",
		"2018-07-02",
		"2018-07-03",
		"2018-07-04",
		"2018-07-05",
	})

	// set close price from 1 to n on mockdata
	for i, data := range mockdata {
		bar := data.(*gbt.Bar)
		bar.Close = float64(i + 1)
		mockdata[i] = bar
	}

	var testCases = []struct {
		msg       string
		mockdata  []gbt.DataEvent
		period    int
		runBefore int
		expOk     bool
		expErr    error
	}{
		{msg: "test too few data points",
			mockdata: mockdata[:1],
			period:   5,
			expOk:    false,
			expErr:   fmt.Errorf("invalid value length for indicator sma"),
		},
		{msg: "test normal run",
			mockdata:  mockdata[:3],
			period:    3,
			runBefore: 2,
			expOk:     true,
			expErr:    nil,
		},
	}

	for _, tc := range testCases {
		// set up data handler
		data := &gbt.Data{}
		data.SetStream(tc.mockdata)
		event, _ := data.Next()

		// set up strategy
		strategy := &gbt.Strategy{}
		strategy.SetData(data)
		strategy.SetEvent(event)

		// run the backtest n times to  pull data from stream and fill data.list
		for i := 0; i < tc.runBefore; i++ {
			data.Next()
		}

		// create Algo
		algo := SMA(tc.period)

		ok, err := algo.Run(strategy)
		if (ok != tc.expOk) || !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("%v: SMA(%v): \nexpected %v %#v, \nactual   %v %#v",
				tc.msg, tc.period, tc.expOk, tc.expErr, ok, err)
		}

	}

}
