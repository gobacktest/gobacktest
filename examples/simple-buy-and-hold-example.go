package main

import (
	"log"

	"github.com/dirkolbrich/gobacktest"
)

func main() {
	// define symbols
	var symbols []string
	symbols = append(symbols, "DBK.DE")

	bt := gobacktest.New(symbols)
	log.Printf("bt: [%T] %v\n", bt, bt)
	bt.Run()

	bt2 := gobacktest.Test{}
	log.Printf("[%T] %v\n", bt2, bt2)
}
