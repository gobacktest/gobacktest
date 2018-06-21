package backtest

import (
	"reflect"
	"testing"
)

func TestNodeName(t *testing.T) {
	var testCases = []struct {
		msg  string
		node *Node
		exp  string
	}{
		{"test getting name of node:",
			&Node{name: "test"},
			"test",
		},
	}

	for _, tc := range testCases {
		name := tc.node.Name()
		if !reflect.DeepEqual(name, tc.exp) {
			t.Errorf("%v Name(): \nexpected %#v, \nactual %#v",
				tc.msg, tc.exp, name)
		}
	}
}

func TestNodeParent(t *testing.T) {
	var testCases = []struct {
		msg   string
		node  *Node
		exp   *Node
		expOk bool
	}{
		{"test getting parent of node:",
			&Node{
				name:   "test",
				parent: &Node{name: "parent"},
			},
			&Node{name: "parent"},
			true,
		},
		{"test getting parent of root node:",
			&Node{
				name: "test",
			},
			&Node{},
			false,
		},
	}

	for _, tc := range testCases {
		parent, ok := tc.node.Parent()
		if !reflect.DeepEqual(parent, tc.exp) || (ok != tc.expOk) {
			t.Errorf("%v Parent(): \nexpected %#v, %v, \nactual   %#v, %v",
				tc.msg, tc.exp, tc.expOk, parent, ok)
		}
	}
}

func TestNodeSetParent(t *testing.T) {
	var testCases = []struct {
		msg    string
		node   *Node
		parent *Node
		exp    *Node
	}{
		{"test setting parent of node:",
			&Node{name: "test"},
			&Node{name: "parent"},
			&Node{
				name:   "test",
				parent: &Node{name: "parent"},
			},
		},
	}

	for _, tc := range testCases {
		node := tc.node.SetParent(tc.parent)
		if !reflect.DeepEqual(node, tc.exp) {
			t.Errorf("%v SetParent(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.exp, node)
		}
	}
}

func TestNodeChildren(t *testing.T) {
	var testCases = []struct {
		msg   string
		node  *Node
		exp   []NodeHandler
		expOk bool
	}{
		{"basic test with nil children:",
			&Node{},
			[]NodeHandler{},
			false,
		},
		{"basic test with one children:",
			&Node{
				children: []NodeHandler{
					&Asset{},
				},
			},
			[]NodeHandler{
				&Asset{},
			},
			true,
		},
		{"basic test with multiple children:",
			&Node{
				children: []NodeHandler{
					&Strategy{},
					&Asset{},
				},
			},
			[]NodeHandler{
				&Strategy{},
				&Asset{},
			},
			true,
		},
	}

	for _, tc := range testCases {
		children, ok := tc.node.Children()
		if !reflect.DeepEqual(children, tc.exp) || (ok != tc.expOk) {
			t.Errorf("%v Children(): \nexpected %#v, %v, \nactual  %#v, %v",
				tc.msg, tc.exp, tc.expOk, children, ok)
		}
	}
}
