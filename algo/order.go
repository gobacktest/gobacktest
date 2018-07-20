package algo

import (
	gbt "github.com/dirkolbrich/gobacktest"
)

// boolAlgo is a base Algo which return a set bool value
type orderAlgo struct {
	gbt.Algo
	order *gbt.Order
	value float64
}

// CreateOrder creates an order with a specified value.
func CreateOrder(direction string, value float64) gbt.AlgoHandler {
	return &orderAlgo{value: value}
}

// Run runs the algo, returns the bool value of the algo
func (algo orderAlgo) Run(s gbt.StrategyHandler) (bool, error) {
	data, _ := s.Data()
	portfolio, _ := s.Portfolio()
	dataEvent, _ := s.Event()
	symbol := dataEvent.Symbol()

	event := &gbt.Event{}
	event.SetTime(dataEvent.Time())
	event.SetSymbol(symbol)

	initialOrder := &gbt.Order{
		Event: *event,
	}
	initialOrder.SetDirection(gbt.BOT)

	// fetch latest known data for the symbol
	latest := data.Latest(symbol)

	sizeManager := portfolio.(*gbt.Portfolio).SizeManager()
	sizedOrder, err := sizeManager.SizeOrder(initialOrder, latest, portfolio)
	if err != nil {
	}

	riskManager := portfolio.(*gbt.Portfolio).RiskManager()
	order, err := riskManager.EvaluateOrder(sizedOrder, latest, portfolio.(*gbt.Portfolio).Holdings())
	if err != nil {
	}

	err = s.AddOrder(order)
	if err != nil {
		return false, err
	}

	return true, nil
}
