package internal

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dirkolbrich/gobacktest/internal/utils"
)

// DataHandler is the basic data interface
type DataHandler interface {
	Load([]string) error
	LoadAll() error
}

// Data is a basic data struct
type Data struct {
	dataStream []EventHandler
}

// BarEventFromCSVFileData is a data struct, which loads the market data from csv files
type BarEventFromCSVFileData struct {
	dataStream barStream
	FilePath   string
}

// Load loads single data endpoints into a stream ordered by date (latest first).
func (d *BarEventFromCSVFileData) Load(symbols []string) error {
	log.Printf("Loading %d symbols from file: %v\n", len(symbols), symbols)

	startLoad := time.Now()

	// load elements for each symbol
	for _, symbol := range symbols {
		// set filename and filepath
		fileName := symbol + ".csv"

		startReadCSV := time.Now()
		// open file for corresponding symbol
		lines, err := readCSVFile(fileName, d.FilePath)
		if err != nil {
			log.Println(err)
			return err
		}

		// for each found record create an event
		for _, l := range lines {
			event := createBarEventFromLine(l, symbol)
			d.dataStream = append(d.dataStream, event)
		}
		// log.Printf("d.dataStream: [%T] %+v\n", d.dataStream, d.dataStream)
		elapsedReadCSV := time.Since(startReadCSV)
		log.Printf("Loading %s file with %d entrys took %s", fileName, len(lines), elapsedReadCSV)
	}

	elapsed := time.Since(startLoad)
	log.Printf("Loading %d files took %s", len(symbols), elapsed)

	// order dataStream by Date
	sort.Sort(d.dataStream)

	return nil
}

// LoadAll loads the complete data.
func (d *BarEventFromCSVFileData) LoadAll() error {
	return nil
}

func readCSVFile(fileName, filePath string) ([]map[string]string, error) {
	// open file
	file, err := os.Open(filePath + fileName)
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

func createBarEventFromLine(l map[string]string, symbol string) BarEvent {
	// create in between struct to hold temporary data
	type record struct {
		date          time.Time
		openPrice     float64
		highPrice     float64
		lowPrice      float64
		closePrice    float64
		adjClosePrice float64
		volume        int64
	}
	r := record{}

	// parse each string in line to corresponding record value
	r.date, _ = time.Parse("2006-01-02", l["Date"])
	r.openPrice, _ = strconv.ParseFloat(l["Open"], 64)
	r.highPrice, _ = strconv.ParseFloat(l["High"], 64)
	r.lowPrice, _ = strconv.ParseFloat(l["Low"], 64)
	r.closePrice, _ = strconv.ParseFloat(l["Close"], 64)
	r.adjClosePrice, _ = strconv.ParseFloat(l["Adj Close"], 64)
	r.volume, _ = strconv.ParseInt(l["Volume"], 10, 64)

	// create PRiceParser to convert float64 into int64
	pp := utils.PriceParser{}

	be := BarEvent{
		date:          r.date,
		symbol:        strings.ToUpper(symbol),
		openPrice:     pp.Parse(r.openPrice),
		highPrice:     pp.Parse(r.highPrice),
		lowPrice:      pp.Parse(r.lowPrice),
		closePrice:    pp.Parse(r.closePrice),
		adjClosePrice: pp.Parse(r.adjClosePrice),
		volume:        r.volume,
	}

	return be
}

// implementing sorting function into stream
type barStream []BarEvent

func (s barStream) Len() int {
	return len(s)
}

func (s barStream) Less(i, j int) bool {
	// if date is equal sort by symbol
	if s[i].date.Equal(s[j].date) {
		return s[i].symbol < s[j].symbol
	}

	return s[i].date.Before(s[j].date)
}

func (s barStream) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
