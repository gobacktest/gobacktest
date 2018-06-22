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

func TestNodeRoot(t *testing.T) {
	var testCases = []struct {
		msg  string
		node *Node
		exp  bool
	}{
		{"test root node:",
			&Node{name: "test", root: true},
			true,
		},
		{"test child node:",
			&Node{name: "test", root: false},
			false,
		},
	}

	for _, tc := range testCases {
		ok := tc.node.Root()
		if ok != tc.exp {
			t.Errorf("%v Root(): \nexpected %v, \nactual %v",
				tc.msg, tc.exp, ok)
		}
	}
}
func TestNodeSetRoot(t *testing.T) {
	var testCases = []struct {
		msg  string
		node *Node
		root bool
		exp  *Node
	}{
		{"test set root true on non root:",
			&Node{name: "test", root: false},
			true,
			&Node{name: "test", root: true},
		},
		{"test set root true on root:",
			&Node{name: "test", root: true},
			true,
			&Node{name: "test", root: true},
		},
		{"test set root false on root:",
			&Node{name: "test", root: true},
			false,
			&Node{name: "test", root: false},
		},
		{"test set root false on non root:",
			&Node{name: "test", root: false},
			false,
			&Node{name: "test", root: false},
		},
	}

	for _, tc := range testCases {
		tc.node.SetRoot(tc.root)
		if !reflect.DeepEqual(tc.node, tc.exp) {
			t.Errorf("%v SetRoot(%v): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.root, tc.exp, tc.node)
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
		{"test with nil children:",
			&Node{},
			[]NodeHandler{},
			false,
		},
		{"test with one children:",
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
		{"test with multiple children:",
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

func TestNodeSingleSetChildren(t *testing.T) {
	expNode := &Node{
		name: "test",
		root: true,
		children: []NodeHandler{
			&Node{
				name: "child",
				root: false,
			},
		},
	}
	node := &Node{name: "test", root: true}
	child := &Node{name: "child", root: true}
	testNode := node.SetChildren(child)

	if !reflect.DeepEqual(testNode, expNode) {
		t.Errorf("set single child on node SetChildren(): \nexpected %#v, \nactual %#v", expNode, node)
	}
}

func TestNodeMultipleSetChildren(t *testing.T) {
	expNode := &Node{
		name: "test",
		root: true,
		children: []NodeHandler{
			&Node{
				name: "child1",
				root: false,
			},
			&Node{
				name: "child2",
				root: false,
			},
		},
	}
	node := &Node{name: "test", root: true}
	child1 := &Node{name: "child1", root: true}
	child2 := &Node{name: "child2", root: true}
	testNode := node.SetChildren(child1, child2)

	if !reflect.DeepEqual(testNode, expNode) {
		t.Errorf("set single child on node SetChildren(): \nexpected %#v, \nactual %#v", expNode, node)
	}
}
