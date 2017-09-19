package backtest

import (
	"testing"
	"time"
)

type testEvent struct {
}

func (t testEvent) IsEvent() bool {
	return true
}

func (t testEvent) GetTime() time.Time {
	return time.Now()
}

func (t testEvent) GetSymbol() string {
	return "testEvent"
}

// queueTests is a table for testing the event queue
var queueTests = []struct {
	test     Test           // Test struct
	expEvent EventHandler // expected Event interface return
	expBool  bool           // expected bool return
}{
	{Test{}, nil, false}, // Test.eventQueue is empty
	{Test{
		eventQueue: []EventHandler{
			testEvent{},
		},
	}, testEvent{}, true},
}

func TestNextEvent(t *testing.T) {
	for _, tt := range queueTests {
		event, ok := tt.test.nextEvent()
		if ok != tt.expBool {
			t.Errorf("nextEvent(): expected %v %v, actual %v %v", tt.expEvent, tt.expBool, event, ok)
		}
	}
}
