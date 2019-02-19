package gobacktest

// SignalDirection defines which direction a signal indicates.
type SignalDirection int

// different types of order directions
const (
	// Buy
	BuySignal SignalDirection = iota // 0
	// Sell
	SellSignal
	// Hold
	HoldSignal
	// Exit
	ExitSignal
)

func (sd SignalDirection) String() string {
	switch sd {
	case BuySignal:
		return "BOT, buy direction"
	case SellSignal:
		return "SLD, sell direction"
	case HoldSignal:
		return "HLD, hold direction"
	case ExitSignal:
		return "EXT, exit direction"
	}

	return ""
}

// Signal declares a basic signal event.
type Signal struct {
	direction SignalDirection // long, short, exit or hold
}

// Direction returns the direction of a Signal.
func (s Signal) Direction() SignalDirection {
	return s.direction
}

// SetDirection sets the directions of a Signal.
func (s *Signal) SetDirection(dir SignalDirection) {
	s.direction = dir
}
