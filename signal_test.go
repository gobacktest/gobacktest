package gobacktest

import (
	"reflect"
	"testing"
)

func TestSignalSetDirection(t *testing.T) {
	var testCases = []struct {
		msg       string
		signal    Signal
		dir       Direction
		expSignal Signal
	}{
		{"simple direction:",
			Signal{},
			BOT,
			Signal{direction: BOT},
		},
	}

	for _, tc := range testCases {
		tc.signal.SetDirection(tc.dir)
		if !reflect.DeepEqual(tc.signal, tc.expSignal) {
			t.Errorf("%v SetDirection(%v): \nexpected %#v, \nactual %#v",
				tc.msg, tc.dir, tc.expSignal, tc.signal)
		}
	}
}

func TestSignalGetDirection(t *testing.T) {
	var testCases = []struct {
		msg    string
		signal Signal
		expDir Direction
	}{
		{"simple direction:",
			Signal{direction: BOT},
			BOT,
		},
	}

	for _, tc := range testCases {
		dir := tc.signal.Direction()
		if dir != tc.expDir {
			t.Errorf("%v Direction(): \nexpected %#v, \nactual %#v",
				tc.msg, tc.expDir, dir)
		}
	}
}
