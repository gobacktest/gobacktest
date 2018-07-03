package algo

import (
	"testing"

	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

func TestAlgoTrue(t *testing.T) {
	algo := &TrueAlgo{}

	ok, err := algo.Run(&bt.Strategy{})
	if !ok || (err != nil) {
		t.Errorf("TrueAlgo(): \nexpected %#v %v, \nactual   %#v %v", nil, true, err, ok)
	}
}

func TestAlgoFalse(t *testing.T) {
	algo := &FalseAlgo{}

	ok, err := algo.Run(&bt.Strategy{})
	if ok || (err != nil) {
		t.Errorf("FalseAlgo(): \nexpected %#v %v, \nactual   %#v %v", nil, false, err, ok)
	}
}
