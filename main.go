package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"worksheetparser/debug"
	"worksheetparser/generator"
	"worksheetparser/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please drag and drop a worksheet file onto this executable")
		fmt.Println("or run worksheetpager <filename> in a terminal")
		return
	}

	inputPath := os.Args[1]
	dir := filepath.Dir(inputPath)

	interactiveOut := filepath.Join(dir, "interactive.html")
	printableOut := filepath.Join(dir, "printable.html")

	worksheet, err := parser.Parse(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	debug.LogPrintln("-------------------results-------------------")
	debug.LogPrint(worksheet)

	fmt.Println("Generating HTML worksheets...")

	if err := generator.GenerateInteractive(worksheet, interactiveOut); err != nil {
		fmt.Println("Error generating interactive version: ", err)
	}

	if err := generator.GeneratePrintable(worksheet, printableOut); err != nil {
		fmt.Println("Error generating printable version: ", err)
	}

	fmt.Printf("Done generating files!\nCreated:\n - %s\n - %s\n", interactiveOut, printableOut)
}
