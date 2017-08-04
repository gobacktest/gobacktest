package internal

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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
	StreamIsEmpty() bool
	StreamHistory() []EventHandler
}

/***** Basic Data struct with implemented interface methods *****/

// Data is a basic data struct
type Data struct {
	dataCurrent       map[string]EventHandler
	dataList          map[string][]EventHandler
	dataStream        []EventHandler
	dataStreamHistory []EventHandler
}

// Load loads data endpoints into a stream
func (d *Data) Load(s []string) error {
	return nil
}

// Stream returns the data stream
func (d *Data) Stream() []EventHandler {
	return d.dataStream
}

// Next returns the first element of the data stream
// deletes it from the stream and appends it to history
func (d *Data) Next() (event EventHandler, ok bool) {
	if len(d.dataStream) == 0 {
		return event, false
	}

	event = d.dataStream[0]
	d.dataStream = d.dataStream[1:] // delete from dataStream
	d.dataStreamHistory = append(d.dataStream, event)

	// update list of current data events
	d.updateCurrent(event)

	// update list of data events for single symbol
	d.updateList(event)

	return event, true
}

// StreamIsEmpty checks if the data stream is empty
func (d *Data) StreamIsEmpty() bool {
	return false
}

// StreamHistory returns the historic data stream
func (d *Data) StreamHistory() []EventHandler {
	return d.dataStreamHistory
}

func (d *Data) updateCurrent(e EventHandler) {
	// Check for nil map, else initialise the map
	if d.dataCurrent == nil {
		d.dataCurrent = make(map[string]EventHandler)
	}

	d.dataCurrent[e.Symbol()] = e
}

func (d *Data) updateList(e EventHandler) {
	// Check for nil map, else initialise the map
	if d.dataList == nil {
		d.dataList = make(map[string][]EventHandler)
	}

	d.dataList[e.Symbol()] = append(d.dataList[e.Symbol()], e)
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
	// load elements for each symbol
	for _, symbol := range symbols {
		// set filename and filepath
		fileName := symbol + ".csv"
		// open file for corresponding symbol
		lines, err := readCSVFile(fileName, d.FileDir)
		if err != nil {
			log.Println(err)
			return err
		}

		// for each found record create an event
		for _, l := range lines {
			event, err := createBarEventFromLine(l, symbol)
			if err != nil {
				log.Println(err)
			}
			d.dataStream = append(d.dataStream, event)
		}
	}

	// sort data stream
	d.dataStream = sortStream(d.Stream())

	return nil
}

// readCSVFile opens and reads a csv file line by line
// and returns a slice with a key/value map for each line
func readCSVFile(fileName, fileDir string) ([]map[string]string, error) {
	// open file
	file, err := os.Open(fileDir + fileName)
	if err != nil {
		log.Println(err)
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

// createBarEventFromLine takes a key/value line and buils a BarEvent struct
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
	be := BarEvent{
		Event:         Event{timestamp: date, symbol: strings.ToUpper(symbol)},
		OpenPrice:     openPrice,
		HighPrice:     highPrice,
		LowPrice:      lowPrice,
		ClosePrice:    closePrice,
		AdjClosePrice: adjClosePrice,
		Volume:        volume,
	}

	return be, nil
}

// sortStream sorts the dataStream
func sortStream(stream []EventHandler) []EventHandler {
	sort.Slice(stream, func(i, j int) bool {
		// cast EventHandler interface{} to concrete BarEvent{} implementation
		bar1 := stream[i].(BarEvent)
		bar2 := stream[j].(BarEvent)

		// if date is equal sort by symbol
		if bar1.Timestamp().Equal(bar2.Timestamp()) {
			return bar1.Symbol() < bar2.Symbol()
		}

		// else sort by date
		return bar1.Timestamp().Before(bar2.Timestamp())
	})

	return stream
}
