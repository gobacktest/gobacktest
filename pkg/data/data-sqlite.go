package data

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dirkolbrich/gobacktest/pkg/backtest"
)

// BarEventFromSQLiteData is a data struct, which loads the market data from a SQLite database.
// It expands the underlying data struct
type BarEventFromSQLiteData struct {
	backtest.Data
	FileDir string
}

// Load loads single data endpoints into a stream ordered by date (latest first).
func (d *BarEventFromSQLiteData) Load(symbols []string) (err error) {
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
				log.Println(err)
			}
			d.stream = append(d.stream, event)
		}
	}

	// sort data stream
	d.sortStream()

	return nil
}

// fetchFilesFromDir returns a map of all filenames in a directory
// e.g map{"BAS.DE": "BAS.DE.csv"}
func fetchTableFromDB(dir string) (m map[string]string, err error) {
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

// createBarEventFromLine takes a key/value map and a string and builds a bar struct
func createBarEventFromEntry(line map[string]string, symbol string) (backtest.BarEvent, error) {
	// parse each string in line to corresponding record value
	date, _ := time.Parse("2006-01-02", line["Date"])
	openPrice, _ := strconv.ParseFloat(line["Open"], 64)
	highPrice, _ := strconv.ParseFloat(line["High"], 64)
	lowPrice, _ := strconv.ParseFloat(line["Low"], 64)
	closePrice, _ := strconv.ParseFloat(line["Close"], 64)
	adjClosePrice, _ := strconv.ParseFloat(line["Adj Close"], 64)
	volume, _ := strconv.ParseInt(line["Volume"], 10, 64)

	// create and populate new event
	event := bar{
		event:         event{timestamp: date, symbol: strings.ToUpper(symbol)},
		openPrice:     openPrice,
		highPrice:     highPrice,
		lowPrice:      lowPrice,
		closePrice:    closePrice,
		adjClosePrice: adjClosePrice,
		volume:        volume,
	}

	return event, nil
}
