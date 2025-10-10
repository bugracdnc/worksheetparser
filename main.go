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

	isInteractiveGenerated := true
	if err := generator.GenerateInteractive(worksheet, interactiveOut); err != nil {
		fmt.Println("Error generating interactive version: ", err)
		isInteractiveGenerated = false
	}

	isPrintableGenerated := true
	if err := generator.GeneratePrintable(worksheet, printableOut); err != nil {
		fmt.Println("Error generating printable version: ", err)
		isPrintableGenerated = false
	}

	if isInteractiveGenerated || isPrintableGenerated {
		fmt.Println("Done generating files!\nCreated:")
		if isInteractiveGenerated {
			fmt.Println(" - ", interactiveOut)
		}
		if isPrintableGenerated {
			fmt.Println(" - ", printableOut)
		}
	} else {
		fmt.Println("Could not generate any files...")
	}

}
