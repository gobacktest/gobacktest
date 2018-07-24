package gobacktest

import (
	"errors"
	"reflect"
	"testing"
)

func TestMetricAdd(t *testing.T) {
	var testCases = []struct {
		msg    string
		metric Metric
		key    string
		value  float64
		expM   Metric
		expErr error
	}{
		{
			msg:    "empty key:",
			metric: Metric{},
			key:    "",
			value:  123.45,
			expM:   map[string]float64{},
			expErr: errors.New("invalid key given"),
		},
		{
			msg:    "simple key value:",
			metric: Metric{},
			key:    "test",
			value:  123.45,
			expM:   map[string]float64{"test": 123.45},
			expErr: nil,
		},
		{
			msg:    "add to existing key value:",
			metric: map[string]float64{"abcd": 678.90},
			key:    "test",
			value:  123.45,
			expM: map[string]float64{
				"test": 123.45,
				"abcd": 678.90,
			},
			expErr: nil,
		},
		{
			msg:    "overwrite existing key value:",
			metric: map[string]float64{"test": 678.90},
			key:    "test",
			value:  123.45,
			expM:   map[string]float64{"test": 123.45},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		err := tc.metric.Add(tc.key, tc.value)
		if !reflect.DeepEqual(tc.metric, tc.expM) || !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("%v Add(%v, %v): \nexpected %#v %#v\nactual   %#v %#v", tc.msg, tc.key, tc.value, tc.expM, tc.expErr, tc.metric, err)
		}
	}
}

func TestMetricGet(t *testing.T) {
	var testCases = []struct {
		msg    string
		metric Metric
		key    string
		expVal float64
		expOk  bool
	}{
		{
			msg:    "empty metric:",
			metric: Metric{},
			key:    "test",
			expVal: 0,
			expOk:  false,
		},
		{
			msg:    "simple metric:",
			metric: map[string]float64{"test": 123.45},
			key:    "test",
			expVal: 123.45,
			expOk:  true,
		},
		{
			msg: "multiple metric:",
			metric: map[string]float64{
				"test": 123.45,
				"abcd": 678.90,
			},
			key:    "test",
			expVal: 123.45,
			expOk:  true,
		},
	}

	for _, tc := range testCases {
		value, ok := tc.metric.Get(tc.key)
		if (value != tc.expVal) || (ok != tc.expOk) {
			t.Errorf("%v Metric(%v): \nexpected %#v %v\nactual   %#v %v", tc.msg, tc.key, tc.expVal, tc.expOk, value, ok)
		}
	}
}
