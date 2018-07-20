package gobacktest

import (
	"reflect"
	"testing"
)

func TestDataReset(t *testing.T) {
	var testCases = []struct {
		msg     string
		data    DataHandler
		expData DataHandler
	}{
		{"test with empty data stream",
			&Data{
				latest: map[string]DataEvent{
					"TEST.DE": &Bar{Close: 100},
				},
				list: map[string][]DataEvent{
					"TEST.DE": {
						&Bar{Close: 100},
						&Bar{Close: 110},
						&Bar{Close: 95},
					},
				},
				stream: []DataEvent{},
				streamHistory: []DataEvent{
					&Bar{Close: 100},
					&Bar{Close: 110},
					&Bar{Close: 95},
				},
			},
			&Data{
				stream: []DataEvent{
					&Bar{Close: 100},
					&Bar{Close: 110},
					&Bar{Close: 95},
				},
			},
		},
		{"test with empty data",
			&Data{},
			&Data{},
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

func TestDataNext(t *testing.T) {
	var testCases = []struct {
		msg      string
		data     DataHandler
		expData  DataHandler
		expEvent DataEvent
		expOk    bool
	}{
		{"testing multiple data events",
			&Data{
				stream: []DataEvent{
					&Bar{Event: Event{symbol: "TEST.DE"}, Close: 110},
					&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
					&Bar{Event: Event{symbol: "TEST.DE"}, Close: 90},
				},
			},
			&Data{
				latest: map[string]DataEvent{
					"TEST.DE": &Bar{Event: Event{symbol: "TEST.DE"}, Close: 110},
				},
				list: map[string][]DataEvent{
					"TEST.DE": {
						&Bar{Event: Event{symbol: "TEST.DE"}, Close: 110},
					},
				},
				stream: []DataEvent{
					&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
					&Bar{Event: Event{symbol: "TEST.DE"}, Close: 90},
				},
				streamHistory: []DataEvent{
					&Bar{Event: Event{symbol: "TEST.DE"}, Close: 110},
				},
			},
			&Bar{Event: Event{symbol: "TEST.DE"}, Close: 110},
			true,
		},
		{"testing single data events",
			&Data{
				stream: []DataEvent{
					&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
				},
			},
			&Data{
				latest: map[string]DataEvent{
					"TEST.DE": &Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
				},
				list: map[string][]DataEvent{
					"TEST.DE": {
						&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
					},
				},
				stream: []DataEvent{},
				streamHistory: []DataEvent{
					&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
				},
			},
			&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
			true,
		},
		{"testing empty data events",
			&Data{},
			&Data{},
			nil,
			false,
		},
	}
	for _, tc := range testCases {
		event, ok := tc.data.Next()
		if !reflect.DeepEqual(event, tc.expEvent) || (ok != tc.expOk) {
			t.Errorf("%v Next(): \nexpected %#v %v, \nactual   %#v %v",
				tc.msg, tc.expEvent, tc.expOk, event, ok)
		}
		if !reflect.DeepEqual(tc.data, tc.expData) {
			t.Errorf("%v Next(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expData, tc.data)
		}
	}
}

func TestUpdateLatest(t *testing.T) {
	var testCases = []struct {
		msg     string
		data    Data
		event   DataEvent
		expData Data
	}{
		{"test update filled latest",
			Data{
				latest: map[string]DataEvent{
					"TEST.DE": &Bar{Event: Event{symbol: "TEST.DE"}, Close: 80},
				},
			},
			&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
			Data{
				latest: map[string]DataEvent{
					"TEST.DE": &Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
				},
			},
		},
		{"test update empty latest",
			Data{},
			&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
			Data{
				latest: map[string]DataEvent{
					"TEST.DE": &Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
				},
			},
		},
		{"test update filled latest with other symbol data",
			Data{
				latest: map[string]DataEvent{
					"TEST.DE": &Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
				},
			},
			&Bar{Event: Event{symbol: "BAS.DE"}, Close: 90},
			Data{
				latest: map[string]DataEvent{
					"TEST.DE": &Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
					"BAS.DE":  &Bar{Event: Event{symbol: "BAS.DE"}, Close: 90},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc.data.updateLatest(tc.event)
		if !reflect.DeepEqual(tc.data.latest, tc.expData.latest) {
			t.Errorf("%v updateLatest(%v): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.event, tc.expData.latest, tc.data.latest)
		}
	}
}

func TestUpdateList(t *testing.T) {
	var testCases = []struct {
		msg     string
		data    Data
		event   DataEvent
		expData Data
	}{
		{"test update filled list",
			Data{
				list: map[string][]DataEvent{
					"TEST.DE": {
						&Bar{Event: Event{symbol: "TEST.DE"}, Close: 90},
					},
				},
			},
			&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
			Data{
				list: map[string][]DataEvent{
					"TEST.DE": {
						&Bar{Event: Event{symbol: "TEST.DE"}, Close: 90},
						&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
					},
				},
			},
		},
		{"test update empty list",
			Data{},
			&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
			Data{
				list: map[string][]DataEvent{
					"TEST.DE": {
						&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
					},
				},
			},
		},
		{"test update filled list with other symbol data",
			Data{
				list: map[string][]DataEvent{
					"TEST.DE": {
						&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
					},
				},
			},
			&Bar{Event: Event{symbol: "BAS.DE"}, Close: 90},
			Data{
				list: map[string][]DataEvent{
					"TEST.DE": {
						&Bar{Event: Event{symbol: "TEST.DE"}, Close: 100},
					},
					"BAS.DE": {
						&Bar{Event: Event{symbol: "BAS.DE"}, Close: 90},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc.data.updateList(tc.event)
		if !reflect.DeepEqual(tc.data.list, tc.expData.list) {
			t.Errorf("%v updateList(%v): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.event, tc.expData.list, tc.data.list)
		}
	}
}
