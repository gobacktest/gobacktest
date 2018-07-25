package algo

import (
	gbt "github.com/dirkolbrich/gobacktest"
)

// signalAlgo creates a signal
type signalAlgo struct {
	gbt.Algo
	signal    *gbt.Signal
	direction string
}

// CreateSignal creates a signal with a specified direction.
func CreateSignal(direction string) gbt.AlgoHandler {
	return &signalAlgo{direction: direction}
}

// Run runs the algo, returns the bool value of the algo
func (algo signalAlgo) Run(s gbt.StrategyHandler) (bool, error) {
	// data, _ := s.Data()
	dataEvent, _ := s.Event()
	symbol := dataEvent.Symbol()

	event := &gbt.Event{}
	event.SetTime(dataEvent.Time())
	event.SetSymbol(symbol)

	signal := &gbt.Signal{
		Event: *event,
	}

	switch {
	case (algo.direction == "buy") || (algo.direction == "long") || (algo.direction == "BOT"):
		signal.SetDirection(gbt.BOT)
	case (algo.direction == "sell") || (algo.direction == "short") || (algo.direction == "SLD"):
		signal.SetDirection(gbt.SLD)
	case (algo.direction == "exit") || (algo.direction == "EXT"):
		signal.SetDirection(gbt.EXT)
	default:
		signal.SetDirection(gbt.HLD)
	}

	err := s.AddSignal(signal)
	if err != nil {
		return false, err
	}

	return true, nil
}
