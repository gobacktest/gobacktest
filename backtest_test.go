package gobacktest

import (
	"testing"

	"github.com/dirkolbrich/gobacktest/internal"
)

// queueTests is a table for testing parsing bar data into a BarEvent
var queueTests = []struct {
	test      Test // start input
	expReturn bool // expected return
}{
	{Test{}, true},
	{Test{
		eventQueue: []internal.EventHandler{
			internal.BarEvent{},
		},
	}, false},
}

func TestQueueIsEmpty(t *testing.T) {

	for _, tt := range queueTests {
		ok := tt.test.queueIsEmpty()
		if ok != tt.expReturn {
			t.Errorf("queueIsEmpty(): expected %v, actual %v", tt.expReturn, ok)
		}
	}
}
