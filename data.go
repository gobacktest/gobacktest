package gobacktest

import (
	"sort"
)

// DataHandler is the combined data interface.
type DataHandler interface {
	DataLoader
	DataStreamer
	Reseter
}

// DataLoader interface definse how to load data into the data stream.
type DataLoader interface {
	Load([]string) error
}

// DataStreamer is the interface defining data streams functionality.
type DataStreamer interface {
	Next() (DataEvent, bool)
	Stream() []DataEvent
	History() []DataEvent
	Latest(string) DataEvent
	List(string) []DataEvent
}

// Data is a basic data struct.
type Data struct {
	latest  map[string]DataEvent
	list    map[string][]DataEvent
	stream  []DataEvent
	history []DataEvent
}

// Load data events into a stream.
// This method satisfies the DataLoader interface, but should be overwritten
// by the specific data loading implementation.
func (d *Data) Load(s []string) error {
	return nil
}

// Reset implements the Reseter interface and resets the data struct to a clean state with loaded data events.
func (d *Data) Reset() error {
	d.latest = nil
	d.list = nil
	d.stream = d.history
	d.history = nil
	return nil
}

// Stream returns the data stream.
func (d *Data) Stream() []DataEvent {
	return d.stream
}

// SetStream sets the data stream.
func (d *Data) SetStream(stream []DataEvent) {
	d.stream = stream
}

// Next returns the first element of the data stream,
// deletes it from the data stream and appends it to the historic data stream.
func (d *Data) Next() (dh DataEvent, ok bool) {
	// check for element in datastream
	if len(d.stream) == 0 {
		return dh, false
	}

	dh = d.stream[0]
	d.stream = d.stream[1:] // delete first element from stream
	d.history = append(d.history, dh)

	// update list of current data events
	d.updateLatest(dh)
	// update list of data events for single symbol
	d.updateList(dh)

	return dh, true
}

// History returns the historic data stream.
func (d *Data) History() []DataEvent {
	return d.history
}

// Latest returns the last known data event for a symbol.
func (d *Data) Latest(symbol string) DataEvent {
	return d.latest[symbol]
}

// List returns the data event list for a symbol.
func (d *Data) List(symbol string) []DataEvent {
	return d.list[symbol]
}

// SortStream sorts the data stream in ascending order.
func (d *Data) SortStream() {
	sort.Slice(d.stream, func(i, j int) bool {
		b1 := d.stream[i]
		b2 := d.stream[j]

		// if date is equal sort by symbol
		if b1.Time().Equal(b2.Time()) {
			return b1.Symbol() < b2.Symbol()
		}
		// else sort by date
		return b1.Time().Before(b2.Time())
	})
}

// updateLatest puts the last current data event to the current list.
func (d *Data) updateLatest(event DataEvent) {
	// check for nil map, else initialise the map
	if d.latest == nil {
		d.latest = make(map[string]DataEvent)
	}

	d.latest[event.Symbol()] = event
}

// updateList appends a data event to the data list for the corresponding symbol.
func (d *Data) updateList(event DataEvent) {
	// Check for nil map, else initialise the map
	if d.list == nil {
		d.list = make(map[string][]DataEvent)
	}

	d.list[event.Symbol()] = append(d.list[event.Symbol()], event)
}
