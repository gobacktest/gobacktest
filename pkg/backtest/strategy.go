package backtest

// StrategyHandler is a basic strategy interface.
type StrategyHandler interface {
	Data() (DataHandler, bool)
	SetData(d DataHandler) error
	Portfolio() (PortfolioHandler, bool)
	SetPortfolio(p PortfolioHandler) error
	Event() (DataEventHandler, bool)
	Strategies() ([]StrategyHandler, bool)
	Assets() ([]*Asset, bool)
	OnData(DataEventHandler) (SignalEvent, error)
}

// Strategy implements NodeHandler via Node, used as a strategy building block.
type Strategy struct {
	Node
	algos     AlgoStack
	data      DataHandler
	portfolio PortfolioHandler
	event     DataEventHandler
}

// NewStrategy return a new strategy node ready to use.
func NewStrategy(name string) *Strategy {
	var s = &Strategy{}
	s.SetName(name)
	s.SetRoot(true)
	return s
}

// Data returns the underlying data property.
func (s *Strategy) Data() (DataHandler, bool) {
	if s.data == nil {
		return nil, false
	}

	return s.data, true
}

// SetData sets the data property.
func (s *Strategy) SetData(data DataHandler) error {
	s.data = data

	// check for sub strategies and set their data as well
	subStrategies, _ := s.Strategies()

	for _, sub := range subStrategies {
		err := sub.SetData(data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Portfolio returns the underlying portfolio property.
func (s *Strategy) Portfolio() (PortfolioHandler, bool) {
	if s.portfolio == nil {
		return nil, false
	}

	return s.portfolio, true
}

// SetPortfolio sets the portfolio property.
func (s *Strategy) SetPortfolio(portfolio PortfolioHandler) error {
	s.portfolio = portfolio

	// check for sub strategies and set their portfolio as well
	subStrategies, _ := s.Strategies()

	for _, sub := range subStrategies {
		err := sub.SetPortfolio(portfolio)
		if err != nil {
			return err
		}
	}

	return nil
}

// Event returns the underlying data property.
func (s *Strategy) Event() (DataEventHandler, bool) {
	if s.event == nil {
		return nil, false
	}

	return s.event, true
}

// SetAlgo sets the algo stack for the Strategy
func (s *Strategy) SetAlgo(algos ...AlgoHandler) *Strategy {
	for _, algo := range algos {
		s.algos.stack = append(s.algos.stack, algo)
	}
	return s
}

// Strategies return all children which are a strategy.
func (s *Strategy) Strategies() ([]StrategyHandler, bool) {
	var strategies []StrategyHandler

	// get all children
	children, ok := s.Children()

	// no children means no sub strategies
	if !ok {
		return strategies, false
	}

	// check each child if it is a strategy
	for _, child := range children {
		switch c := child.(type) {
		case *Strategy:
			strategies = append(strategies, c)
		}
	}

	// no sub strategies in children
	if len(strategies) == 0 {
		return strategies, false
	}

	return strategies, true
}

// Assets return all children which are a strategy.
func (s *Strategy) Assets() ([]*Asset, bool) {
	var assets []*Asset

	// get all children
	children, ok := s.Children()

	// no children means no sub strategies
	if !ok {
		return assets, false
	}

	// check each child if it is a strategy
	for _, child := range children {
		switch c := child.(type) {
		case *Asset:
			assets = append(assets, c)
		}
	}

	// no sub strategies in children
	if len(assets) == 0 {
		return assets, false
	}

	return assets, true
}

// OnData handles an incoming data event. It runs the algo stack on this data.
func (s *Strategy) OnData(event DataEventHandler) (SignalEvent, error) {
	s.event = event

	// create Signal
	se := &Signal{}

	// // type switch for event type
	// switch e := event.(type) {
	// case *Bar:
	// 	// fill Signal
	// 	se.Event = Event{Timestamp: e.GetTime(), Symbol: e.GetSymbol()}
	// 	se.Direction = "long"
	// }

	// run the algo stack of this strategy
	ok, err := s.algos.Run(s)
	if !ok {
		return se, err
	}

	// pass data event down to child strategies
	if strategies, ok := s.Strategies(); ok {
		for _, strategy := range strategies {
			strategy.OnData(event)
		}
	}

	return se, nil
}
