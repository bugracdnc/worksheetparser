package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"worksheetparser/debug"
	"worksheetparser/internal/generator"
	"worksheetparser/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please drag and drop a worksheet file onto this executable")
		fmt.Println("or run worksheetpager <filename> in a terminal")
		return
	}

	inputPath := os.Args[1]
	dir := filepath.Dir(inputPath)

	worksheet, err := parser.Parse(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	interactiveOut := filepath.Join(dir, "interactive.html")
	printableOut := filepath.Join(dir, "printable.html")

	debug.LogPrintln("-------------------results-------------------")
	debug.LogPrint(worksheet)

	fmt.Println("Generating HTML worksheets...")

	fileInteractive, err := os.Create(interactiveOut)
	if err != nil {
		log.Fatalf("Error creating interactive file: %v\n", err)
	} else {
		defer fileInteractive.Close()

		if err := generator.RenderInteractive(fileInteractive, worksheet); err != nil {
			log.Fatalf("Error rendering interactive: %v\n", err)
		} else {
			fmt.Println(" - Created: ", interactiveOut)
		}
	}

	filePrintable, err := os.Create(printableOut)
	if err != nil {
		log.Fatalf("Error creating interactive file: %v\n", err)
	} else {
		defer filePrintable.Close()

		if err := generator.RenderPrintable(filePrintable, worksheet); err != nil {
			log.Fatalf("Error rendering interactive: %v\n", err)
		} else {
			fmt.Println(" - Created: ", printableOut)
		}
	}

}
