package backtest

import (
	"sort"
)

// DataHandler is the combined data interface
type DataHandler interface {
	DataLoader
	DataStreamer
	Reseter
}

// DataLoader is the interface loading the data into the data stream
type DataLoader interface {
	Load([]string) error
}

// DataStreamer is the interface returning the data streams
type DataStreamer interface {
	Next() (DataEventHandler, bool)
	Stream() []DataEventHandler
	History() []DataEventHandler
	Latest(string) DataEventHandler
	List(string) []DataEventHandler
}

// Data is a basic data struct
type Data struct {
	latest        map[string]DataEventHandler
	list          map[string][]DataEventHandler
	stream        []DataEventHandler
	streamHistory []DataEventHandler
}

// Load loads data endpoints into a stream.
// This method satisfies the DataLoeder interface, but should be overwritten
// by the specific data loading implamentation.
func (d *Data) Load(s []string) error {
	return nil
}

// Reset implements the Reseter interface and rests the data struct to a clean state with loaded data points
func (d *Data) Reset() {
	d.latest = nil
	d.list = nil
	d.stream = d.streamHistory
	d.streamHistory = nil
}

// SetStream sets the data stream
func (d *Data) SetStream(stream []DataEventHandler) {
	d.stream = stream
}

// Stream returns the data stream
func (d *Data) Stream() []DataEventHandler {
	return d.stream
}

// Next returns the first element of the data stream
// deletes it from the stream and appends it to history
func (d *Data) Next() (dh DataEventHandler, ok bool) {
	// check for element in datastream
	if len(d.stream) == 0 {
		return dh, false
	}

	dh = d.stream[0]
	d.stream = d.stream[1:] // delete first element from stream
	d.streamHistory = append(d.streamHistory, dh)

	// update list of current data events
	d.updateLatest(dh)
	// update list of data events for single symbol
	d.updateList(dh)

	return dh, true
}

// History returns the historic data stream
func (d *Data) History() []DataEventHandler {
	return d.streamHistory
}

// Latest returns the last known data event for a symbol.
func (d *Data) Latest(symbol string) DataEventHandler {
	return d.latest[symbol]
}

// List returns the data event list for a symbol.
func (d *Data) List(symbol string) []DataEventHandler {
	return d.list[symbol]
}

// SortStream sorts the dataStream
func (d *Data) SortStream() {
	sort.Slice(d.stream, func(i, j int) bool {
		b1 := d.stream[i]
		b2 := d.stream[j]

		// if date is equal sort by symbol
		if b1.GetTime().Equal(b2.GetTime()) {
			return b1.GetSymbol() < b2.GetSymbol()
		}
		// else sort by date
		return b1.GetTime().Before(b2.GetTime())
	})
}

// updateLatest puts the last current data event to the current list.
func (d *Data) updateLatest(event DataEventHandler) {
	// check for nil map, else initialise the map
	if d.latest == nil {
		d.latest = make(map[string]DataEventHandler)
	}

	d.latest[event.GetSymbol()] = event
}

// updateList appends an event to the data list.
func (d *Data) updateList(event DataEventHandler) {
	// Check for nil map, else initialise the map
	if d.list == nil {
		d.list = make(map[string][]DataEventHandler)
	}

	d.list[event.GetSymbol()] = append(d.list[event.GetSymbol()], event)
}
