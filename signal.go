package gobacktest

// Signal declares a basic signal event
type Signal struct {
	Event
	direction OrderDirection // long or short
}

// Direction returns the Direction of a Signal
func (s Signal) Direction() OrderDirection {
	return s.direction
}

// SetDirection sets the Directions field of a Signal
func (s *Signal) SetDirection(dir OrderDirection) {
	s.direction = dir
}
