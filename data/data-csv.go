package data

import (
	"encoding/csv"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	gbt "github.com/dirkolbrich/gobacktest"
)

// BarEventFromCSVFile is a data struct, which loads the market data from csv files.
// It expands the underlying data struct
type BarEventFromCSVFile struct {
	gbt.Data
	FileDir string
}

// Load loads single data endpoints into a stream ordered by date (latest first).
func (d *BarEventFromCSVFile) Load(symbols []string) (err error) {
	// check file location
	if len(d.FileDir) == 0 {
		return errors.New("no directory for data provided: ")
	}

	// create a map for holding the file name for each symbol
	files := make(map[string]string)

	// read all files from directory
	if len(symbols) == 0 {
		files, err = fetchFilesFromDir(d.FileDir)
		if err != nil {
			return err
		}
		log.Printf("%v data files found.\n", len(files))
	}

	// construct filenames for provided symbols
	for _, symbol := range symbols {
		file := symbol + ".csv"
		files[symbol] = file
	}
	log.Printf("Loading %v symbol files.\n", len(files))

	// read file for each fileName
	for symbol, file := range files {
		log.Printf("Loading %s file for %s symbol.\n", file, symbol)

		// open file for corresponding symbol
		lines, err := readCSVFile(d.FileDir + file)
		if err != nil {
			return err
		}
		log.Printf("%v data lines found.\n", len(lines))

		// for each found record create an event
		for _, line := range lines {
			event, err := createBarEventFromLine(line, symbol)
			if err != nil {
				// what happens if line could not be parsed - needs logging
				// log.Println(line)
				// log.Println(err)
				continue
			}
			// append event to data stream
			d.Data.SetStream(append(d.Data.Stream(), event))
		}
	}
	// sort data stream
	d.Data.SortStream()

	return nil
}

// fetchFilesFromDir returns a map of all filenames in a directory
// e.g map{"BAS.DE": "BAS.DE.csv"}
func fetchFilesFromDir(dir string) (m map[string]string, err error) {
	// read filenames from directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return m, err
	}

	// initialise the map
	m = make(map[string]string)

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

		name := filename[0 : len(filename)-len(extension)]
		m[name] = filename
	}
	return m, nil
}

// readCSVFile opens and reads a csv file line by line
// and returns a slice with a key/value map for each line
func readCSVFile(path string) (lines []map[string]string, err error) {
	log.Printf("Loading from %s.\n", path)
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
	// var lines []map[string]string

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

// createBarEventFromLine takes a key/value map and a string and builds a bar struct
func createBarEventFromLine(line map[string]string, symbol string) (bar *gbt.Bar, err error) {
	// parse each string in line to corresponding record value
	date, err := time.Parse("2006-01-02", line["Date"])
	if err != nil {
		return bar, err
	}

	openPrice, err := strconv.ParseFloat(line["Open"], 64)
	if err != nil {
		return bar, err
	}

	highPrice, err := strconv.ParseFloat(line["High"], 64)
	if err != nil {
		return bar, err
	}

	lowPrice, err := strconv.ParseFloat(line["Low"], 64)
	if err != nil {
		return bar, err
	}

	closePrice, err := strconv.ParseFloat(line["Close"], 64)
	if err != nil {
		return bar, err
	}

	adjClosePrice, err := strconv.ParseFloat(line["Adj Close"], 64)
	if err != nil {
		return bar, err
	}
	volume, err := strconv.ParseInt(line["Volume"], 10, 64)
	if err != nil {
		return bar, err
	}

	// create and populate new event
	event := &gbt.Event{}
	event.SetTime(date)
	event.SetSymbol(strings.ToUpper(symbol))

	metric := &gbt.Metric{}

	bar = &gbt.Bar{
		Event:    *event,
		Metric:   *metric,
		Open:     openPrice,
		High:     highPrice,
		Low:      lowPrice,
		Close:    closePrice,
		AdjClose: adjClosePrice,
		Volume:   volume,
	}

	return bar, nil
}
