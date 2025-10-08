package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"worksheetparser/parser"
)

func GenerateInteractive(worksheet parser.Worksheet, outPath string) error {
	tmplPath := filepath.Join("templates", "interactive.html.tmpl")
	fmt.Println("Loading template:", tmplPath)

	tmpl, err := template.New("interactive.html.tmpl").Funcs(template.FuncMap{
		"inc": func(i int) int { return i + 1 }}).
		ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Printf("Title: %s\nQuestions: %d\n", worksheet.Title, len(worksheet.Questions))

	err = tmpl.Execute(f, worksheet)
	if err != nil {
		fmt.Println("Template execution error:", err)
		return err
	}
	return nil
}
