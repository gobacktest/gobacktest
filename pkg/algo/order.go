package algo

import (
	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

// boolAlgo is a base Algo which return a set bool value
type orderAlgo struct {
	bt.Algo
	order *bt.Order
	value float64
}

// CreateOrder creates an order with a specified value.
func CreateOrder(direction string, value float64) bt.AlgoHandler {
	return &orderAlgo{value: value}
}

// Run runs the algo, returns the bool value of the algo
func (algo orderAlgo) Run(s bt.StrategyHandler) (bool, error) {
	data, _ := s.Data()
	portfolio, _ := s.Portfolio()
	dataEvent, _ := s.Event()
	symbol := dataEvent.Symbol()
	orderType := "MKT"

	event := &bt.Event{}
	event.SetTime(dataEvent.Time())
	event.SetSymbol(symbol)

	initialOrder := &bt.Order{
		Event:     *event,
		Direction: "BOT",
		// Qty should be set by PositionSizer
		OrderType: orderType,
		Limit:     0,
	}

	// fetch latest known data for the symbol
	latest := data.Latest(symbol)

	sizeManager := portfolio.(*bt.Portfolio).SizeManager()
	sizedOrder, err := sizeManager.SizeOrder(initialOrder, latest, portfolio)
	if err != nil {
	}

	riskManager := portfolio.(*bt.Portfolio).RiskManager()
	order, err := riskManager.EvaluateOrder(sizedOrder, latest, portfolio.(*bt.Portfolio).Holdings())
	if err != nil {
	}

	err = s.AddOrder(order)
	if err != nil {
		return false, err
	}

	return true, nil
}
