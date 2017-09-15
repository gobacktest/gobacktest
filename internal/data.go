package internal

import (
	"sort"
)

// DataHandler is the combined data interface
type DataHandler interface {
	DataLoader
	DataStreamer
}

// DataLoader is the interface loading the data into the data stream
type DataLoader interface {
	Load([]string) error
}

// DataStreamer is the interface returning the data streams
type DataStreamer interface {
	Next() (DataEvent, bool)
	Stream() []DataEvent
	History() []DataEvent
	Latest(string) DataEvent
	List(string) []DataEvent
}

// Data is a basic data struct
type Data struct {
	latest        map[string]DataEvent
	list          map[string][]DataEvent
	stream        []DataEvent
	streamHistory []DataEvent
}

// Load loads data endpoints into a stream
func (d *Data) Load(s []string) error {
	return nil
}

// Stream returns the data stream
func (d *Data) Stream() []DataEvent {
	return d.stream
}

// Next returns the first element of the data stream
// deletes it from the stream and appends it to history
func (d *Data) Next() (event DataEvent, ok bool) {
	// check for element in datastream
	if len(d.stream) == 0 {
		return event, false
	}

	event = d.stream[0]
	d.stream = d.stream[1:] // delete first element from stream
	d.streamHistory = append(d.stream, event)

	// update list of current data events
	d.updateLatest(event)
	// update list of data events for single symbol
	d.updateList(event)

	return event, true
}

// History returns the historic data stream
func (d *Data) History() []DataEvent {
	return d.streamHistory
}

// Latest returns the last known data event for a symbol.
func (d *Data) Latest(symbol string) DataEvent {
	return d.latest[symbol]
}

// updateCurrent puts the last current data event to the current list.
func (d *Data) updateLatest(event DataEvent) {
	// check for nil map, else initialise the map
	if d.latest == nil {
		d.latest = make(map[string]DataEvent)
	}

	d.latest[event.Symbol()] = event
}

// List returns the data event list for a symbol.
func (d *Data) List(symbol string) []DataEvent {
	return d.list[symbol]
}

// updateList appends an event to the data list.
func (d *Data) updateList(event DataEvent) {
	// Check for nil map, else initialise the map
	if d.list == nil {
		d.list = make(map[string][]DataEvent)
	}

	d.list[event.Symbol()] = append(d.list[event.Symbol()], event)
}

// sortStream sorts the dataStream
func (d *Data) sortStream() {
	sort.Slice(d.stream, func(i, j int) bool {
		b1 := d.stream[i]
		b2 := d.stream[j]

		// if date is equal sort by symbol
		if b1.Timestamp().Equal(b2.Timestamp()) {
			return b1.Symbol() < b2.Symbol()
		}
		// else sort by date
		return b1.Timestamp().Before(b2.Timestamp())
	})
}
