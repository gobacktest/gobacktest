package internal

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"io/ioutil"
	"path/filepath"
)

/***** Define DataHandler interface *****/

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
	Next() (EventHandler, bool)
	Stream() []EventHandler
	History() []EventHandler
	Current(string) EventHandler
	List(string) []EventHandler
}

/***** Basic Data struct with implemented interface methods *****/

// Data is a basic data struct
type Data struct {
	current       map[string]EventHandler
	list          map[string][]EventHandler
	stream        []EventHandler
	streamHistory []EventHandler
}

// Load loads data endpoints into a stream
func (d *Data) Load(s []string) error {
	return nil
}

// Stream returns the data stream
func (d *Data) Stream() []EventHandler {
	return d.stream
}

// Next returns the first element of the data stream
// deletes it from the stream and appends it to history
func (d *Data) Next() (event EventHandler, ok bool) {
	// check for element in datastream
	if len(d.stream) == 0 {
		return event, false
	}

	event = d.stream[0]
	d.stream = d.stream[1:] // delete first element from stream
	d.streamHistory = append(d.stream, event)

	if event.(type) == BarEvent {
		event = d.calculateBarMetrics(event, d.List(event.Symbol()))
	}

	// update list of current data events
	d.updateCurrent(event)

	// update list of data events for single symbol
	d.updateList(event)

	return event, true
}

// History returns the historic data stream
func (d *Data) History() []EventHandler {
	return d.streamHistory
}

// Current returns the latest data event for a symbol.
func (d *Data) Current(symbol string) EventHandler {
	return d.current[symbol]
}

// updateCurrent puts the last current data event to the current list.
func (d *Data) updateCurrent(event EventHandler) {
	// check for nil map, else initialise the map
	if d.current == nil {
		d.current = make(map[string]EventHandler)
	}

	d.current[event.Symbol()] = event
}

// List returns the data event list for a symbol.
func (d *Data) List(symbol string) []EventHandler {
	return d.list[symbol]
}

// updateList appends an event to the data list.
func (d *Data) updateList(event EventHandler) {
	// Check for nil map, else initialise the map
	if d.list == nil {
		d.list = make(map[string][]EventHandler)
	}

	d.list[event.Symbol()] = append(d.list[event.Symbol()], event)
}

// calculateBarMetrics calculates metrics for a bar event
func (d *Data) calculateBarMetrics(bar BarEvent, list []EventHandler) BarEvent {
	return bar
}

/***** Concrete BarEventFromCSVFileData struct *****/

// BarEventFromCSVFileData is a data struct, which loads the market data from csv files.
// It expands the underlying data struct
type BarEventFromCSVFileData struct {
	Data
	FileDir string
}

// Load loads single data endpoints into a stream ordered by date (latest first).
func (d *BarEventFromCSVFileData) Load(symbols []string) error {
	// check file location
	if d.FileDir == nil {
		return error.New("No directory for data provided.")
	}

	// create a map for holding the file name for each symbol
	var files map[string]string

	// read all files from directory
	if len(symbols) == 0 {
		files, err = fetchFilesFromDir(d.FileDir)
		if err != nil {
			return err
		}
	}

	// construct filenames for provided symbols
	for _, symbol := range symbols {
		file := symbol + ".csv"
		files[symbol] = file
	}

	// read file for each fileName
	for symbol, file := range files {
		// open file for corresponding symbol
		lines, err := readCSVFile(d.FileDir + file)
		if err != nil {
			return err
		}

		// for each found record create an event
		for _, line := range lines {
			event, err := createBarEventFromLine(line, symbol)
			if err != nil {
				log.Println(err)
			}
			d.stream = append(d.stream, event)
		}
	}

	// sort data stream
	d.stream = sortStream(d.Stream())

	return nil
}

// fetchFilesFromDir returns a map of all filenames in a directory
// e.g map{"BAS.DE": "BAS.DE.csv"}
func fetchFilesFromDir(dir string) (m map[string]string, error) {
	// read filenames from directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return m, err
	}
	// read filenames from directory
	for _, file := range files {
		// file is directory
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		extension := filepath.Ext(filename)
		// file is not CSV
		if extension != ".csv" {
			continue
		}

		name := filename[0 : len(filename) - len(extension)]
		m[name] = filename
	}
	return m, nil
}

// readCSVFile opens and reads a csv file line by line
// and returns a slice with a key/value map for each line
func readCSVFile(path string) (lines []map[string]string, error) {
	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create scanner on top of file
	reader := csv.NewReader(file)
	// set delimeter
	reader.Comma = ','
	// read first line for keys and fill in array
	keys, err := reader.Read()

	// create a slice for holding the different maps of each line
	var lines []map[string]string

	// read each line and create a map of values combined to the keys
	for line, err := reader.Read(); err == nil; line, err = reader.Read() {
		l := make(map[string]string)
		for i, v := range line {
			l[keys[i]] = v
		}
		// put found line as map into stream holder item
		lines = append(lines, l)
	}

	return lines, nil
}

// createBarEventFromLine takes a key/value line and builds a BarEvent struct
func createBarEventFromLine(line map[string]string, symbol string) (BarEvent, error) {
	// parse each string in line to corresponding record value
	date, _ := time.Parse("2006-01-02", line["Date"])
	openPrice, _ := strconv.ParseFloat(line["Open"], 64)
	highPrice, _ := strconv.ParseFloat(line["High"], 64)
	lowPrice, _ := strconv.ParseFloat(line["Low"], 64)
	closePrice, _ := strconv.ParseFloat(line["Close"], 64)
	adjClosePrice, _ := strconv.ParseFloat(line["Adj Close"], 64)
	volume, _ := strconv.ParseInt(line["Volume"], 10, 64)

	// create and populate new event
	event := BarEvent{
		Event:         Event{timestamp: date, symbol: strings.ToUpper(symbol)},
		OpenPrice:     openPrice,
		HighPrice:     highPrice,
		LowPrice:      lowPrice,
		ClosePrice:    closePrice,
		AdjClosePrice: adjClosePrice,
		Volume:        volume,
	}

	return event, nil
}

// sortStream sorts the dataStream
func sortStream(stream []EventHandler) []EventHandler {
	sort.Slice(stream, func(i, j int) bool {
		// cast EventHandler interface{} to concrete BarEvent{} implementation
		b1 := stream[i].(BarEvent)
		b2 := stream[j].(BarEvent)

		// if date is equal sort by symbol
		if b1.Timestamp().Equal(b2.Timestamp()) {
			return b1.Symbol() < b2.Symbol()
		}
		// else sort by date
		return b1.Timestamp().Before(b2.Timestamp())
	})

	return stream
}
