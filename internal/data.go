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
	Stream() []EventHandler
}

// Data is a basic data struct
type Data struct {
	dataStream []EventHandler
}

// Load loads data endpoints into a stream
func (d *Data) Load(s []string) error {
	return nil
}

// Stream returns the data stream
func (d *Data) Stream() []EventHandler {
	return d.dataStream
}

// BarEventFromCSVFileData is a data struct, which loads the market data from csv files
type BarEventFromCSVFileData struct {
	dataStream []EventHandler
	FileDir    string
}

// Load loads single data endpoints into a stream ordered by date (latest first).
func (d *BarEventFromCSVFileData) Load(symbols []string) error {
	//log.Printf("Loading %d symbols from file: %v\n", len(symbols), symbols)
	// startLoad := time.Now()

	// load elements for each symbol
	for _, symbol := range symbols {
		// set filename and filepath
		fileName := symbol + ".csv"

		// startReadCSV := time.Now()
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

		// elapsedReadCSV := time.Since(startReadCSV)
		// log.Printf("Loading %s file with %d entrys took %s", fileName, len(lines), elapsedReadCSV)
	}

	// elapsed := time.Since(startLoad)
	// log.Printf("Loading %d files took %s", len(symbols), elapsed)

	//log.Printf("[%T] %v\n", d.dataStream, d.dataStream)

	startSortStream := time.Now()
	d.dataStream = sortStream(d.Stream())
	elapsedSortStream := time.Since(startSortStream)
	log.Printf("Sorting %T with %d entrys took %s", d.dataStream, len(d.dataStream), elapsedSortStream)

	//log.Printf("[%T] %+v\n", d.dataStream, d.dataStream)

	return nil
}

// Stream returns the data stream
func (d *BarEventFromCSVFileData) Stream() []EventHandler {
	return d.dataStream
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
		Date:          date,
		Symbol:        strings.ToUpper(symbol),
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
		if bar1.Date.Equal(bar2.Date) {
			return bar1.Symbol < bar2.Symbol
		}

		// else sort by date
		return bar1.Date.Before(bar2.Date)
	})

	return stream
}
