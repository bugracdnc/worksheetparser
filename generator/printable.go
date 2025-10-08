package generator

import (
	"html/template"
	"os"
	"path/filepath"

	"worksheetparser/parser"
)

// GeneratePrintable renders a printable HTML version (no interactivity).
func GeneratePrintable(ws parser.Worksheet, outPath string) error {
	tmplPath := filepath.Join("templates", "printable.html.tmpl")

	tmpl, err := template.New("printable.html.tmpl").
		Funcs(template.FuncMap{
			"inc": func(i int) int { return i + 1 },
			"add": func(a, b int) int { return a + b },
		}).
		ParseFiles(tmplPath)

	if err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, ws)
}
