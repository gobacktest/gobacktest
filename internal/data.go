package internal

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// DataHandler is the basic data interface
type DataHandler interface {
	Load([]string) error
	LoadAll() error
}

// Data is a basic data struct
type Data struct {
	dataStream []string
}

// Load loads single data endpoints.
func (d Data) Load(symbols []string) error {
	for _, v := range symbols {
		filename := v + ".csv"
		file, err := os.Open("../data/" + filename)
		if err != nil {
			log.Println(err)
			return err
		}
		defer file.Close()

		// creat scanner on top of file
		reader := csv.NewReader(file)
		// set delimeter
		reader.Comma = ','

		// read first line for keys and fill in array
		keys, err := reader.Read()

		// create stream, a slice for holding the different maps
		var stream []map[string]string

		// read each line and create a map of values combined to the keys
		for line, err := reader.Read(); err == nil; line, err = reader.Read() {
			m := make(map[string]string)
			for i, v := range line {
				m[keys[i]] = v
			}
			// put found record as map into stream holder item
			stream = append(stream, m)
		}

		// for each found record create an event
		for i, r := range stream {
			if i > 5 {
				break
			}
			event := createBarEvent(r)
			log.Printf("event [%T] %+v\n", event, event)
		}
	}

	return nil
}

// LoadAll loads the complete data.
func (d Data) LoadAll() error {
	return nil
}

func createBarEvent(record map[string]string) BarEvent {
	log.Printf("record [%T] %v\n", record, record)

	// parse strings to float
	floats := make(map[string]float64)
	for k, v := range record {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Fatal("Could not parse string to float")
		}
	}

	log.Printf("record [%T] %v\n", record, record)

	be := BarEvent{
		// {date: {0, 0, nil}},
		openPrice: record["Open"],
		highPrice: record["High"],
		lowPrice:  record["Low"],
		// {closePrice: record[4]},
		// {volume: record[6]},
	}

	return be
}
