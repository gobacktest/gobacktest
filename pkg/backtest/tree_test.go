package backtest

import (
	"reflect"
	"testing"
)

func TestNodeName(t *testing.T) {
	var testCases = []struct {
		msg  string
		node Node
		exp  string
	}{
		{"test getting name of node:",
			Node{name: "test"},
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
