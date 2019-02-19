package gobacktest

import (
	"time"
)

// DataHandler is the combined data interface.
type DataHandler interface {
	DataLoader
	DataStreamer
	Reseter
}

// DataLoader defines how to load data into the data stream.
type DataLoader interface {
	Load([]string) error
}

// DataStreamer defines data stream functionality.
type DataStreamer interface {
	NextFromStream() (Data, bool)
	Stream() []Data
	History() []Data
	Latest(string) DataEvent
	List(string) []DataEvent
}

// DataStore is a basic data provider struct.
type DataStore struct {
	latest  map[string]Data
	list    map[string][]Data
	stream  []Data
	history []Data
}

// Load data events into a stream.
// This method satisfies the DataLoader interface, but should be overwritten
// by the specific data loading implementation.
func (d *DataStore) Load(s []string) error {
	return nil
}

// Reset implements Reseter to reset the data struct to a clean state with loaded data events.
func (d *DataStore) Reset() error {
	d.latest = nil
	d.list = nil
	d.stream = d.history
	d.history = nil
	return nil
}

// Stream returns the data stream.
func (d *DataStore) Stream() []Data {
	return d.stream
}

// SetStream sets the data stream.
func (d *DataStore) SetStream(stream []Data) {
	d.stream = stream
}

// NextFromStream returns the first element of the data stream,
// deletes it from the data stream and appends it to the historic data stream.
func (d *DataStore) NextFromStream() (dh Data, ok bool) {
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
func (d *DataStore) History() []Data {
	return d.history
}

// Latest returns the last known data event for a symbol.
func (d *DataStore) Latest(symbol string) Data {
	return d.latest[symbol]
}

// List returns the data event list for a symbol.
func (d *DataStore) List(symbol string) []Data {
	return d.list[symbol]
}

// updateLatest puts the last current data event to the current list.
func (d *DataStore) updateLatest(data Data) {
	// check for nil map, else initialise the map
	if d.latest == nil {
		d.latest = make(map[string]Data)
	}

	d.latest[data.Symbol()] = data
}

// updateList appends a data event to the data list for the corresponding symbol.
func (d *DataStore) updateList(data Data) {
	// Check for nil map, else initialise the map
	if d.list == nil {
		d.list = make(map[string][]Data)
	}

	d.list[data.Symbol()] = append(d.list[data.Symbol()], data)
}

// Data defines the basic interface for all types of data events.
type Data interface {
	Pricer
}

// Pricer defines the handling of the latest Price information
type Pricer interface {
	Price() float64
}

// BarData defines a bar event interface.
type BarData interface {
	Data
}

// Bar declares a data event for an OHLCV bar.
type Bar struct {
	Metric
	Time     time.Time
	Symbol   string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	AdjClose float64
	Volume   int64
}

// Price returns the close price of the bar event.
func (b Bar) Price() float64 {
	return b.Close
}

// TickData defines a tick event interface.
type TickData interface {
	Spreader
}

// Spreader declares functionality to get the spread of a tick.
type Spreader interface {
	Spread() float64
}

// Tick declares a data event for a price tick.
type Tick struct {
	Metric
	Time      time.Time
	Symbol    string
	Bid       float64
	Ask       float64
	BidVolume int64
	AskVolume int64
}

// Price returns the mid point of Bid and Ask prices.
func (t Tick) Price() float64 {
	latest := (t.Bid + t.Ask) / float64(2)
	return latest
}

// Spread returns the difference or spread of Bid and Ask prices.
func (t Tick) Spread() float64 {
	return t.Bid - t.Ask
}
