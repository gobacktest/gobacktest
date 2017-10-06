package backtest

import (
	"reflect"
	"testing"
)

func TestDataReset(t *testing.T) {
	var testCases = []struct {
		msg     string
		data    Data
		expData Data
	}{
		{"test with empty data stream",
			Data{
				latest: map[string]DataEventHandler{
					"TEST.DE": Bar{Close: 100},
				},
				list: map[string][]DataEventHandler{
					"TEST.DE": {
						Bar{Close: 100},
						Bar{Close: 110},
						Bar{Close: 95},
					},
				},
				stream: []DataEventHandler{},
				streamHistory: []DataEventHandler{
					Bar{Close: 100},
					Bar{Close: 110},
					Bar{Close: 95},
				},
			},
			Data{
				stream: []DataEventHandler{
					Bar{Close: 100},
					Bar{Close: 110},
					Bar{Close: 95},
				},
			},
		},
		{"test with empty data",
			Data{},
			Data{},
		},
	}

	for _, tc := range testCases {
		tc.data.Reset()
		if !reflect.DeepEqual(tc.data, tc.expData) {
			t.Errorf("%v Reset(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expData, tc.data)
		}
	}
}

func TestUpdateLatest(t *testing.T) {
	var testCases = []struct {
		msg     string
		data    Data
		event   DataEventHandler
		expData Data
	}{
		{"test update filled latest",
			Data{
				latest: map[string]DataEventHandler{
					"TEST.DE": Bar{
						Event: Event{Symbol: "TEST.DE"},
						Close: 80,
					},
				},
			},
			Bar{
				Event: Event{Symbol: "TEST.DE"},
				Close: 100,
			},
			Data{
				latest: map[string]DataEventHandler{
					"TEST.DE": Bar{
						Event: Event{Symbol: "TEST.DE"},
						Close: 100,
					},
				},
			},
		},
		{"test update empty latest",
			Data{},
			Bar{
				Event: Event{Symbol: "TEST.DE"},
				Close: 100,
			},
			Data{
				latest: map[string]DataEventHandler{
					"TEST.DE": Bar{
						Event: Event{Symbol: "TEST.DE"},
						Close: 100,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc.data.updateLatest(tc.event)
		if !reflect.DeepEqual(tc.data, tc.expData) {
			t.Errorf("%v updateLatest(%v): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.event, tc.expData, tc.data)
		}
	}
}

func TestUpdateList(t *testing.T) {
	var testCases = []struct {
		msg     string
		data    Data
		event   DataEventHandler
		expData Data
	}{
		{"test update filled list",
			Data{
				list: map[string][]DataEventHandler{
					"TEST.DE": {
						Bar{
							Event: Event{Symbol: "TEST.DE"},
							Close: 90,
						},
					},
				},
			},
			Bar{
				Event: Event{Symbol: "TEST.DE"},
				Close: 100,
			},
			Data{
				list: map[string][]DataEventHandler{
					"TEST.DE": {
						Bar{
							Event: Event{Symbol: "TEST.DE"},
							Close: 90,
						},
						Bar{
							Event: Event{Symbol: "TEST.DE"},
							Close: 100,
						},
					},
				},
			},
		},
		{"test update empty list",
			Data{},
			Bar{
				Event: Event{Symbol: "TEST.DE"},
				Close: 100,
			},
			Data{
				list: map[string][]DataEventHandler{
					"TEST.DE": {
						Bar{
							Event: Event{Symbol: "TEST.DE"},
							Close: 100,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc.data.updateList(tc.event)
		if !reflect.DeepEqual(tc.data, tc.expData) {
			t.Errorf("%v updateList(%v): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.event, tc.expData, tc.data)
		}
	}
}
