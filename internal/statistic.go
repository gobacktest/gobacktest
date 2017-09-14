package internal

// StatisticHandler is a basic statistic interface
type StatisticHandler interface {
	Update(DataEvent, PortfolioHandler)
	TrackEvent(e Event)
	Events() []Event
}

// Statistic is a basic test statistic, which holds simple lists of historic events
type Statistic struct {
	eventHistory []Event
}

// Update updates the complete statistics to a given data event
func (s *Statistic) Update(d DataEvent, p PortfolioHandler) {
	return
}

// TrackEvent tracks an event
func (s *Statistic) TrackEvent(e Event) {
	s.eventHistory = append(s.eventHistory, e)
}

// Events returns the complete events history
func (s Statistic) Events() []Event {
	return s.eventHistory
}
