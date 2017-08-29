package internal

import "time"

// StrategyHandler is a basic strategy interface
type StrategyHandler interface {
	CalculateSignal(DataEvent) (SignalEvent, bool)
}

// SimpleStrategy is a basic test strategy, which interprets every DataEvent as a signal to buy
type SimpleStrategy struct {
	eventStream []Event
}

// CalculateSignal handles the single Event
func (s *SimpleStrategy) CalculateSignal(data DataEvent) (se SignalEvent, ok bool) {
	// log.Printf("reveived event, adding to eventStream: %#v\n", event)
	s.eventStream = append(s.eventStream, data)

	// type switch for event type
	switch ev := data.(type) {
	case bar:
		se = &signal{
			event:     event{timestamp: time.Now(), symbol: ev.Symbol()},
			direction: "long",
		}
	}

	return se, true
}
