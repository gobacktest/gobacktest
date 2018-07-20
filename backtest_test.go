package gobacktest

import (
	"testing"
	"time"
)

// setup mock for an event
type testEvent struct {
}

func (t testEvent) Time() time.Time {
	return time.Now()
}

func (t *testEvent) SetTime(time time.Time) {
}

func (t testEvent) Symbol() string {
	return "testEvent"
}

func (t *testEvent) SetSymbol(s string) {
}

// queueTests is a table for testing the event queue
var queueTests = []struct {
	test     Backtest     // Test struct
	expEvent EventHandler // expected Event interface return
	expBool  bool         // expected bool return
}{
	{Backtest{}, nil, false}, // Test.eventQueue is empty
	{Backtest{
		eventQueue: []EventHandler{
			&testEvent{},
		},
	}, &testEvent{}, true},
}

func TestNextEvent(t *testing.T) {
	for _, tt := range queueTests {
		event, ok := tt.test.nextEvent()
		if ok != tt.expBool {
			t.Errorf("nextEvent(): expected %v %v, actual %v %v", tt.expEvent, tt.expBool, event, ok)
		}
	}
}
