package gobacktest

// StrategyHandler is a basic strategy interface.
type StrategyHandler interface {
	Data() (DataHandler, bool)
	SetData(d DataHandler) error
	Portfolio() (PortfolioHandler, bool)
	SetPortfolio(p PortfolioHandler) error
	Event() (DataEvent, bool)
	SetEvent(DataEvent) error
	Signals() ([]SignalEvent, bool)
	AddSignal(...SignalEvent) error
	Strategies() ([]StrategyHandler, bool)
	Assets() ([]*Asset, bool)
	OnData(DataEvent) ([]SignalEvent, error)
}

// Strategy implements NodeHandler via Node, used as a strategy building block.
type Strategy struct {
	Node
	algos     AlgoStack
	data      DataHandler
	portfolio PortfolioHandler
	event     DataEvent
	signals   []SignalEvent
}

// NewStrategy return a new strategy node ready to use.
func NewStrategy(name string) *Strategy {
	var s = &Strategy{}
	s.SetName(name)
	s.SetRoot(true)
	s.SetWeight(1)
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
func (s *Strategy) Event() (DataEvent, bool) {
	if s.event == nil {
		return nil, false
	}

	return s.event, true
}

// SetEvent sets the event property.
func (s *Strategy) SetEvent(event DataEvent) error {
	s.event = event
	return nil
}

// Signals returns a slice of all from th ealgo loop created signals.
func (s *Strategy) Signals() ([]SignalEvent, bool) {
	if len(s.signals) == 0 {
		return s.signals, false
	}

	return s.signals, true
}

// AddSignal sets the data property.
func (s *Strategy) AddSignal(signals ...SignalEvent) error {
	for _, signal := range signals {
		s.signals = append(s.signals, signal)
	}
	return nil
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
func (s *Strategy) OnData(event DataEvent) (signals []SignalEvent, err error) {
	s.SetEvent(event)

	// run the algo stack of this strategy
	ok, err := s.algos.Run(s)
	if !ok {
		return nil, err
	}

	// pass data event down to child strategies
	if strategies, ok := s.Strategies(); ok {
		for _, strategy := range strategies {
			signals, err := strategy.OnData(event)
			if err != nil {
				return nil, err
			}
			s.AddSignal(signals...)
		}
	}

	signals, ok = s.Signals()
	if !ok {
		return nil, nil
	}

	return signals, nil
}
