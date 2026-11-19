package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bugracdnc/worksheetparser/internal/generator"
	"github.com/bugracdnc/worksheetparser/internal/models"
	"github.com/bugracdnc/worksheetparser/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Drag and drop a worksheet file onto this executable,")
		fmt.Println("or run: worksheetparser <file.wsp.txt>")
		return
	}

	inputPath := os.Args[1]
	worksheet, warnings, err := parser.ParseFile(inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, w := range warnings {
		fmt.Fprintln(os.Stderr, " ! Warning:", w)
	}

	dir, file := filepath.Split(inputPath)
	fmt.Println("Generating HTML worksheets...")
	for _, out := range []struct {
		suffix string
		render func(io.Writer, models.Worksheet) error
	}{
		{"_interactive.html", generator.RenderInteractive},
		{"_printable.html", generator.RenderPrintable},
	} {
		path := filepath.Join(dir, file+out.suffix)
		if err := writeFile(path, worksheet, out.render); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(" - Created:", path)
	}
}

func writeFile(path string, ws models.Worksheet, render func(io.Writer, models.Worksheet) error) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating %s: %w", path, err)
	}
	if err := render(f, ws); err != nil {
		f.Close()
		return fmt.Errorf("rendering %s: %w", path, err)
	}
	return f.Close()
}
