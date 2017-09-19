package backtest

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

// Load loads data endpoints into a stream
func (d *Data) Load(s []string) error {
	return nil
}

// SetStream sets the data stream
func (d *Data) SetStream(stream []DataEventHandler)  {
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
	d.streamHistory = append(d.stream, dh)

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

// updateLatest puts the last current data event to the current list.
func (d *Data) updateLatest(event DataEventHandler) {
	// check for nil map, else initialise the map
	if d.latest == nil {
		d.latest = make(map[string]DataEventHandler)
	}

	d.latest[event.GetSymbol()] = event
}

// List returns the data event list for a symbol.
func (d *Data) List(symbol string) []DataEventHandler {
	return d.list[symbol]
}

// updateList appends an event to the data list.
func (d *Data) updateList(event DataEventHandler) {
	// Check for nil map, else initialise the map
	if d.list == nil {
		d.list = make(map[string][]DataEventHandler)
	}

	d.list[event.GetSymbol()] = append(d.list[event.GetSymbol()], event)
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
