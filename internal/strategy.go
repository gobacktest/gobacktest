package internal

import "time"

// StrategyHandler is a basic strategy interface
type StrategyHandler interface {
	CalculateSignal(DataEvent, DataHandler) (SignalEvent, error)
}

// SimpleStrategy is a basic test strategy, which interprets every DataEvent as a signal to buy
type SimpleStrategy struct {
	eventStream []Event
}

// CalculateSignal handles the single Event
func (s *SimpleStrategy) CalculateSignal(e DataEvent, data DataHandler) (se SignalEvent, err error) {
	// log.Printf("reveived event, adding to eventStream: %#v\n", event)
	s.eventStream = append(s.eventStream, e)

	// type switch for event type
	switch e := e.(type) {
	case bar:
		se = &signal{
			event:     event{timestamp: time.Now(), symbol: e.Symbol()},
			direction: "long",
		}
	}

	return se, nil
}
