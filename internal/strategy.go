package internal

// Strategy is a basic strategy interface
type Strategy interface {
	CalculateSignal(EventHandler) (EventHandler, error)
}

// SimpleStrategy is a basic test strategy, which interprets every DataEvent as a sign to by
type SimpleStrategy struct {
	events []EventHandler
}

// CalculateSignal handles the single Event
func (s SimpleStrategy) CalculateSignal(e EventHandler) (event EventHandler, err error) {
	return
}
