package gobacktest

import "time"

// Event defines the basic event interface.
type Event interface {
	Time() time.Time
	Symbol() string
}

// DataEvent represent the different types of data.
type DataEvent struct {
	Data Data
}

// Time returns the timestamp of the data event.
func (d DataEvent) Time() (t time.Time) {
	switch data := d.Data.(type) {
	case Bar:
		return data.Time
	case Tick:
		return data.Time
	}

	return t
}

// Symbol returns the symbol of the data event.
func (d DataEvent) Symbol() string {
	return d.Data.Symbol()
}

// SignalEvent represents the different Signals.
type SignalEvent struct {
	Signal Signal
}

// OrderEvent represents the different Order types.
type OrderEvent struct {
	Order Order
}

// FillEvent prepresent a filled order
type FillEvent struct {
	Fill Fill
}

// EventHandler declares the basic event interface.
type EventHandler interface {
	Reseter
	AppendToQueue(Event)
	NextFromQueue() (Event, bool)
	History() ([]Event, bool)
}

// EventStore is a basic implementation of an event storage system.
type EventStore struct {
	queue   []Event
	history []Event
}

// AppendToQueue appends the event to the end of the queue.
func (es EventStore) AppendToQueue(e Event) {
	es.queue = append(es.queue, e)
}

// NextFromQueue pulls the first event from the queue.
func (es EventStore) NextFromQueue() (e Event, ok bool) {
	// if event queue empty return false
	if len(es.queue) == 0 {
		return e, false
	}

	// return first element from the event queue
	e = es.queue[0]
	es.queue = es.queue[1:]

	// append this event to the history
	es.history = append(es.history, e)

	return e, true
}

// History returns a slice of all historic events.
func (es EventStore) History() (e []Event, ok bool) {
	if len(es.history) == 0 {
		return es.history, false
	}

	return es.history, true
}

// Reset implements the Reseter interface and brings the EventHandler into a clean state.
func (es EventStore) Reset() error {
	es.queue = nil
	es.history = nil

	return nil
}
