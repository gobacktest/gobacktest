package gobacktest

// Direction defines which direction a signal indicates
type Direction int

// different types of order directions
const (
	// Buy
	BOT Direction = iota // 0
	// Sell
	SLD
	// Hold
	HLD
	// Exit
	EXT
)

// Signal declares a basic signal event
type Signal struct {
	Event
	direction Direction // long, short, exit or hold
}

// Direction returns the Direction of a Signal
func (s Signal) Direction() Direction {
	return s.direction
}

// SetDirection sets the Directions field of a Signal
func (s *Signal) SetDirection(dir Direction) {
	s.direction = dir
}
