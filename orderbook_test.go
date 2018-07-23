package gobacktest

import (
	"errors"
	"reflect"
	"testing"
)

func TestOrderbookAdd(t *testing.T) {
	var testCases = []struct {
		msg    string
		order  *Order
		ob     OrderBook
		expOB  OrderBook
		expErr error
	}{
		{
			msg:   "add order to empty orderbook:",
			order: &Order{},
			ob:    OrderBook{},
			expOB: OrderBook{
				counter: 1,
				orders: []OrderEvent{
					&Order{id: 1},
				},
			},
			expErr: nil,
		},
		{
			msg:   "add second order:",
			order: &Order{},
			ob: OrderBook{
				counter: 1,
				orders: []OrderEvent{
					&Order{id: 1},
				},
			},
			expOB: OrderBook{
				counter: 2,
				orders: []OrderEvent{
					&Order{id: 1},
					&Order{id: 2},
				},
			},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		err := tc.ob.Add(tc.order)
		if !reflect.DeepEqual(tc.ob, tc.expOB) || !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("%v Add(%v): \nexpected %#v %#v\nactual   %#v %#v", tc.msg, tc.order, tc.expOB, tc.expErr, tc.ob, err)
		}
	}
}

func TestOrderbookRemove(t *testing.T) {
	var testCases = []struct {
		msg    string
		id     int
		ob     OrderBook
		expOB  OrderBook
		expErr error
	}{
		{
			msg:    "remove from empty orderbook:",
			id:     1,
			ob:     OrderBook{},
			expOB:  OrderBook{},
			expErr: errors.New("order with id 1 not found"),
		},
		{
			msg: "remove order with invalid id:",
			id:  2,
			ob: OrderBook{
				counter: 1,
				orders: []OrderEvent{
					&Order{id: 1},
				},
			},
			expOB: OrderBook{
				counter: 1,
				orders: []OrderEvent{
					&Order{id: 1},
				},
			},
			expErr: errors.New("order with id 2 not found"),
		},
		{
			msg: "remove single order from orderbook:",
			id:  1,
			ob: OrderBook{
				counter: 1,
				orders: []OrderEvent{
					&Order{id: 1},
				},
			},
			expOB: OrderBook{
				counter: 1,
				orders:  []OrderEvent{},
				history: []OrderEvent{
					&Order{id: 1},
				},
			},
			expErr: nil,
		},
		{
			msg: "remove order from multiple order:",
			id:  2,
			ob: OrderBook{
				counter: 3,
				orders: []OrderEvent{
					&Order{id: 1},
					&Order{id: 2},
					&Order{id: 3},
				},
			},
			expOB: OrderBook{
				counter: 3,
				orders: []OrderEvent{
					&Order{id: 1},
					&Order{id: 3},
				},
				history: []OrderEvent{
					&Order{id: 2},
				},
			},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		err := tc.ob.Remove(tc.id)
		if !reflect.DeepEqual(tc.ob, tc.expOB) || !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("%v Remove(%v): \nexpected %#v %#v\nactual   %#v %#v", tc.msg, tc.id, tc.expOB, tc.expErr, tc.ob, err)
		}
	}
}
