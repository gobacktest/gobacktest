package backtest

import (
	"reflect"
	"testing"
)

func TestNewAsset(t *testing.T) {
	var testCases = []struct {
		msg  string
		name string
		exp  *Asset
	}{
		{"setup new asset:",
			"test",
			&Asset{
				Node{name: "test"},
			},
		},
	}

	for _, tc := range testCases {
		asset := NewAsset(tc.name)
		if !reflect.DeepEqual(asset, tc.exp) {
			t.Errorf("%v NewAsset(%s): \nexpected %#v, \nactual %#v",
				tc.msg, tc.name, tc.exp, asset)
		}
	}
}

func TestAssetChildren(t *testing.T) {
	var testCases = []struct {
		msg       string
		asset     Asset
		expReturn []NodeHandler
	}{
		{"basic test with nil children:",
			Asset{},
			[]NodeHandler{},
		},
		{"basic test with one children:",
			Asset{
				Node{
					children: []NodeHandler{
						&Asset{},
					},
				},
			},
			[]NodeHandler{},
		},
		{"basic test with multiple children:",
			Asset{
				Node{
					children: []NodeHandler{
						&Strategy{},
						&Asset{},
					},
				},
			},
			[]NodeHandler{},
		},
	}

	for _, tc := range testCases {
		children := tc.asset.Children()
		if !reflect.DeepEqual(children, tc.expReturn) {
			t.Errorf("%v Children(): \nexpected %#v, \nactual %#v",
				tc.msg, tc.expReturn, children)
		}
	}
}

func TestAssetSetChildren(t *testing.T) {
	var testCases = []struct {
		msg   string
		asset Asset
		child NodeHandler
		exp   NodeHandler
	}{
		{"test setting single children:",
			Asset{},
			&Strategy{},
			&Asset{},
		},
	}

	for _, tc := range testCases {
		asset := tc.asset.SetChildren(tc.child)
		if !reflect.DeepEqual(asset, tc.exp) {
			t.Errorf("%v SetChildren(): \nexpected %#v, \nactual %#v",
				tc.msg, tc.exp, asset)
		}
	}
}
