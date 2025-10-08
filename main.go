package main

import (
	"log"

	"worksheetparser/debug"
	"worksheetparser/generator"
	"worksheetparser/parser"
)

func main() {
	worksheet, err := parser.Parse("examples/example.wsp.txt")
	if err != nil {
		log.Fatal(err)
	}

	debug.LogPrintln("-------------------results-------------------")
	debug.LogPrint(worksheet)

	generator.GenerateInteractive(worksheet, "interactive.html")
	generator.GeneratePrintable(worksheet, "printable.html")
}
