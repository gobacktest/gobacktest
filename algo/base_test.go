package algo

import (
	"reflect"
	"testing"

	gbt "github.com/dirkolbrich/gobacktest"
)

func TestBoolAlgo(t *testing.T) {
	var testCases = []struct {
		msg     string
		option  bool
		expBool bool
		expErr  error
	}{
		{"test true option",
			true,
			true,
			nil,
		},
		{"test false option",
			false,
			false,
			nil,
		},
	}

	for _, tc := range testCases {
		algo := BoolAlgo(tc.option)
		ok, err := algo.Run(&gbt.Strategy{})
		if (ok != tc.expBool) || !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("%v: BoolAlgo(%v): \nexpected %#v %v, \nactual   %#v %v", tc.msg, tc.option, tc.expBool, tc.expErr, ok, err)
		}
	}

}
