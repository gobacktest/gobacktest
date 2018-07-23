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

// CreateSignal creates an signal with a specified direktion.
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

	switch algo.direction {
	case "buy":
		signal.SetDirection(gbt.BOT)
	case "sell":
		signal.SetDirection(gbt.SLD)
	case "exit":
		signal.SetDirection(gbt.EXT)
	}

	err := s.AddSignal(signal)
	if err != nil {
		return false, err
	}

	return true, nil
}
