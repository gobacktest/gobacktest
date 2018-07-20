package gobacktest

// StrategyHandler is a basic strategy interface.
type StrategyHandler interface {
	Data() (DataHandler, bool)
	SetData(d DataHandler) error
	Portfolio() (PortfolioHandler, bool)
	SetPortfolio(p PortfolioHandler) error
	Event() (DataEvent, bool)
	SetEvent(DataEvent) error
	Orders() ([]OrderEvent, bool)
	AddOrder(...OrderEvent) error
	Strategies() ([]StrategyHandler, bool)
	Assets() ([]*Asset, bool)
	OnData(DataEvent) ([]OrderEvent, error)
}

// Strategy implements NodeHandler via Node, used as a strategy building block.
type Strategy struct {
	Node
	algos     AlgoStack
	data      DataHandler
	portfolio PortfolioHandler
	event     DataEvent
	orders    []OrderEvent
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

// Orders the orderbook, a slice of all known orders.
func (s *Strategy) Orders() ([]OrderEvent, bool) {
	if len(s.orders) == 0 {
		return s.orders, false
	}

	return s.orders, true
}

// AddOrder sets the data property.
func (s *Strategy) AddOrder(orders ...OrderEvent) error {
	for _, order := range orders {
		s.orders = append(s.orders, order)
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
func (s *Strategy) OnData(event DataEvent) (orders []OrderEvent, err error) {
	s.SetEvent(event)

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
		return nil, err
	}

	// pass data event down to child strategies
	if strategies, ok := s.Strategies(); ok {
		for _, strategy := range strategies {
			orders, err := strategy.OnData(event)
			if err != nil {
				return nil, err
			}
			s.AddOrder(orders...)
		}
	}

	orders, ok = s.Orders()
	if !ok {
		return []OrderEvent{}, nil
	}

	return orders, nil
}
