package main

import (
	"log"

	"worksheetparser/debug"
	"worksheetparser/parser"
)

func main() {
	worksheet, err := parser.Parse("examples/example.wsp.txt")
	if err != nil {
		log.Fatal(err)
	}

	debug.LogPrintln("-------------------results-------------------")
	debug.LogPrint(worksheet)
}
