package internal

import "time"

// StrategyHandler is a basic strategy interface
type StrategyHandler interface {
	CalculateSignal(EventHandler) (SignalEvent, bool)
}

// SimpleStrategy is a basic test strategy, which interprets every DataEvent as a signal to buy
type SimpleStrategy struct {
	eventStream []EventHandler
}

// CalculateSignal handles the single Event
func (s *SimpleStrategy) CalculateSignal(event EventHandler) (signal SignalEvent, ok bool) {
	// log.Printf("reveived event, adding to eventStream: %#v\n", event)
	s.eventStream = append(s.eventStream, event)

	// type switch for event type
	switch ev := event.(type) {
	case BarEvent:
		signal = SignalEvent{
			Event:        Event{timestamp: time.Now(), symbol: ev.Symbol()},
			Direction:    "long",
		}
	}

	return signal, true
}
